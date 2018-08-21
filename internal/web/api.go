package web

import (
	"encoding/json"
	"fmt"
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

	msgChan := make(chan string, 50)
	wg := new(sync.WaitGroup)

	for _, f := range r.MultipartForm.File["files"] {
		wg.Add(1)
		go func(header *multipart.FileHeader) {
			defer wg.Done()
			err := storage.UploadFile(header, []string{})
			if err != nil {
				msgChan <- fmt.Sprintf("%s: %s", header.Filename, err)
			} else {
				msgChan <- fmt.Sprintf("%s: %s", header.Filename, "done")
			}
		}(f)
	}
	wg.Wait()
	close(msgChan)

	var messages []string

	for msg := range msgChan {
		messages = append(messages, msg)
	}

	json.NewEncoder(w).Encode(messages)
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

	msgChan := make(chan string, 50)
	wg := new(sync.WaitGroup)

	for _, filename := range filenames {
		wg.Add(1)
		go func(f string) {
			err := storage.DeleteFile(f)
			// Log only errors
			if err != nil {
				msgChan <- fmt.Sprintf("%s: %s", f, err.Error())
			}
			wg.Done()
		}(filename)
	}

	wg.Wait()
	close(msgChan)

	var messages []string
	for msg := range msgChan {
		messages = append(messages, msg)
	}

	json.NewEncoder(w).Encode(messages)
}

// PUT /api/files?oldname=123&newname=567
//
// Response: -
//
func renameFile(w http.ResponseWriter, r *http.Request) {
	var (
		oldName = r.FormValue("oldname")
		newName = r.FormValue("newname")
	)

	if oldName == "" || newName == "" {
		http.Error(w, ErrEmptyFilename.Error(), http.StatusBadRequest)
		return
	}

	// We can skip checking of invalid characters, because Go will return an error
	err := storage.RenameFile(oldName, newName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// PUT /api/tags?file=filename&tags=first,second,third
//
// Response: -
//
func changeTags(w http.ResponseWriter, r *http.Request) {
	var (
		filename = r.FormValue("file")
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

	err := storage.ChangeTags(filename, tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
