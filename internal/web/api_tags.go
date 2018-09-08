package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ShoshinNikita/tags-drive/internal/params"
	"github.com/ShoshinNikita/tags-drive/internal/storage/tags"
)

// GET /api/tags
//
// Response: json list of all tags
//
func returnTags(w http.ResponseWriter, r *http.Request) {
	allTags := tags.GetAllTags()

	enc := json.NewEncoder(w)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(allTags)
}

// TODO add param color
// POST /api/tags?tag=newtag
//
// Response: -
//
func addTag(w http.ResponseWriter, r *http.Request) {
	tagName := r.FormValue("tag")
	if tagName == "" {
		Error(w, "tag is empty", http.StatusBadRequest)
		return
	}

	tags.AddTag(tags.Tag{Name: tagName, Color: tags.DefaultColor})
	w.WriteHeader(http.StatusCreated)
}

// PUT /api/tags?id=tagID&new-color=new-color&new-name=new-name
// new-color shouldn't contain '#'
//
// Response: -
//
func changeTag(w http.ResponseWriter, r *http.Request) {
	var (
		tagID    = r.FormValue("id")
		newName  = r.FormValue("new-name")
		newColor = r.FormValue("new-color")
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
}
