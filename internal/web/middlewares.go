package web

import (
	"net/http"

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
				http.Error(w, "need auth", http.StatusForbidden)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func decryptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
