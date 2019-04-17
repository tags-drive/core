package web

import (
	"net/http"
	"sync"
)

// processError is a wrapper over http.Error
func (s Server) processError(w http.ResponseWriter, err string, code int) {
	if s.config.Debug {
		s.logger.Errorf("request error: %s (code: %d)\n", err, code)
	} else {
		// We should log server errors
		if 500 <= code && code < 600 {
			s.logger.Errorf("request error: %s (code: %d)\n", err, code)
		}
	}

	http.Error(w, err, code)
}

func runPool(n int, data <-chan interface{}, worker func(data <-chan interface{})) {
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			worker(data)
			wg.Done()
		}()
	}

	wg.Wait()
}

// setDebugHeaders sets headers:
//   "Access-Control-Allow-Origin" - "*"
//   "Access-Control-Allow-Methods" - "POST, GET, OPTIONS, PUT, DELETE"
//   "Access-Control-Allow-Headers" - "Content-Type"
func setDebugHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
