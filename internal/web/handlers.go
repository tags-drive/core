package web

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ShoshinNikita/tags-drive/internal/params"
	"github.com/ShoshinNikita/tags-drive/internal/web/auth"
)

func index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, f)
}

func login(w http.ResponseWriter, r *http.Request) {
	// Redirect to / if user is authorized
	c, err := r.Cookie(params.AuthCookieName)
	if err == nil && auth.CheckToken(c.Value) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	f, err := os.Open("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, f)
}

func logout(w http.ResponseWriter, r *http.Request) {
	// We can skip err, because we already checked is there a cookie in authMiddleware()
	c, _ := r.Cookie(params.AuthCookieName)

	token := c.Value
	auth.DeleteToken(token)
	// Delete cookie
	http.SetCookie(w, &http.Cookie{Name: params.AuthCookieName, Expires: time.Unix(0, 0)})
}

func authentication(w http.ResponseWriter, r *http.Request) {
	encrypt := func(s string) string {
		hash := sha256.Sum256([]byte(s))
		for i := 0; i < 10; i++ {
			hash = sha256.Sum256([]byte(hex.EncodeToString(hash[:])))
		}
		return hex.EncodeToString(hash[:])
	}

	var (
		login    = r.FormValue("login")
		password = r.FormValue("password")
	)

	if password != encrypt(params.Password) || login != params.Login {
		if login != params.Login {
			http.Error(w, "invalid login", http.StatusBadRequest)
		} else {
			http.Error(w, "invalid password", http.StatusBadRequest)
		}
		return
	}

	token := auth.GenerateToken()
	auth.AddToken(token)
	http.SetCookie(w, &http.Cookie{Name: params.AuthCookieName, Value: token, HttpOnly: true, Expires: time.Now().Add(params.MaxTokenLife)})
}

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
