package web

import (
	"fmt"
	"net/http"
	"strings"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/minio/sio"

	"github.com/tags-drive/core/internal/params"
)

func (s Server) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if params.SkipLogin {
			h.ServeHTTP(w, r)
			return
		}

		validToken := func() bool {
			c, err := r.Cookie(params.AuthCookieName)
			if err != nil {
				return false
			}

			token := c.Value
			return s.authService.CheckToken(token)
		}()

		if !validToken {
			// Redirect won't help
			if r.Method != "GET" {
				s.processError(w, "need auth", http.StatusForbidden)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (s Server) decryptMiddleware(dir http.Dir) http.Handler {
	if !params.Encrypt {
		return http.FileServer(dir)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Path
		f, err := dir.Open(fileName)
		if err != nil {
			s.processError(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		_, err = sio.Decrypt(w, f, sio.Config{Key: params.PassPhrase[:]})
		if err != nil {
			s.processError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// debugMiddleware logs requests and sets debug headers
func (s Server) debugMiddleware(h http.Handler) http.Handler {
	const (
		// Can be changed to debug
		logDataRequests = false
		logExtRequests  = false
		logStaticFiles  = false
		logFavicon      = false
	)

	// time len + space (1) + [DBG] (5) + space (1) + method len (?) + space (1)
	const prefixOffset = len(clog.DefaultTimeLayout) + 1 + 5 + 1 + 1

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			setDebugHeaders(w, r)

			h.ServeHTTP(w, r)
		}()

		// Don't log favicon.ico
		if !logFavicon && strings.HasSuffix(r.URL.Path, "favicon.ico") {
			return
		}

		// Don't log data requests
		if !logDataRequests && strings.HasPrefix(r.URL.Path, "/data/") {
			return
		}

		// Don't log extensions requests
		if !logExtRequests && strings.HasPrefix(r.URL.Path, "/ext/") {
			return
		}

		// Don't log static files
		if !logStaticFiles && strings.HasPrefix(r.URL.Path, "/static/") {
			return
		}

		builder := new(strings.Builder)
		builder.WriteString(r.Method)
		builder.WriteByte(' ')
		builder.WriteString(r.URL.Path)
		builder.WriteByte('\n')

		r.ParseForm()
		if len(r.Form) > 0 {

			prefix := strings.Repeat(" ", prefixOffset+len(r.Method))

			space := 0
			for k := range r.Form {
				if space < len(k) {
					space = len(k)
				}
			}
			space++

			for k, v := range r.Form {
				p := strings.Repeat(" ", space-len(k))
				builder.WriteString(prefix)
				builder.WriteString(k)
				builder.WriteString(p)
				fmt.Fprint(builder, v)
				builder.WriteByte('\n')
			}
		}

		s.logger.Debug(builder.String())
	})
}

// cacheMiddleware sets "Cache-Control" header
func cacheMiddleware(h http.Handler, maxAge int64) http.Handler {
	maxAgeString := fmt.Sprintf("max-age=%d", maxAge)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", maxAgeString)
		w.Header().Add("Cache-Control", "private")
		h.ServeHTTP(w, r)
	})
}
