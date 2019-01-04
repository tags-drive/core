package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ShoshinNikita/log"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage"
	"github.com/tags-drive/core/internal/storage/files"
)

const (
	// It is trade-off between memory and I/O
	// If maxSize == 50MB, program takes too much memory
	// If maxSize == 2MB, there're too many I/O-operations
	maxSize = 10 << 20 // 10MB
)

// multiplyResponse is used as a response by POST /api/files and DELETE /api/files
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
//   - offset: lower bound [offset:]
//
// Response: json array
//
func returnFiles(w http.ResponseWriter, r *http.Request) {
	var (
		order  = getParam("asc", r.FormValue("order"), "asc", "desc")
		expr   = r.FormValue("expr")
		search = r.FormValue("search")

		offset   = 0
		sortMode = files.SortByNameAsc
	)

	// Get offset
	offset = func() int {
		param := r.FormValue("offset")
		if param == "" {
			return 0
		}

		r, err := strconv.ParseInt(param, 10, 0)
		if err != nil || r < 0 {
			return 0
		}

		return int(r)
	}()

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

	files, err := storage.Files.Get(expr, sortMode, search, offset)
	if err != nil {
		if err == storage.ErrBadExpessionSyntax || err == storage.ErrOffsetOutOfBounds {
			Error(w, err.Error(), http.StatusBadRequest)
		} else {
			Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(files)
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

	files := storage.Files.GetRecent(number)

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
//   - ids: list of ids of files for downloading separated by comma `ids=1,2,54,9`
//
// Response: zip archive
//
func downloadFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := func() (res []int) {
		strIDs := r.FormValue("ids")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.ParseInt(strID, 10, 0)
			if err == nil {
				res = append(res, int(id))
			}
		}
		return
	}()

	body, err := storage.Files.Archive(ids)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	if _, err := io.Copy(w, body); err != nil {
		log.Errorf("Can't copy zip file to response body: %s\n", err)
	}
}

// POST /api/files
//
// Body must be "multipart/form-data"
//
// Params:
//   - tags: list of tags, separated by comma (`tags=1,2,3`)
//
// Response: json array
//
func upload(w http.ResponseWriter, r *http.Request) {
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

	var response []multiplyResponse
	for _, header := range r.MultipartForm.File["files"] {
		err := storage.Files.Upload(header, tags)
		if err != nil {
			response = append(response, multiplyResponse{
				Filename: header.Filename,
				IsError:  true,
				Error:    err.Error(),
			})
			log.Errorf("Can't load a file %s: %s\n", header.Filename, err)
		} else {
			response = append(response, multiplyResponse{Filename: header.Filename, Status: "uploaded"})
		}
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
//   - ids: list ids of files for recovering separated by comma `ids=1,2,54,9`
//
// Response: -
//
func recoverFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := func() (res []int) {
		strIDs := r.FormValue("ids")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.ParseInt(strID, 10, 0)
			if err == nil {
				res = append(res, int(id))
			}
		}
		return
	}()

	if len(ids) == 0 {
		Error(w, "list of ids of files for recovering can't be empty", http.StatusBadRequest)
		return
	}

	for _, id := range ids {
		storage.Files.Recover(id)
	}
}

// PUT /api/files/name
//
// Params:
//   - id: file id
//   - new-name: new filename
//
//  Response: -
//
func changeFilename(w http.ResponseWriter, r *http.Request) {
	strID := r.FormValue("id")
	id, err := strconv.ParseInt(strID, 10, 0)
	if err != nil {
		Error(w, "bad id syntax", http.StatusBadRequest)
		return
	}

	newName := r.FormValue("new-name")
	if newName == "" {
		Error(w, "new-name param can't be empty", http.StatusBadRequest)
		return
	}

	// We can skip checking of invalid characters, because Go will return an error
	err = storage.Files.Rename(int(id), newName)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /api/files/tags
//
// Params:
//   - id: file id
//   - tags: updated list of tags, separated by comma (`tags=1,2,3`)
//
// Response: -
//
func changeFileTags(w http.ResponseWriter, r *http.Request) {
	strID := r.FormValue("id")
	id, err := strconv.ParseInt(strID, 10, 0)
	if err != nil {
		Error(w, "bad id syntax", http.StatusBadRequest)
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

	var goodTags []int
	for _, id := range tags {
		if storage.Tags.Check(id) {
			goodTags = append(goodTags, id)
		}
	}

	err = storage.Files.ChangeTags(int(id), goodTags)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /api/files/description
//
// Params:
//   - id: file id
//   - description: updated description
//
// Response: -
//
func changeFileDescription(w http.ResponseWriter, r *http.Request) {
	strID := r.FormValue("id")
	id, err := strconv.ParseInt(strID, 10, 0)
	if err != nil {
		Error(w, "bad id syntax", http.StatusBadRequest)
		return
	}

	newDescription := r.FormValue("description")
	err = storage.Files.ChangeDescription(int(id), newDescription)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DELETE /api/files
//
// Params:
//   - ids: list of ids of files for deleting separated by comma `ids=1,2,54,9`
//   - force: should file be deleted right now
//     (if it isn't empty, file will be deleted right now)
//
// Response: json array
//
func deleteFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := func() (res []int) {
		strIDs := r.FormValue("ids")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.ParseInt(strID, 10, 0)
			if err == nil {
				res = append(res, int(id))
			}
		}
		return
	}()

	force := func() bool {
		return r.FormValue("force") != ""
	}()

	var response []multiplyResponse
	for _, id := range ids {
		file, err := storage.Files.GetFile(id)
		if err != nil {
			msg := err.Error()
			if err == storage.ErrFileIsNotExist {
				msg = fmt.Sprintf("file with id \"%d\" doesn't exist", id)
			}

			response = append(response, multiplyResponse{
				Filename: "",
				IsError:  true,
				Error:    msg,
			})

			// We can skip non-existent file
			continue
		}

		deleteFunc := storage.Files.Delete
		// We will use status if deleteFunc returns nil error
		status := "added into trash"
		if force {
			deleteFunc = storage.Files.DeleteForce
			status = "deleted"
		}

		err = deleteFunc(id)
		if err != nil {
			response = append(response, multiplyResponse{
				Filename: file.Filename,
				IsError:  true,
				Error:    err.Error(),
			})
		} else {
			// Use pre-defined var status
			response = append(response, multiplyResponse{
				Filename: file.Filename,
				Status:   status,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(response)
}
