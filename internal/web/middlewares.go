package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/minio/sio"

	"github.com/ShoshinNikita/tags-drive/internal/params"
	"github.com/ShoshinNikita/tags-drive/internal/web/auth"
)

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validToken := func() bool {
			c, err := r.Cookie(params.AuthCookieName)
			if err != nil {
				return false
			}

			token := c.Value
			if !auth.CheckToken(token) {
				return false
			}

			return true
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
			fmt.Printf("%s ", r.Method)
			fmt.Print(r.URL.Path)
			if r.URL.RawQuery != "" {
				fmt.Printf("?%s", r.URL.RawQuery)
			}
			fmt.Println()
		}

		h.ServeHTTP(w, r)
	})
}
