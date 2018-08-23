package web

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tags-drive/internal/params"
	"github.com/ShoshinNikita/tags-drive/internal/storage"
)

const (
	maxSize = 50000000 // 50MB
)

var (
	ErrEmptyFilename = errors.New("name of a file can't be empty")
)

// multiplyResponse is used as response by POST /api/files and DELETE /api/files
type multiplyResponse struct {
	Filename string `json:"filename"`
	IsError  bool   `json:"isError"`
	Error    string `json:"error"`
	Status   string `json:"status"` // Status isn't empty when IsError == false
}

/* Files */

// POST /api/files (multipart/form-data)
//
// Response: json list of strings with status of files uploading
//
func upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		switch err {
		case http.ErrNotMultipart:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	respChan := make(chan multiplyResponse, 50)
	wg := new(sync.WaitGroup)

	for _, f := range r.MultipartForm.File["files"] {
		wg.Add(1)
		go func(header *multipart.FileHeader) {
			defer wg.Done()
			err := storage.UploadFile(header, []string{})
			if err != nil {
				respChan <- multiplyResponse{Filename: header.Filename, IsError: true, Error: err.Error()}
			} else {
				respChan <- multiplyResponse{Filename: header.Filename, Status: "uploaded"}
			}
		}(f)
	}
	wg.Wait()
	close(respChan)

	var response []multiplyResponse
	for resp := range respChan {
		response = append(response, resp)
	}

	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(response)
}

// DELETE /api/files?file=file1&file=file2
//
// Response: -
//
func deleteFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filenames := r.Form["file"]
	if len(filenames) == 0 {
		http.Error(w, "list of files for deleting can't be empty", http.StatusBadRequest)
		return
	}

	respChan := make(chan multiplyResponse, 50)
	wg := new(sync.WaitGroup)

	for _, filename := range filenames {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			err := storage.DeleteFile(f)
			if err != nil {
				respChan <- multiplyResponse{Filename: f, IsError: true, Error: err.Error()}
			} else {
				respChan <- multiplyResponse{Filename: f, Status: "deleted"}
			}
		}(filename)
	}

	wg.Wait()
	close(respChan)

	var response []multiplyResponse
	for resp := range respChan {
		response = append(response, resp)
	}

	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(response)
}

// PUT /api/files?file=123&newname=567&tags=tag1,tag2,tag3
// newname or tags can be skipped
//
// Response: -
//
func changeFile(w http.ResponseWriter, r *http.Request) {
	var (
		filename = r.FormValue("file")
		newName  = r.FormValue("newname")
		tags     = func() []string {
			t := r.FormValue("tags")
			if t == "" {
				return []string{}
			}

			return strings.Split(t, ",")
		}()
	)

	if filename == "" {
		http.Error(w, ErrEmptyFilename.Error(), http.StatusBadRequest)
		return
	}

	// Change tags, at first
	if len(tags) != 0 {
		err := storage.ChangeTags(filename, tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Change filename
	if newName != "" {
		// We can skip checking of invalid characters, because Go will return an error
		err := storage.RenameFile(filename, newName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GET /api/files?sort=(name|size|time)&order(asc|desc)&tags=first,second,third&mode=(or|and|not)&search=abc
// tags - list of tags separated by ',' (can be empty, then all files will be returned)
// First elements in params is default (name, asc and etc.)
//
// Response: json array of files
//
func returnFiles(w http.ResponseWriter, r *http.Request) {
	var (
		order = getParam("asc", r.FormValue("order"), "asc", "desc")
		tags  = func() []string {
			t := r.FormValue("tags")
			if t == "" {
				return []string{}
			}

			return strings.Split(t, ",")
		}()
		search = r.FormValue("search")

		tagMode  = storage.ModeOr
		sortMode = storage.SortByNameAsc
	)

	// Set sortMode
	// Can skip default
	switch r.FormValue("sort") {
	case "name":
		if order == "asc" {
			sortMode = storage.SortByNameAsc
		} else {
			sortMode = storage.SortByNameDesc
		}
	case "size":
		if order == "asc" {
			sortMode = storage.SortBySizeAsc
		} else {
			sortMode = storage.SortBySizeDecs
		}
	case "time":
		if order == "asc" {
			sortMode = storage.SortByTimeAsc
		} else {
			sortMode = storage.SortByTimeDesc
		}
	}

	// Set tagMode
	// Can skip default
	switch r.FormValue("mode") {
	case "or":
		tagMode = storage.ModeOr
	case "and":
		tagMode = storage.ModeAnd
	case "not":
		tagMode = storage.ModeNot
	}

	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(storage.Get(tagMode, sortMode, tags, search))
}

func getParam(def, passed string, options ...string) (s string) {
	s = def
	if passed == def {
		return
	}
	for _, opt := range options {
		if passed == opt {
			return passed
		}
	}

	return
}

// GET /api/files/recent?number=5 (5 is a default value)
//
// Response: json array of files
//
func returnRecentFiles(w http.ResponseWriter, r *http.Request) {
	number := func() int {
		s := r.FormValue("number")
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			n = 5
		}
		return int(n)
	}()

	files := storage.GetRecent(number)

	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(files)
}

/* Tags */
