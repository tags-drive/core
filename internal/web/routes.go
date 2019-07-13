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

	needAuth           bool // true by default
	shareable          bool // false by default
	openGraphSupported bool // false by default
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

func (r *route) supportOpenGraph() *route {
	r.openGraphSupported = true
	return r
}

func (s *Server) addDefaultRoutes(router *mux.Router) {
	const (
		GET    = http.MethodGet
		POST   = http.MethodPost
		PUT    = http.MethodPut
		DELETE = http.MethodDelete
	)

	routes := []*route{
		// Pages
		newRoute("/", GET, s.index),
		newRoute("/mobile", GET, s.mobile),
		newRoute("/share", GET, s.share).enableShare(),
		newRoute("/login", GET, s.login).disableAuth(),
		newRoute("/version", GET, s.backendVersion).disableAuth(),

		// Auth
		newRoute("/api/user", GET, s.checkUser).disableAuth(),
		newRoute("/api/login", POST, s.authentication).disableAuth(),
		newRoute("/api/logout", POST, s.logout),
		// deprecated
		newRoute("/login", POST, s.authentication).disableAuth(),
		newRoute("/logout", POST, s.logout),

		// Files
		newRoute("/api/file/{id:\\d+}", GET, s.returnSingleFile).enableShare(),
		newRoute("/api/files", GET, s.returnFiles).enableShare(),
		newRoute("/api/files/recent", GET, s.returnRecentFiles),
		newRoute("/api/files/download", GET, s.downloadFiles).enableShare(),
		// upload new files
		newRoute("/api/files", POST, s.upload),
		// change file info
		newRoute("/api/file/{id:\\d+}/name", PUT, s.changeFilename),
		newRoute("/api/file/{id:\\d+}/tags", PUT, s.changeFileTags),
		newRoute("/api/file/{id:\\d+}/description", PUT, s.changeFileDescription),
		// bulk tags changing
		newRoute("/api/files/tags", POST, s.addTagsToFiles),
		newRoute("/api/files/tags", DELETE, s.removeTagsFromFiles),
		// remove or recover files
		newRoute("/api/files", DELETE, s.deleteFile),
		newRoute("/api/files/recover", POST, s.recoverFile),

		// Tags
		newRoute("/api/tags", GET, s.returnTags).enableShare(),
		newRoute("/api/tags", POST, s.addTag),
		newRoute("/api/tag/{id:\\d+}", PUT, s.changeTag),
		newRoute("/api/tags", DELETE, s.deleteTag),

		// Share
		newRoute("/api/share/tokens", GET, s.getAllShareTokens),
		newRoute("/api/share/token/{token}", GET, s.getFilesSharedByToken),
		newRoute("/api/share/token", POST, s.createShareToken),
		newRoute("/api/share/token/{token}", DELETE, s.deleteShareToken),
	}

	for _, r := range routes {
		var handler http.Handler = r.handler

		if r.needAuth {
			handler = s.authMiddleware(handler, r.shareable)
		}

		// s.openGraphMiddleware must be the last Middleware here to be the first Middleware to be executed
		// (s.openGraphMiddleware has the greatest priority)
		if r.openGraphSupported {
			handler = s.openGraphMiddleware(handler)
		}

		router.Path(r.path).Methods(r.methods).Handler(handler)
	}
}

func (s *Server) addDebugRoutes(router *mux.Router) {
	const OPTIONS = http.MethodOptions

	routes := []route{
		{path: "/login", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/logout", methods: OPTIONS, handler: setDebugHeaders},
		//
		{path: "/api/file/{id:\\d+}", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/files", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/files/tags", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/files/recover", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/file/{id:\\d+}/tags", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/file/{id:\\d+}/name", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/file/{id:\\d+}/description", methods: OPTIONS, handler: setDebugHeaders},
		//
		{path: "/api/tags", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/tag/{id:\\d+}", methods: OPTIONS, handler: setDebugHeaders},
		//
		{path: "/api/share/token", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/share/tokens", methods: OPTIONS, handler: setDebugHeaders},
		{path: "/api/share/token/{token}", methods: OPTIONS, handler: setDebugHeaders},
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
