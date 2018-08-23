// Package web is responsible for serving browser and API requests
package web

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/gorilla/mux"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

// Start starts the server. It has to run in goroutine
//
// Functions stops when stopChan is closed. If there's any error, function will send it into errChan
// After stopping the server function sends http.ErrServerClosed into errChan
func Start(stopChan chan struct{}, errChan chan<- error) {
	router := mux.NewRouter()
	// For static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	// For uploaded files
	router.PathPrefix("/data/").Handler(http.StripPrefix("/data/", decryptMiddleware(http.Dir(params.DataFolder+"/"))))
	for _, r := range routes {
		var handler http.Handler = r.handler
		if r.needAuth {
			handler = authMiddleware(r.handler)
		}
		router.Path(r.path).Methods(r.methods).Handler(handler)
	}

	server := &http.Server{Addr: params.Port, Handler: router}

	go func() {
		log.Infoln("Start web server")
		if !params.IsTLS {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		} else {
			errChan <- errors.New("TLS isn't available")
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
