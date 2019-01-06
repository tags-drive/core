// Package web is responsible for serving browser and API requests
package web

import (
	"context"
	"net/http"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/gorilla/mux"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
)

type Server struct {
	fileStorage cmd.FileStorageInterface
	tagStorage  cmd.TagStorageInterface

	logger *log.Logger
}

type route struct {
	path     string
	methods  string
	handler  http.HandlerFunc
	needAuth bool
}

// NewWebServer just creates new Web struct. It doesn't call any Init functions
func NewWebServer(fs cmd.FileStorageInterface, ts cmd.TagStorageInterface, lg *log.Logger) *Server {
	return &Server{
		fileStorage: fs,
		tagStorage:  ts,
		logger:      lg,
	}
}

// Start starts the server. It has to be ran in goroutine
//
// Server stops when ctx.Done()
func (s Server) Start(ctx context.Context) error {
	// Init routes
	var routes = []route{
		{"/", "GET", s.index, true},
		// auth
		{"/login", "GET", s.login, false},
		{"/login", "POST", s.authentication, false},
		{"/logout", "POST", s.logout, true},
		// files
		{"/api/files", "GET", s.returnFiles, true},
		{"/api/files/download", "GET", s.downloadFiles, true},
		{"/api/files", "POST", s.upload, true},
		{"/api/files/recover", "POST", s.recoverFile, true},
		{"/api/files/tags", "PUT", s.changeFileTags, true},
		{"/api/files/name", "PUT", s.changeFilename, true},
		{"/api/files/description", "PUT", s.changeFileDescription, true},
		{"/api/files", "DELETE", s.deleteFile, true},
		{"/api/files/recent", "GET", s.returnRecentFiles, true},
		// tags
		{"/api/tags", "GET", s.returnTags, true},
		{"/api/tags", "POST", s.addTag, true},
		{"/api/tags", "PUT", s.changeTag, true},
		{"/api/tags", "DELETE", s.deleteTag, true},
	}

	var debugRoutes = []route{
		{"/login", "OPTIONS", setDebugHeaders, false},
		{"/logout", "OPTIONS", setDebugHeaders, false},
		{"/api/files", "OPTIONS", setDebugHeaders, false},
		{"/api/files/recover", "OPTIONS", setDebugHeaders, false},
		{"/api/files/tags", "OPTIONS", setDebugHeaders, false},
		{"/api/files/name", "OPTIONS", setDebugHeaders, false},
		{"/api/files/description", "OPTIONS", setDebugHeaders, false},
		{"/api/tags", "OPTIONS", setDebugHeaders, false},
	}

	router := mux.NewRouter()

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/")))
	uploadedFilesHandler := http.StripPrefix("/data/", decryptMiddleware(http.Dir(params.DataFolder+"/")))
	exitensionsHandler := http.StripPrefix("/ext/", s.extensionHandler(http.Dir("./web/static/ext/48px/")))

	// For static files
	router.PathPrefix("/static/").Handler(staticHandler)
	// For uploaded files
	router.PathPrefix("/data/").Handler(cacheMiddleware(uploadedFilesHandler, 365*24*60*60)) // cache for 365 days
	// For exitensions
	router.PathPrefix("/ext/").Handler(cacheMiddleware(exitensionsHandler, 7*24*60*60)) // cache for 7 days

	// Add usual routes
	for _, r := range routes {
		var handler http.Handler = r.handler
		if r.needAuth {
			handler = authMiddleware(r.handler)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}

	if params.Debug {
		// Add debug routes
		for _, r := range debugRoutes {
			var handler http.Handler = r.handler
			if r.needAuth {
				handler = authMiddleware(r.handler)
			}
			router.Path(r.path).Methods(r.methods).Handler(handler)
		}
	}

	var handler http.Handler = router
	if params.Debug {
		handler = debugMiddleware(router)
	}

	server := &http.Server{Addr: params.Port, Handler: handler}

	go func() {
		<-ctx.Done()

		s.logger.Infoln("Stopping web server")

		shutdown, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(shutdown); err != nil {
			s.logger.Errorf("can't gracefully shutdown server: %s\n", err)
		}
	}()

	s.logger.Infoln("Start web server")

	listenAndServe := server.ListenAndServe
	if params.IsTLS {
		listenAndServe = func() error {
			return server.ListenAndServeTLS("ssl/cert.cert", "ssl/key.key")
		}
	}

	// http.ErrServerClosed is a valid error
	if err := listenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// TODO

// processError is a wrapper over http.Error
func processError(w http.ResponseWriter, err string, code int) {
	if params.Debug {
		log.Errorf("Request error: %s (code: %d)\n", err, code)
	} else {
		// We should log server errors
		if 500 <= code && code < 600 {
			log.Errorf("Request error: %s (code: %d)\n", err, code)
		}
	}

	http.Error(w, err, code)
}
