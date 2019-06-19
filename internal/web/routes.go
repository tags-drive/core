package web

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

type route struct {
	path    string
	methods string
	handler http.HandlerFunc

	needAuth  bool // true by default
	shareable bool // false by default
}

func newRoute(path, method string, handler http.HandlerFunc) *route {
	return &route{
		path:      path,
		methods:   method,
		handler:   handler,
		needAuth:  true,
		shareable: false,
	}
}

func (r *route) enableShare() *route {
	r.shareable = true
	return r
}

func (r *route) disableAuth() *route {
	r.needAuth = false
	return r
}

func (s *Server) addDefaultRoutes(router *mux.Router) {
	routes := []*route{
		// Pages
		newRoute("/", "GET", s.index),
		newRoute("/mobile", "GET", s.mobile),
		newRoute("/share", "GET", s.share).enableShare(),
		newRoute("/login", "GET", s.login).disableAuth(),
		newRoute("/version", "GET", s.backendVersion).disableAuth(),

		// Auth
		newRoute("/api/user", "GET", s.checkUser).disableAuth(),
		newRoute("/api/login", "POST", s.authentication).disableAuth(),
		newRoute("/api/logout", "POST", s.logout),
		// deprecated
		newRoute("/login", "POST", s.authentication).disableAuth(),
		newRoute("/logout", "POST", s.logout),

		// Files
		newRoute("/api/file/{id:\\d+}", "GET", s.returnSingleFile).enableShare(),
		newRoute("/api/files", "GET", s.returnFiles).enableShare(),
		newRoute("/api/files/recent", "GET", s.returnRecentFiles),
		newRoute("/api/files/download", "GET", s.downloadFiles).enableShare(),
		// upload new files
		newRoute("/api/files", "POST", s.upload),
		// change file info
		newRoute("/api/file/{id:\\d+}/name", "PUT", s.changeFilename),
		newRoute("/api/file/{id:\\d+}/tags", "PUT", s.changeFileTags),
		newRoute("/api/file/{id:\\d+}/description", "PUT", s.changeFileDescription),
		// bulk tags changing
		newRoute("/api/files/tags", "POST", s.addTagsToFiles),
		newRoute("/api/files/tags", "DELETE", s.removeTagsFromFiles),
		// remove or recover files
		newRoute("/api/files", "DELETE", s.deleteFile),
		newRoute("/api/files/recover", "POST", s.recoverFile),

		// Tags
		newRoute("/api/tags", "GET", s.returnTags).enableShare(),
		newRoute("/api/tags", "POST", s.addTag),
		newRoute("/api/tag/{id:\\d+}", "PUT", s.changeTag),
		newRoute("/api/tags", "DELETE", s.deleteTag),

		// Share
		newRoute("/api/share/tokens", "GET", s.getAllShareTokens),
		newRoute("/api/share/token/{token}", "GET", s.getFilesSharedByToken),
		newRoute("/api/share/token", "POST", s.createShareToken),
		newRoute("/api/share/token/{token}", "DELETE", s.deleteShareToken),
	}

	for _, r := range routes {
		var handler http.Handler = r.handler
		if r.needAuth {
			handler = s.authMiddleware(r.handler, r.shareable)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}
}

func (s *Server) addDebugRoutes(router *mux.Router) {
	routes := []route{
		{"/login", "OPTIONS", setDebugHeaders, false, false},
		{"/logout", "OPTIONS", setDebugHeaders, false, false},
		//
		{"/api/file/{id:\\d+}", "OPTIONS", setDebugHeaders, false, false},
		{"/api/files", "OPTIONS", setDebugHeaders, false, false},
		{"/api/files/tags", "OPTIONS", setDebugHeaders, false, false},
		{"/api/files/recover", "OPTIONS", setDebugHeaders, false, false},
		{"/api/file/{id:\\d+}/tags", "OPTIONS", setDebugHeaders, false, false},
		{"/api/file/{id:\\d+}/name", "OPTIONS", setDebugHeaders, false, false},
		{"/api/file/{id:\\d+}/description", "OPTIONS", setDebugHeaders, false, false},
		//
		{"/api/tags", "OPTIONS", setDebugHeaders, false, false},
		{"/api/tag/{id:\\d+}", "OPTIONS", setDebugHeaders, false, false},
		//
		{"/api/share/token", "OPTIONS", setDebugHeaders, false, false},
		{"/api/share/tokens", "OPTIONS", setDebugHeaders, false, false},
		{"/api/share/token/{token}", "OPTIONS", setDebugHeaders, false, false},
	}

	for _, r := range routes {
		var handler http.Handler = r.handler
		if r.needAuth {
			handler = s.authMiddleware(r.handler, r.shareable)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}
}

func (s *Server) addPprofRoutes(router *mux.Router) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/debug/pprof/", pprof.Index)
	serveMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	serveMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	serveMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	serveMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.PathPrefix("/debug/pprof/").Handler(serveMux)
}
