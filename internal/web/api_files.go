package web

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/files/logical-parser"
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

// GET /api/files
//
// Params:
//   - expr: logical expression
//   - search: text for search
//   - sort: name | size | time
//   - order: asc | desc
//
// Response: json array
//
func returnFiles(w http.ResponseWriter, r *http.Request) {
	var (
		order  = getParam("asc", r.FormValue("order"), "asc", "desc")
		expr   = r.FormValue("expr")
		search = r.FormValue("search")

		sortMode = files.SortByNameAsc
	)

	// Set sortMode
	// Can skip default
	switch r.FormValue("sort") {
	case "name":
		if order == "asc" {
			sortMode = files.SortByNameAsc
		} else {
			sortMode = files.SortByNameDesc
		}
	case "size":
		if order == "asc" {
			sortMode = files.SortBySizeAsc
		} else {
			sortMode = files.SortBySizeDecs
		}
	case "time":
		if order == "asc" {
			sortMode = files.SortByTimeAsc
		} else {
			sortMode = files.SortByTimeDesc
		}
	}

	parsedExpr, err := parser.Parse(expr)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(files.Get(parsedExpr, sortMode, search))
}

// GET /api/files/recent
//
// Params:
//   - number: number of returned files (5 is a default value)
//
// Response: same as `GET /api/files`
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

	files := files.GetRecent(number)

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(files)
}

// GET /api/files/download
//
// Params:
//   - file file for downloading
//     (to download multiple files at a time, use `file` several times: `file=123.jp  file=hello.png`)
//
// Response: zip archive
//
func downloadFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filenames := r.Form["file"]
	if len(filenames) == 0 {
		Error(w, "list of files can't be empty", http.StatusBadRequest)
		return
	}

	body, err := files.ArchiveFiles(filenames)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	if _, err := io.Copy(w, body); err != nil {
		log.Errorf("can't copy zip file to response body: %s\n", err)
	}
}

// POST /api/files
//
// Body must be "multipart/form-data"
//
// Response: json array
//
func upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		switch err {
		case http.ErrNotMultipart:
			Error(w, err.Error(), http.StatusBadRequest)
		default:
			Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	respChan := make(chan multiplyResponse, 50)
	wg := new(sync.WaitGroup)

	for _, f := range r.MultipartForm.File["files"] {
		wg.Add(1)
		go func(header *multipart.FileHeader) {
			defer wg.Done()
			err := files.UploadFile(header)
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
		// Log an error
		if resp.IsError {
			log.Errorf("Can't load a file %s: %s\n", resp.Filename, resp.Error)
		}
		response = append(response, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(response)
}

// POST /api/files/recover
//
// Params:
//   - file: file for recovering
//     (to recover multiple files at a time, use `file` several times:`file=123.jpg&file=hello.png`)
//
// Response: -
//
func recoverFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filenames := r.Form["file"]
	if len(filenames) == 0 {
		Error(w, "list of files for recovering can't be empty", http.StatusBadRequest)
		return
	}

	for _, f := range filenames {
		files.RecoverFile(f)
	}
}

// PUT /api/files/name
//
// Params:
//   - file: old filename
//   - new-name: new filename
//
//  Response: -
//
func changeFilename(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("file")
	if filename == "" {
		Error(w, ErrEmptyFilename.Error(), http.StatusBadRequest)
		return
	}

	newName := r.FormValue("new-name")
	if newName == "" {
		Error(w, "new-name param can't be empty", http.StatusBadRequest)
		return
	}

	// We can skip checking of invalid characters, because Go will return an error
	err := files.RenameFile(filename, newName)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /api/files/tags
//
// Params:
//   - file: filename
//   - tags: updated list of tags, separated by comma (`tags=1,2,3`)
//
// Response: -
//
func changeFileTags(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("file")
	if filename == "" {
		Error(w, ErrEmptyFilename.Error(), http.StatusBadRequest)
		return
	}

	tags := func() []int {
		t := r.FormValue("tags")
		if t == "" {
			return []int{}
		}

		res := []int{}
		for _, s := range strings.Split(t, ",") {
			if id, err := strconv.Atoi(s); err == nil {
				res = append(res, id)
			}
		}
		return res
	}()

	err := files.ChangeTags(filename, tags)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /api/files/description
//
// Params:
//   - file: filename
//   - description: updated description
//
// Response: -
//
func changeFileDescription(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("file")
	if filename == "" {
		Error(w, ErrEmptyFilename.Error(), http.StatusBadRequest)
		return
	}

	newDescription := r.FormValue("description")
	err := files.ChangeDescription(filename, newDescription)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DELETE /api/files
//
// Params:
//   - file: file for deleting
//     (to delete multiplefiles at a time, use `file` several times:`file=123.jpg&file=hello.png`)
//   - force: should file be deleted right now
//     (if it isn't empty, file will be deleted right now)
//
// Response: json array
//
func deleteFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filenames := r.Form["file"]
	if len(filenames) == 0 {
		Error(w, "list of files for deleting can't be empty", http.StatusBadRequest)
		return
	}

	force := func() bool {
		return r.FormValue("force") != ""
	}()

	respChan := make(chan multiplyResponse, 50)
	wg := new(sync.WaitGroup)

	for _, filename := range filenames {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()

			if !force {
				err := files.DeleteFile(f)
				if err != nil {
					respChan <- multiplyResponse{Filename: f, IsError: true, Error: err.Error()}
				} else {
					respChan <- multiplyResponse{Filename: f, Status: "added into trash"}
				}
			} else {
				err := files.DeleteFileForce(f)
				if err != nil {
					respChan <- multiplyResponse{Filename: f, IsError: true, Error: err.Error()}
				} else {
					respChan <- multiplyResponse{Filename: f, Status: "deleted"}
				}
			}
		}(filename)
	}

	wg.Wait()
	close(respChan)

	var response []multiplyResponse
	for resp := range respChan {
		response = append(response, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(response)
}
