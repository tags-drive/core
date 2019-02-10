// Package web is responsible for serving browser and API requests
package web

import (
	"context"
	"net/http"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/web/auth"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Server struct {
	fileStorage cmd.FileStorageInterface
	tagStorage  cmd.TagStorageInterface
	authService cmd.AuthService

	logger *clog.Logger

	httpServer *http.Server
}

type route struct {
	path     string
	methods  string
	handler  http.HandlerFunc
	needAuth bool
}

// NewWebServer just creates new Web struct. It doesn't call any Init functions
func NewWebServer(fs cmd.FileStorageInterface, ts cmd.TagStorageInterface, lg *clog.Logger) (*Server, error) {
	s := &Server{
		fileStorage: fs,
		tagStorage:  ts,
		logger:      lg,
	}

	var err error

	s.authService, err = auth.NewAuthService(lg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Start starts the server. It has to be ran in goroutine
//
// Server stops when ctx.Done()
func (s *Server) Start() error {
	// Init routes
	var routes = []route{
		{"/", "GET", s.index, true},

		// Auth
		{"/login", "GET", s.login, false},
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

	var debugRoutes = []route{
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

	router := mux.NewRouter()

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/")))
	uploadedFilesHandler := http.StripPrefix("/data/", s.decryptMiddleware(http.Dir(params.DataFolder+"/")))
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
			handler = s.authMiddleware(r.handler)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}

	if params.Debug {
		// Add debug routes
		for _, r := range debugRoutes {
			var handler http.Handler = r.handler
			if r.needAuth {
				handler = s.authMiddleware(r.handler)
			}
			router.Path(r.path).Methods(r.methods).Handler(handler)
		}
	}

	var handler http.Handler = router
	if params.Debug {
		handler = s.debugMiddleware(router)
	}

	s.httpServer = &http.Server{Addr: params.Port, Handler: handler}

	s.logger.Infoln("start web server")

	listenAndServe := s.httpServer.ListenAndServe
	if params.IsTLS {
		listenAndServe = func() error {
			return s.httpServer.ListenAndServeTLS("ssl/cert.cert", "ssl/key.key")
		}
	}

	// http.ErrServerClosed is a valid error
	if err := listenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s Server) Shutdown() error {
	shutdown, release := context.WithTimeout(context.Background(), time.Second*10)
	defer release()

	s.httpServer.SetKeepAlivesEnabled(false)

	serverErr := s.httpServer.Shutdown(shutdown)

	// Shutdown auth service
	if err := s.authService.Shutdown(); err != nil {
		s.logger.Warnf("can't shutdown authService gracefully: %s\n", err)
	}

	return serverErr
}

// processError is a wrapper over http.Error
func (s Server) processError(w http.ResponseWriter, err string, code int) {
	if params.Debug {
		s.logger.Errorf("request error: %s (code: %d)\n", err, code)
	} else {
		// We should log server errors
		if 500 <= code && code < 600 {
			s.logger.Errorf("request error: %s (code: %d)\n", err, code)
		}
	}

	http.Error(w, err, code)
}
