package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/tags-drive/core/internal/storage/share_tokens"
)

// GET /api/share/tokens
//
// Params: -
//
// Response: json map with tokens and ids of shared files
//
func (s Server) getAllShareTokens(w http.ResponseWriter, r *http.Request) {
	allTokens := s.shareService.GetAllTokens()

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(allTokens)
}

// GET /api/share/token/{token}
//
// Params:
//   - token: share token
//
// Response: json array with ids of shared files
//
func (s Server) getFilesSharedByToken(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	if token == "" {
		s.processError(w, "share token can't be empty", http.StatusBadRequest)
		return
	}

	sharedFiles, err := s.shareService.GetFilesIDs(token)
	if err != nil {
		if err == share.ErrInvalidToken {
			s.processError(w, "invalid share token", http.StatusBadRequest)
		} else {
			s.processError(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if s.config.Debug {
		enc.SetIndent("", "  ")
	}

	enc.Encode(sharedFiles)
}

// POST /api/share/token
//
// Params:
//   - ids: list of ids of files to share separated by commas (example: "1,2,3")
//
// Response: { "token": "created token" }
//
func (s Server) createShareToken(w http.ResponseWriter, r *http.Request) {
	ids := func() []int {
		t := r.FormValue("ids")
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

	var goodIDs []int
	for _, id := range ids {
		if s.fileStorage.CheckFile(id) {
			goodIDs = append(goodIDs, id)
		}
	}

	token := s.shareService.CreateToken(goodIDs)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, `{"token":"%s"}`, token)
}

// DELETE /api/share/token/{token}
//
// Params:
//   - token: share token
//
// Response: -
//
func (s Server) deleteShareToken(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	if token == "" {
		s.processError(w, "share token can't be empty", http.StatusBadRequest)
		return
	}

	s.shareService.DeleteToken(token)
}
