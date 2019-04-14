package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) addDefaultRoutes(router *mux.Router) {
	routes := []route{
		// Pages
		{"/", "GET", s.index, true},
		{"/mobile", "GET", s.mobile, true},
		{"/login", "GET", s.login, false},

		// Auth
		{"/api/login", "POST", s.authentication, false},
		{"/api/logout", "POST", s.logout, true},
		// deprecated
		{"/login", "POST", s.authentication, false},
		{"/logout", "POST", s.logout, true},

		// Files
		{"/api/file/{id:\\d+}", "GET", s.returnSingleFile, false},
		{"/api/files", "GET", s.returnFiles, true},
		{"/api/files/recent", "GET", s.returnRecentFiles, true},
		{"/api/files/download", "GET", s.downloadFiles, true},
		{"/api/files", "POST", s.upload, true},
		// change file info
		{"/api/file/{id:\\d+}/name", "PUT", s.changeFilename, true},
		{"/api/file/{id:\\d+}/tags", "PUT", s.changeFileTags, true},
		{"/api/file/{id:\\d+}/description", "PUT", s.changeFileDescription, true},
		// bulk tags changing
		{"/api/files/tags", "POST", s.addTagsToFiles, true},
		{"/api/files/tags", "DELETE", s.removeTagsFromFiles, true},
		// remove or recover files
		{"/api/files", "DELETE", s.deleteFile, true},
		{"/api/files/recover", "POST", s.recoverFile, true},

		// Tags
		{"/api/tags", "GET", s.returnTags, true},
		{"/api/tags", "POST", s.addTag, true},
		{"/api/tag/{id:\\d+}", "PUT", s.changeTag, true},
		{"/api/tags", "DELETE", s.deleteTag, true},
	}

	for _, r := range routes {
		var handler http.Handler = r.handler
		if r.needAuth {
			handler = s.authMiddleware(r.handler)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}
}

func (s *Server) addDebugRoutes(router *mux.Router) {
	routes := []route{
		{"/login", "OPTIONS", setDebugHeaders, false},
		{"/logout", "OPTIONS", setDebugHeaders, false},
		{"/api/file/{id:\\d+}", "OPTIONS", setDebugHeaders, false},
		{"/api/files", "OPTIONS", setDebugHeaders, false},
		{"/api/files/tags", "OPTIONS", setDebugHeaders, false},
		{"/api/files/recover", "OPTIONS", setDebugHeaders, false},
		{"/api/file/{id:\\d+}/tags", "OPTIONS", setDebugHeaders, false},
		{"/api/file/{id:\\d+}/name", "OPTIONS", setDebugHeaders, false},
		{"/api/file/{id:\\d+}/description", "OPTIONS", setDebugHeaders, false},
		{"/api/tags", "OPTIONS", setDebugHeaders, false},
		{"/api/tag/{id:\\d+}", "OPTIONS", setDebugHeaders, false},
	}

	for _, r := range routes {
		var handler http.Handler = r.handler
		if r.needAuth {
			handler = s.authMiddleware(r.handler)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}
}
