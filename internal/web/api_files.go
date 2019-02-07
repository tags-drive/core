package web

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
	filesPck "github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
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

// GET /api/file/{id}
//
// Params:
//   - id: id of a file
//
// Response: json object
//
func (s Server) returnSingleFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		s.processError(w, "invalid id", http.StatusBadRequest)
		return
	}

	file, err := s.fileStorage.GetFile(id)
	if err != nil {
		if err == filesPck.ErrFileIsNotExist {
			s.processError(w, "file doesn't exist", http.StatusNotFound)
			return
		}

		s.processError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(file)
}

// GET /api/files
//
// Params:
//   - expr: logical expression
//   - search: text for search
//   - sort: name | size | time
//   - order: asc | desc
//   - offset: lower bound [offset:]
//   - count: number of returned files ([offset:offset+count]). If count == 0, all files will be returned. Default is 0
//
// Response: json array
//
func (s Server) returnFiles(w http.ResponseWriter, r *http.Request) {
	var (
		order  = getParam("asc", r.FormValue("order"), "asc", "desc")
		expr   = r.FormValue("expr")
		search = r.FormValue("search")

		offset   = 0
		count    = 0
		sortMode = cmd.SortByNameAsc
	)

	// Get offset
	offset = func() int {
		param := r.FormValue("offset")
		if param == "" {
			return 0
		}

		r, err := strconv.Atoi(param)
		if err != nil || r < 0 {
			return 0
		}

		return r
	}()

	// Get offset
	count = func() int {
		param := r.FormValue("count")
		if param == "" {
			return 0
		}

		r, err := strconv.Atoi(param)
		if err != nil || r < 0 {
			return 0
		}

		return r
	}()

	// Set sortMode
	// Can skip default
	switch r.FormValue("sort") {
	case "name":
		if order == "asc" {
			sortMode = cmd.SortByNameAsc
		} else {
			sortMode = cmd.SortByNameDesc
		}
	case "size":
		if order == "asc" {
			sortMode = cmd.SortBySizeAsc
		} else {
			sortMode = cmd.SortBySizeDecs
		}
	case "time":
		if order == "asc" {
			sortMode = cmd.SortByTimeAsc
		} else {
			sortMode = cmd.SortByTimeDesc
		}
	}

	files, err := s.fileStorage.Get(expr, sortMode, search, offset, count)
	if err != nil {
		if err == aggregation.ErrBadSyntax || err == filesPck.ErrOffsetOutOfBounds {
			s.processError(w, err.Error(), http.StatusBadRequest)
		} else {
			s.processError(w, err.Error(), http.StatusInternalServerError)
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
func (s Server) returnRecentFiles(w http.ResponseWriter, r *http.Request) {
	number := func() int {
		s := r.FormValue("number")
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			n = 5
		}
		return int(n)
	}()

	files := s.fileStorage.GetRecent(number)

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
func (s Server) downloadFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := func() (res []int) {
		strIDs := r.FormValue("ids")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.Atoi(strID)
			if err == nil {
				res = append(res, id)
			}
		}
		return
	}()

	body, err := s.fileStorage.Archive(ids)
	if err != nil {
		s.processError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	if _, err := io.Copy(w, body); err != nil {
		s.logger.Errorf("can't copy zip file to response body: %s\n", err)
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
func (s Server) upload(w http.ResponseWriter, r *http.Request) {
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
			s.processError(w, err.Error(), http.StatusBadRequest)
		default:
			s.processError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var response []multiplyResponse
	for _, header := range r.MultipartForm.File["files"] {
		err := s.fileStorage.Upload(header, tags)
		if err != nil {
			response = append(response, multiplyResponse{
				Filename: header.Filename,
				IsError:  true,
				Error:    err.Error(),
			})
			s.logger.Errorf("can't load a file %s: %s\n", header.Filename, err)
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
func (s Server) recoverFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := func() (res []int) {
		strIDs := r.FormValue("ids")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.Atoi(strID)
			if err == nil {
				res = append(res, id)
			}
		}
		return
	}()

	if len(ids) == 0 {
		s.processError(w, "list of ids of files for recovering can't be empty", http.StatusBadRequest)
		return
	}

	for _, id := range ids {
		s.fileStorage.Recover(id)
	}
}

// PUT /api/file/{id}/name
//
// Params:
//   - id: file id
//   - new-name: new filename
//
//  Response: -
//
func (s Server) changeFilename(w http.ResponseWriter, r *http.Request) {
	strID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		s.processError(w, "bad id syntax", http.StatusBadRequest)
		return
	}

	newName := r.FormValue("new-name")
	if newName == "" {
		s.processError(w, "new-name param can't be empty", http.StatusBadRequest)
		return
	}

	// We can skip checking of invalid characters, because Go will return an error
	err = s.fileStorage.Rename(id, newName)
	if err != nil {
		s.processError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /api/file/{id}/tags
//
// Params:
//   - id: file id
//   - tags: updated list of tags, separated by comma (`tags=1,2,3`)
//
// Response: -
//
func (s Server) changeFileTags(w http.ResponseWriter, r *http.Request) {
	strID := mux.Vars(r)["id"]
	fileID, err := strconv.Atoi(strID)
	if err != nil {
		s.processError(w, "bad id syntax", http.StatusBadRequest)
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
		if s.tagStorage.Check(id) {
			goodTags = append(goodTags, id)
		}
	}

	err = s.fileStorage.ChangeTags(fileID, goodTags)
	if err != nil {
		s.processError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /api/file/{id}/description
//
// Params:
//   - id: file id
//   - description: updated description
//
// Response: -
//
func (s Server) changeFileDescription(w http.ResponseWriter, r *http.Request) {
	strID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		s.processError(w, "bad id syntax", http.StatusBadRequest)
		return
	}

	newDescription := r.FormValue("description")
	err = s.fileStorage.ChangeDescription(id, newDescription)
	if err != nil {
		s.processError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /api/files/tags
//
// Params:
//   - files: file ids (list of ids separated by ',')
//   - tags: tags for adding (list of tags ids separated by ',')
//
// Response: -
//
func (s Server) addTagsToFiles(w http.ResponseWriter, r *http.Request) {
	filesIDs := func() (res []int) {
		strIDs := r.FormValue("files")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.Atoi(strID)
			if err == nil {
				res = append(res, id)
			}
		}
		return res
	}()

	tagsIDs := func() (res []int) {
		strIDs := r.FormValue("tags")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.Atoi(strID)
			// Add only valid tags
			if err == nil && s.tagStorage.Check(id) {
				res = append(res, int(id))
			}
		}
		return res
	}()

	s.fileStorage.AddTagsToFiles(filesIDs, tagsIDs)
}

// DELETE /api/files/tags
//
// Params:
//   - files: file ids (list of ids separated by ',')
//   - tags: tags for deleting (list of tags ids separated by ',')
//
// Response: -
//
func (s Server) removeTagsFromFiles(w http.ResponseWriter, r *http.Request) {
	filesIDs := func() (res []int) {
		strIDs := r.FormValue("files")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.Atoi(strID)
			if err == nil {
				res = append(res, id)
			}
		}
		return res
	}()

	tagsIDs := func() (res []int) {
		strIDs := r.FormValue("tags")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.Atoi(strID)
			// Add only valid tags
			if err == nil && s.tagStorage.Check(id) {
				res = append(res, int(id))
			}
		}
		return res
	}()

	s.fileStorage.RemoveTagsFromFiles(filesIDs, tagsIDs)
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
func (s Server) deleteFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := func() (res []int) {
		strIDs := r.FormValue("ids")
		for _, strID := range strings.Split(strIDs, ",") {
			id, err := strconv.Atoi(strID)
			if err == nil {
				res = append(res, id)
			}
		}
		return
	}()

	force := func() bool {
		return r.FormValue("force") != ""
	}()

	var response []multiplyResponse
	for _, id := range ids {
		file, err := s.fileStorage.GetFile(id)
		if err != nil {
			msg := err.Error()
			if err == filesPck.ErrFileIsNotExist {
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

		deleteFunc := s.fileStorage.Delete
		// We will use status if deleteFunc returns nil error
		status := "added into trash"
		if force {
			deleteFunc = s.fileStorage.DeleteForce
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
