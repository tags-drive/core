package web

import (
	"encoding/json"
	"net/http"

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

// POST /api/tags?tag=newtag
//
// Response: -
//
func addTag(w http.ResponseWriter, r *http.Request) {
	tagName := r.FormValue("tag")
	if tagName == "" {
		http.Error(w, "tag is empty", http.StatusBadRequest)
		return
	}

	tags.AddTag(tags.Tag{Name: tagName, Color: tags.DefaultColor})
	w.WriteHeader(http.StatusCreated)
}

// PUT /api/tag?tag=tagname&color=new-color&new-name=new-name
// new-color shouldn't contain '#'
//
// Response: -
//
func changeTag(w http.ResponseWriter, r *http.Request) {
	var (
		tagName  = r.FormValue("tag")
		newName  = r.FormValue("new-name")
		newColor = r.FormValue("new-color")
	)

	if tagName == "" {
		http.Error(w, "tag is empty", http.StatusBadRequest)
		return
	}

	if !tags.Check(tagName) {
		http.Error(w, tags.ErrTagIsNotExist.Error(), http.StatusBadRequest)
		return
	}

	tags.Change(tagName, newName, newColor)
}

// DELETE /api/tags?tag=tagname
//
// Response: -
//
func deleteTag(w http.ResponseWriter, r *http.Request) {
	tagName := r.FormValue("tag")
	if tagName == "" {
		http.Error(w, "tag is empty", http.StatusBadRequest)
		return
	}
	tags.DeleteTag(tagName)
}
