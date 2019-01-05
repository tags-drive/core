// Package web is responsible for serving browser and API requests
package web

import (
	"context"
	"net/http"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/gorilla/mux"

	"github.com/tags-drive/core/internal/params"
)

// Start starts the server. It has to be ran in goroutine
//
// Server stops when ctx.Done()
func Start(ctx context.Context) error {
	router := mux.NewRouter()

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/")))
	uploadedFilesHandler := http.StripPrefix("/data/", decryptMiddleware(http.Dir(params.DataFolder+"/")))
	exitensionsHandler := http.StripPrefix("/ext/", extensionHandler(http.Dir("./web/static/ext/48px/")))

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

		log.Infoln("Stopping web server")

		shutdown, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(shutdown); err != nil {
			log.Errorf("can't gracefully shutdown server: %s\n", err)
		}
	}()

	log.Infoln("Start web server")

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
