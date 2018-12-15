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
// Functions stops when stopChan is closed. If there's any error, function will send it into errChan
func Start(stopChan chan struct{}, errChan chan<- error) {
	router := mux.NewRouter()

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/")))
	uploadedFilesHandler := http.StripPrefix("/data/", decryptMiddleware(http.Dir(params.DataFolder+"/")))
	exitensionsHandler := http.StripPrefix("/ext/", extensionHandler(http.Dir("./web/static/ext/48px/")))

	// For static files
	router.PathPrefix("/static/").Handler(staticHandler)
	// For uploaded files
	router.PathPrefix("/data/").Handler(cacheMiddleware(uploadedFilesHandler, 3*60)) // cache for 3 minutes
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
	localErrChan := make(chan error)
	go func() {
		log.Infoln("Start web server")

		listenAndServe := server.ListenAndServe
		if params.IsTLS {
			listenAndServe = func() error {
				return server.ListenAndServeTLS("ssl/cert.cert", "ssl/key.key")
			}
		}

		// http.ErrServerClosed is a valid error
		if err := listenAndServe(); err != nil && err != http.ErrServerClosed {
			localErrChan <- err
		}
	}()

	select {
	case err := <-localErrChan:
		close(localErrChan)
		errChan <- err
		// We don't have to shutdown server gracefully because it is down
		return
	case <-stopChan:
		// We have to shutdown server gracefully
	}

	log.Infoln("Stopping web server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	errChan <- server.Shutdown(ctx)
}

// Error is a wrapper over http.Error
func Error(w http.ResponseWriter, err string, code int) {
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
