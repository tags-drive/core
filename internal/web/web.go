// Package web is responsible for serving browser and API requests
package web

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

var ErrServerClosed = http.ErrServerClosed

// Start starts the server. It has to run in goroutine
//
// Functions stops when stopChan is closed. If there's any error, function will send it into errChan
// After stopping the server function sends ErrServerClosed into errChan
func Start(stopChan chan struct{}, errChan chan<- error) {
	var handler http.Handler

	server := &http.Server{Addr: params.Port, Handler: handler}

	go func() {
		if !params.IsTLS {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		} else {
			errChan <- errors.New("TLS isn't available")
		}
	}()

	<-stopChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		errChan <- err
	} else {
		errChan <- ErrServerClosed
	}
}
