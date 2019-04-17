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
	"github.com/tags-drive/core/internal/web/limiter"
)

const (
	authMaxRequests    = 1 // per second
	authLimiterTimeout = time.Second
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Server struct {
	fileStorage     cmd.FileStorageInterface
	tagStorage      cmd.TagStorageInterface
	authService     cmd.AuthServiceInterface
	authRateLimiter cmd.RateLimiterInterface

	httpServer *http.Server

	logger *clog.Logger
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

	s.authRateLimiter = limiter.NewRateLimiter(authMaxRequests, authLimiterTimeout)

	return s, nil
}

// Start starts the server. It has to be ran in goroutine
//
// Server stops when ctx.Done()
func (s *Server) Start() error {
	router := mux.NewRouter()

	// For static files
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/")))
	router.PathPrefix("/static/").Handler(staticHandler)

	// For uploaded files
	uploadedFilesHandler := http.StripPrefix("/data/", s.decryptMiddleware(http.Dir(params.DataFolder+"/")))
	router.PathPrefix("/data/").Handler(cacheMiddleware(uploadedFilesHandler, 60*60*24*14)) // cache for 14 days

	// For exitensions
	exitensionsHandler := http.StripPrefix("/ext/", s.extensionHandler(http.Dir("./web/static/ext/48px/")))
	router.PathPrefix("/ext/").Handler(cacheMiddleware(exitensionsHandler, 60*60*24*180)) // cache for 180 days

	// Add usual routes
	s.addDefaultRoutes(router)

	if params.Debug {
		s.addDebugRoutes(router)
		s.addPprofRoutes(router)
	}

	var handler http.Handler = router
	if params.Debug {
		handler = s.debugMiddleware(router)
	}

	s.httpServer = &http.Server{Addr: params.Port, Handler: handler}

	listenAndServe := s.httpServer.ListenAndServe
	if params.IsTLS {
		listenAndServe = func() error {
			return s.httpServer.ListenAndServeTLS("ssl/cert.cert", "ssl/key.key")
		}
	}

	s.logger.Infoln("start web server")

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
