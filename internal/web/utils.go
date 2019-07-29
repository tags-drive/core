package web

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// processError is a wrapper over http.Error
func (s Server) processError(w http.ResponseWriter, err string, code int, reasons ...interface{}) {
	// Send error at first
	http.Error(w, err, code)

	msg := fmt.Sprintf("request error: %s (code: %d)", err, code)

	if len(reasons) == 0 {
		// Log an errors if it's a server error or Debug Mode is enabled.
		if 500 <= code && code < 600 || s.config.Debug {
			s.logger.Errorln(msg)
		}
		return
	}

	// We have to build an error message
	var b strings.Builder

	b.WriteString(msg)
	b.WriteString(". Reason: ")
	for i, r := range reasons {
		b.WriteString(fmt.Sprintf("%+v", r))

		if i < len(reasons)-1 {
			b.WriteByte(' ')
		}
	}

	s.logger.Errorln(b.String())
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
