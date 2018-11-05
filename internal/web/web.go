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

// Start starts the server. It has to run in goroutine
//
// Functions stops when stopChan is closed. If there's any error, function will send it into errChan
// After stopping the server function sends http.ErrServerClosed into errChan
func Start(stopChan chan struct{}, errChan chan<- error) {
	router := mux.NewRouter()

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/")))
	uploadedFilesHandler := http.StripPrefix("/data/", decryptMiddleware(http.Dir(params.DataFolder+"/")))
	exitensionsHandler := http.StripPrefix("/ext/", extensionHandler(http.Dir("./web/static/ext/48px/")))

	if params.Debug {
		staticHandler = debugFileServeMiddleware(staticHandler)
		uploadedFilesHandler = debugFileServeMiddleware(uploadedFilesHandler)
		exitensionsHandler = debugFileServeMiddleware(exitensionsHandler)
	}

	// For static files
	router.PathPrefix("/static/").Handler(staticHandler)
	// For uploaded files
	router.PathPrefix("/data/").Handler(uploadedFilesHandler)
	// For exitensions
	router.PathPrefix("/ext/").Handler(exitensionsHandler)

	for _, r := range routes {
		var handler http.Handler = r.handler
		if r.needAuth {
			handler = authMiddleware(r.handler)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}

	var handler http.Handler = router
	if params.Debug {
		handler = debugMiddleware(router)
	}

	server := &http.Server{Addr: params.Port, Handler: handler}

	go func() {
		log.Infoln("Start web server")
		if !params.IsTLS {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		} else {
			if err := server.ListenAndServeTLS("ssl/cert.cert", "ssl/key.key"); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}
	}()

	<-stopChan
	log.Infoln("Stopping web server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		errChan <- err
	} else {
		errChan <- http.ErrServerClosed
	}
}

func Error(w http.ResponseWriter, err string, code int) {
	if params.Debug {
		log.Errorf("Request error: %s (code: %d)\n", err, code)
	}

	http.Error(w, err, code)
}
