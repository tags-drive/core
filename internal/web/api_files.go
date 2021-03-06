package web

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	filesPck "github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
)

const (
	// It is trade-off between memory and I/O
	// If maxSize == 50MB, program takes too much memory
	// If maxSize == 2MB, there're too many I/O-operations
	maxSize = 10 << 20 // 10MB

	maxThreadsInPool = 3
)

// multiplyResponse is used as a response by POST /api/files and DELETE /api/files
type multiplyResponse struct {
	Filename string `json:"filename"`
	IsError  bool   `json:"isError"`
	Error    string `json:"error"`
	Status   string `json:"status"` // Status isn't empty when IsError == false
}

// GET /api/file/{id}
//
// Params:
//   - id: id of a file
//   - shareToken (optional): share token
//
// Response: json object
//
func (s Server) returnSingleFile(w http.ResponseWriter, r *http.Request) {
	state, ok := getRequestState(r.Context())
	if !ok {
		s.processError(w, "can't obtain request state", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		s.processError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if state.shareAccess {
		if !s.shareService.CheckFile(state.shareToken, id) {
			s.processError(w, "share token doesn't grant access to this file", http.StatusForbidden)
			return
		}
	}

	file, err := s.fileStorage.GetFile(id)
	if err != nil {
		if err == filesPck.ErrFileIsNotExist {
			s.processError(w, "file doesn't exist", http.StatusNotFound)
			return
		}

		s.processError(w, "can't get file", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(file)
}

// GET /api/files
//
// Params:
//   - expr: logical expression
//   - search: text for search
//   - regexp: is search a regular expression (it is true when regexp != "")
//   - sort: name | size | time
//   - order: asc | desc
//   - offset: lower bound [offset:]
//   - count: number of returned files ([offset:offset+count]). If count == 0, all files will be returned. Default is 0
//   - shareToken (optional): share token
//
// Response: json array
//
func (s Server) returnFiles(w http.ResponseWriter, r *http.Request) {
	state, ok := getRequestState(r.Context())
	if !ok {
		s.processError(w, "can't obtain request state", http.StatusInternalServerError)
		return
	}

	getSortMode := func(sortType, sortOrder string) filesPck.FilesSortMode {
		// Set default values if needed
		sortType = getParam(sortType, "name", []string{"name", "size", "time"})
		sortOrder = getParam(sortOrder, "asc", []string{"asc", "desc"})

		switch sortType {
		case "name":
			if sortOrder == "asc" {
				return filesPck.SortByNameAsc
			}
			return filesPck.SortByNameDesc
		case "size":
			if sortOrder == "asc" {
				return filesPck.SortBySizeAsc
			}
			return filesPck.SortBySizeDecs
		case "time":
			if sortOrder == "asc" {
				return filesPck.SortByTimeAsc
			}
			return filesPck.SortByTimeDesc
		default:
			return filesPck.SortByNameAsc
		}
	}

	customAtoi := func(value string, defaultValue int) int {
		if value == "" {
			return defaultValue
		}

		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return defaultValue
		}

		return n
	}

	cnf := filesPck.GetFilesConfig{
		Expr:     r.FormValue("expr"),
		Search:   r.FormValue("search"),
		IsRegexp: r.FormValue("regexp") != "",
		SortMode: getSortMode(r.FormValue("sort"), r.FormValue("order")),
		Offset:   customAtoi(r.FormValue("offset"), 0),
		Count:    customAtoi(r.FormValue("count"), 0),
		Filter:   nil,
	}

	// Check if a regexp is valid
	if cnf.IsRegexp {
		if _, err := regexp.Compile(cnf.Search); err != nil {
			s.processError(w, "invalid regular expression", http.StatusBadRequest)
			return
		}
	}

	// Add a filter if needed
	if state.shareAccess {
		cnf.Filter = filesPck.FilterFilesFunction(func(files []filesPck.File) ([]filesPck.File, error) {
			return s.shareService.FilterFiles(state.shareToken, files)
		})
	}

	files, err := s.fileStorage.Get(cnf)
	if err != nil {
		switch err {
		case filesPck.ErrOffsetOutOfBounds:
			s.processError(w, "offset is out of bounds", http.StatusNoContent, err)
		case aggregation.ErrBadSyntax:
			s.processError(w, "bad syntax of logical expression", http.StatusBadRequest, err)
		default:
			s.processError(w, "can't get files", http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(files)
}

func getParam(passedVal, defaultVal string, validOptions []string) string {
	for _, opt := range validOptions {
		if passedVal == opt {
			return passedVal
		}
	}

	return defaultVal
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
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(files)
}

// GET /api/files/download
//
// Params:
//   - ids: list of ids of files for downloading separated by comma `ids=1,2,54,9`
//   - shareToken (optional): share token
//
// Response: zip archive
//
func (s Server) downloadFiles(w http.ResponseWriter, r *http.Request) {
	state, ok := getRequestState(r.Context())
	if !ok {
		s.processError(w, "can't obtain request state", http.StatusInternalServerError)
		return
	}

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

	if state.shareAccess {
		// Have to filter ids
		goodIDs := make([]int, 0, len(ids))
		for _, id := range ids {
			if s.shareService.CheckFile(state.shareToken, id) {
				goodIDs = append(goodIDs, id)
			}
		}
		ids = goodIDs
	}

	body, err := s.fileStorage.Archive(ids)
	if err != nil {
		s.processError(w, "can't archive files", http.StatusInternalServerError, err)
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
			s.processError(w, "invalid form type", http.StatusBadRequest, err)
		default:
			s.processError(w, "can't parse request form", http.StatusInternalServerError, err)
		}
		return
	}

	responses := make([]multiplyResponse, 0, len(r.MultipartForm.File["files"]))
	responsesReady := make(chan struct{})
	responsesChan := make(chan multiplyResponse, 50)

	headersChan := make(chan interface{}, 5)
	// Fill headersChan
	go func() {
		for i := range r.MultipartForm.File["files"] {
			headersChan <- r.MultipartForm.File["files"][i]
		}
		close(headersChan)
	}()

	// Fill responsesChan
	go func() {
		for r := range responsesChan {
			responses = append(responses, r)
		}
		close(responsesReady)
	}()

	runPool(maxThreadsInPool, headersChan, func(data <-chan interface{}) {
		for d := range data {
			header, ok := d.(*multipart.FileHeader)
			if !ok {
				continue
			}

			err := s.fileStorage.Upload(header, tags)
			var resp multiplyResponse
			if err != nil {
				resp = multiplyResponse{
					Filename: header.Filename,
					IsError:  true,
					Error:    err.Error(),
				}
				s.logger.Errorf("can't load a file %s: %s\n", header.Filename, err)
			} else {
				resp = multiplyResponse{Filename: header.Filename, Status: "uploaded"}
			}

			responsesChan <- resp
		}
	})
	close(responsesChan)

	<-responsesReady

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(responses)
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

	idsChan := make(chan interface{}, 5)
	go func() {
		for i := range ids {
			idsChan <- ids[i]
		}
		close(idsChan)
	}()

	runPool(maxThreadsInPool, idsChan, func(data <-chan interface{}) {
		for d := range data {
			id, ok := d.(int)
			if !ok {
				continue
			}
			s.fileStorage.Recover(id)
		}
	})
}

// PUT /api/file/{id}/name
//
// Params:
//   - id: file id
//   - new-name: new filename
//
//  Response: updated file
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
	updatedFile, err := s.fileStorage.Rename(id, newName)
	if err != nil {
		s.processError(w, "can't rename file", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(updatedFile)
}

// PUT /api/file/{id}/tags
//
// Params:
//   - id: file id
//   - tags: updated list of tags, separated by comma (`tags=1,2,3`)
//
// Response: updated file
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

	updatedFile, err := s.fileStorage.ChangeTags(fileID, goodTags)
	if err != nil {
		s.processError(w, "can't change file tags", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(updatedFile)
}

// PUT /api/file/{id}/description
//
// Params:
//   - id: file id
//   - description: updated description
//
// Response: updated file
//
func (s Server) changeFileDescription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		s.processError(w, "bad id syntax", http.StatusBadRequest)
		return
	}
	newDescription := r.FormValue("description")

	updatedFile, err := s.fileStorage.ChangeDescription(id, newDescription)
	if err != nil {
		s.processError(w, "can't change file description", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(updatedFile)
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

	responses := make([]multiplyResponse, 0, len(ids))
	responsesChan := make(chan multiplyResponse, 50)
	responsesReady := make(chan struct{})

	filesIDsChan := make(chan interface{}, 5)
	// Fill filesIDsChan
	go func() {
		for _, id := range ids {
			filesIDsChan <- id
		}
		close(filesIDsChan)
	}()

	// Fill responsesChan
	go func() {
		for r := range responsesChan {
			responses = append(responses, r)
		}
		close(responsesReady)
	}()

	// Used in a worker function
	var (
		deleteFunc = s.fileStorage.Delete
		// We will use status if deleteFunc returns nil error
		respStatus = "added into trash"
	)

	if force {
		deleteFunc = s.fileStorage.DeleteForce
		respStatus = "deleted"
	}

	runPool(maxThreadsInPool, filesIDsChan, func(data <-chan interface{}) {
		for d := range data {
			id, ok := d.(int)
			if !ok {
				continue
			}

			// Check file
			file, err := s.fileStorage.GetFile(id)
			if err != nil {
				msg := err.Error()
				if err == filesPck.ErrFileIsNotExist {
					msg = fmt.Sprintf("file with id \"%d\" doesn't exist", id)
				}

				responsesChan <- multiplyResponse{
					Filename: "",
					IsError:  true,
					Error:    msg,
				}

				// We can skip non-existent file
				continue
			}

			var resp multiplyResponse

			// Delete file
			err = deleteFunc(id)
			if err != nil {
				resp = multiplyResponse{
					Filename: file.Filename,
					IsError:  true,
					Error:    err.Error(),
				}
			} else {
				// Use pre-defined var status
				resp = multiplyResponse{
					Filename: file.Filename,
					Status:   respStatus,
				}
			}

			// Delete the file from Share Storage even if deleting is not permanent
			s.shareService.DeleteFile(id)

			responsesChan <- resp
		}
	})
	close(responsesChan)

	<-responsesReady

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(responses)
}
