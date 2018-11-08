package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/tags"
)

// GET /api/tags
//
// Params: -
//
// Response: json map
//
func returnTags(w http.ResponseWriter, r *http.Request) {
	allTags := tags.GetAllTags()

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
func addTag(w http.ResponseWriter, r *http.Request) {
	tagName := r.FormValue("name")
	tagColor := r.FormValue("color")
	if tagName == "" {
		Error(w, "tag is empty", http.StatusBadRequest)
		return
	}
	if tagColor == "" {
		tagColor = tags.DefaultColor
	}

	tags.AddTag(tags.Tag{Name: tagName, Color: tagColor})
	w.WriteHeader(http.StatusCreated)
}

// PUT /api/tags
//
// Params:
//   - id: id of a tag
//   - name: new name of a tag (can be empty)
//   - color: new color of a tag (can be empty)
//
// Response: -
//
func changeTag(w http.ResponseWriter, r *http.Request) {
	var (
		tagID    = r.FormValue("id")
		newName  = r.FormValue("name")
		newColor = r.FormValue("color")
	)

	var (
		id  int
		err error
	)
	if id, err = strconv.Atoi(tagID); err != nil {
		Error(w, "tag id isn't valid", http.StatusBadRequest)
		return
	}

	tags.Change(id, newName, newColor)
}

// DELETE /api/tags
//
// Params:
//   - id: id of a tag (one tag at a time)
//
// Response: -
//
func deleteTag(w http.ResponseWriter, r *http.Request) {
	tagID := r.FormValue("id")
	var (
		id  int
		err error
	)
	if id, err = strconv.Atoi(tagID); err != nil {
		Error(w, "tag id isn't valid", http.StatusBadRequest)
		return
	}
	tags.DeleteTag(id)
}
