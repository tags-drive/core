package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/tags-drive/core/internal/params"
)

// GET /api/tags
//
// Params: -
//
// Response: json map
//
func (s Server) returnTags(w http.ResponseWriter, r *http.Request) {
	allTags := s.tagStorage.GetAll()

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(allTags)
}

// POST /api/tags
//
// Params:
//   - name: name of a new tag
//   - color: color of a new tag (`#ffffff` by default)
//
// Response: -
//
func (s Server) addTag(w http.ResponseWriter, r *http.Request) {
	tagName := r.FormValue("name")
	tagColor := r.FormValue("color")
	if tagName == "" {
		s.processError(w, "tag is empty", http.StatusBadRequest)
		return
	}
	if tagColor == "" {
		// Default color is white
		tagColor = "#ffffff"
	}

	s.tagStorage.Add(tagName, tagColor)
	w.WriteHeader(http.StatusCreated)
}

// PUT /api/tag/{id}
//
// Params:
//   - id: id of a tag
//   - name: new name of a tag (can be empty)
//   - color: new color of a tag (can be empty)
//
// Response: -
//
func (s Server) changeTag(w http.ResponseWriter, r *http.Request) {
	var (
		tagID    = mux.Vars(r)["id"]
		newName  = r.FormValue("name")
		newColor = r.FormValue("color")
	)

	var (
		id  int
		err error
	)
	if id, err = strconv.Atoi(tagID); err != nil {
		s.processError(w, "tag id isn't valid", http.StatusBadRequest)
		return
	}

	s.tagStorage.Change(id, newName, newColor)
}

// DELETE /api/tags
//
// Params:
//   - id: id of a tag (one tag at a time)
//
// Response: -
//
func (s Server) deleteTag(w http.ResponseWriter, r *http.Request) {
	tagID := r.FormValue("id")
	var (
		id  int
		err error
	)
	if id, err = strconv.Atoi(tagID); err != nil {
		s.processError(w, "tag id isn't valid", http.StatusBadRequest)
		return
	}
	s.tagStorage.Delete(id)
	// Delete refs to tag
	s.fileStorage.DeleteTagFromFiles(id)
}
