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
// Response: map with all tags
//
func returnTags(w http.ResponseWriter, r *http.Request) {
	allTags := tags.GetAllTags()

	if params.Debug {
		setDebugHeaders(w, r)
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(allTags)
}

// POST /api/tags?name=tag-name&color=tag-color
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
	if params.Debug {
		setDebugHeaders(w, r)
	}
	w.WriteHeader(http.StatusCreated)
}

// PUT /api/tags?id=tagID&name=new-name&color=new-color
// new-color shouldn't contain '#'
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

	if params.Debug {
		setDebugHeaders(w, r)
	}
}

// DELETE /api/tags?id=tadID
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

	if params.Debug {
		setDebugHeaders(w, r)
	}
}
