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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		builder := new(strings.Builder)

		// Don't log favicon.ico
		if !strings.HasSuffix(r.URL.Path, "favicon.ico") {
			builder.WriteString(r.Method)
			builder.WriteByte(' ')
			builder.WriteString(r.URL.Path)
			builder.WriteByte('\n')

			r.ParseForm()
			if len(r.Form) > 0 {

				prefix := strings.Repeat(" ", len(clog.DefaultTimeLayout)+len(r.Method)+2)

				space := 0
				for k := range r.Form {
					if space < len(k) {
						space = len(k)
					}
				}

				for k, v := range r.Form {
					p := strings.Repeat(" ", space-len(k))
					builder.WriteString(prefix)
					builder.WriteString(k)
					builder.WriteString(p)
					fmt.Fprint(builder, v)
					builder.WriteByte('\n')
				}
			}

			s.logger.Print(builder.String())
		}

		setDebugHeaders(w, r)

		h.ServeHTTP(w, r)
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
