package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/minio/sio"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/web/auth"
)

func authMiddleware(h http.Handler) http.Handler {
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
			return auth.CheckToken(token)
		}()

		if !validToken {
			// Redirect won't help
			if r.Method != "GET" {
				Error(w, "need auth", http.StatusForbidden)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func decryptMiddleware(dir http.Dir) http.Handler {
	if !params.Encrypt {
		return http.FileServer(dir)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Path
		f, err := dir.Open(fileName)
		if err != nil {
			Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		_, err = sio.Decrypt(w, f, sio.Config{Key: params.Key[:]})
		if err != nil {
			Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func debugMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't log favicon.ico
		if !strings.HasSuffix(r.URL.Path, "favicon.ico") {
			fmt.Printf("%s %s\n", r.Method, r.URL.Path)

			r.ParseForm()
			if len(r.Form) > 0 {
				prefix := strings.Repeat(" ", len(r.Method))

				space := 0
				for k := range r.Form {
					if space < len(k) {
						space = len(k)
					}
				}

				for k, v := range r.Form {
					fmt.Printf("%s %v: ", prefix, k)
					for i := 0; i < space-len(k); i++ {
						fmt.Print(" ")
					}
					fmt.Println(v)
				}
			}
		}

		setDebugHeaders(w, r)

		h.ServeHTTP(w, r)
	})
}

func cacheMiddleware(h http.Handler, maxAge int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
		h.ServeHTTP(w, r)
	})
}
