// Package web is responsible for serving browser and API requests
package web

import (
	"context"
	"net/http"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/share"
	"github.com/tags-drive/core/internal/storage/tags"
	"github.com/tags-drive/core/internal/web/auth"
	"github.com/tags-drive/core/internal/web/limiter"
)

const (
	authMaxRequests    = 1 // per second
	authLimiterTimeout = time.Second
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Server struct {
	config Config

	fileStorage  *files.FileStorage
	tagStorage   *tags.TagStorage
	shareStorage *share.ShareStorage

	authService     *auth.AuthService
	authRateLimiter limiter.RateLimiterInterface

	httpServer *http.Server

	logger *clog.Logger
}

// NewWebServer just creates new Web struct
func NewWebServer(cnf Config, fs *files.FileStorage, ts *tags.TagStorage, lg *clog.Logger) (*Server, error) {
	s := &Server{
		config:      cnf,
		fileStorage: fs,
		tagStorage:  ts,
		logger:      lg,
	}

	var err error

	// Share service
	shareConfig := share.Config{
		ShareTokenJSONFile: cnf.ShareTokensJSONFile,
		Encrypt:            cnf.Encrypt,
		PassPhrase:         cnf.PassPhrase,
	}
	s.shareStorage, err = share.NewShareStorage(shareConfig, fs, lg)
	if err != nil {
		return nil, errors.Wrap(err, "can't init new Share Storage")
	}

	// Auth service
	authConfig := auth.Config{
		Debug:          cnf.Debug,
		TokensJSONFile: cnf.AuthTokensJSONFile,
		Encrypt:        cnf.Encrypt,
		PassPhrase:     cnf.PassPhrase,
		MaxTokenLife:   cnf.MaxTokenLife,
	}
	s.authService, err = auth.NewAuthService(authConfig, lg)
	if err != nil {
		return nil, errors.Wrap(err, "can't init new Auth Service")
	}

	// Rate limiter
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
	router.PathPrefix("/data/").Handler(s.authMiddleware(s.serveData(), true))

	// For exitensions
	exitensionsHandler := http.StripPrefix("/ext/", s.extensionHandler(http.Dir("./web/static/ext/48px/")))
	router.PathPrefix("/ext/").Handler(cacheMiddleware(exitensionsHandler, 60*60*24*180)) // cache for 180 days

	// Add usual routes
	s.addDefaultRoutes(router)

	if s.config.Debug {
		s.addDebugRoutes(router)
		s.addPprofRoutes(router)
	}

	var handler http.Handler = router
	if s.config.Debug {
		handler = s.debugMiddleware(router)
	}

	s.httpServer = &http.Server{Addr: s.config.Port, Handler: handler}

	listenAndServe := s.httpServer.ListenAndServe
	if s.config.IsTLS {
		listenAndServe = func() error {
			return s.httpServer.ListenAndServeTLS("ssl/cert.cert", "ssl/key.key")
		}
	}

	// Start background services
	s.authService.StartBackgroundServices()

	s.logger.Debugln("start web server")

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

	// Shutdown share storage
	if err := s.shareStorage.Shutdown(); err != nil {
		s.logger.Warnf("can't shutdown shareStorage gracefully: %s\n", err)
	}

	return serverErr
}
