package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/tags-drive/core/internal/storage/share_tokens"
	"github.com/tags-drive/core/internal/storage/tags"
)

// GET /api/tags
//
// Params:
//   - shareToken (optional): share token
//
// Response: json map
//
func (s Server) returnTags(w http.ResponseWriter, r *http.Request) {
	state, ok := getRequestState(r.Context())
	if !ok {
		s.processError(w, "can't obtain request state", http.StatusInternalServerError)
		return
	}

	allTags := s.tagStorage.GetAll()

	if state.shareAccess {
		// Have to filter tags
		var err error
		allTags, err = s.shareStorage.FilterTags(state.shareToken, allTags)
		if err != nil {
			if err == share.ErrInvalidToken {
				// Just in case
				s.processError(w, err.Error(), http.StatusBadRequest)
			} else {
				s.processError(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(allTags)
}

// POST /api/tags
//
// Params:
//   - name: name of a new tag
//   - color: color of a new tag (`#ffffff` by default)
//   - group: group of a new tag (empty by default)
//
// Response: -
//
func (s Server) addTag(w http.ResponseWriter, r *http.Request) {
	tagName := r.FormValue("name")
	tagColor := r.FormValue("color")
	tagGroup := r.FormValue("group")

	if tagName == "" {
		s.processError(w, "tag is empty", http.StatusBadRequest)
		return
	}
	if tagColor == "" {
		// Default color is white
		tagColor = "#ffffff"
	}

	s.tagStorage.Add(tagName, tagColor, tagGroup)
	w.WriteHeader(http.StatusCreated)
}

// PUT /api/tag/{id}
//
// Params:
//   - id: id of a tag
//   - name: new name of a tag (can be empty)
//   - color: new color of a tag (can be empty)
//   - group: new group of a tag (can be empty)
//
// Response: update tag
//
func (s Server) changeTag(w http.ResponseWriter, r *http.Request) {
	var (
		tagID    = mux.Vars(r)["id"]
		newName  = r.FormValue("name")
		newColor = r.FormValue("color")
	)

	id, err := strconv.Atoi(tagID)
	if err != nil {
		s.processError(w, "tag id isn't valid", http.StatusBadRequest)
		return
	}

	// Check id
	if !s.tagStorage.Check(id) {
		s.processError(w, "tag with id "+tagID+" doesn't exist", http.StatusBadRequest)
		return
	}

	var updatedTag tags.Tag

	if newName != "" || newColor != "" {
		// name or color was passed, we should update tag
		updatedTag, err = s.tagStorage.UpdateTag(id, newName, newColor)
		if err != nil {
			s.processError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if values, ok := r.Form["group"]; ok && len(values) > 0 {
		// group was passed
		newGroup := values[0]
		updatedTag, err = s.tagStorage.UpdateGroup(id, newGroup)
		if err != nil {
			s.processError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(updatedTag)
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
	s.fileStorage.RemoveTagFromAllFiles(id)
}
