package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage"
)

// GET /api/tags
//
// Params: -
//
// Response: json map
//
func returnTags(w http.ResponseWriter, r *http.Request) {
	allTags := storage.Tags.GetAll()

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
		processError(w, "tag is empty", http.StatusBadRequest)
		return
	}
	if tagColor == "" {
		// Default color is white
		tagColor = "#ffffff"
	}

	storage.Tags.Add(tagName, tagColor)
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
		processError(w, "tag id isn't valid", http.StatusBadRequest)
		return
	}

	storage.Tags.Change(id, newName, newColor)
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
		processError(w, "tag id isn't valid", http.StatusBadRequest)
		return
	}
	storage.Tags.Delete(id)
	// Delete refs to tag
	storage.Files.DeleteTagFromFiles(id)
}
