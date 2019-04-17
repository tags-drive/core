package web

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/tags-drive/core/internal/params"
)

// POST /api/login
//
// Params:
//   - login: user's login
//   - password: password (sha256 checksum repeated 11 times)
//
// Response: cookie with auth token
//
func (s Server) authentication(w http.ResponseWriter, r *http.Request) {
	if !s.authRateLimiter.Take(r.RemoteAddr) {
		s.processError(w, "too many auth requests", http.StatusTooManyRequests)
		return
	}

	encrypt := func(s string) string {
		const repeats = 11

		hash := sha256.Sum256([]byte(s))
		for i := 0; i < repeats-1; i++ {
			hash = sha256.Sum256([]byte(hex.EncodeToString(hash[:])))
		}
		return hex.EncodeToString(hash[:])
	}

	var (
		login = r.FormValue("login")
		// password is already encrypted
		password = r.FormValue("password")
	)

	if password != encrypt(params.Password) || login != params.Login {
		if login != params.Login {
			s.processError(w, "invalid login", http.StatusBadRequest)
		} else {
			s.processError(w, "invalid password", http.StatusBadRequest)
		}

		s.logger.Warnf("%s tried to login with \"%s\" and \"%s\"\n", r.RemoteAddr, login, password)
		return
	}

	s.logger.Warnf("%s successfully logged in\n", r.RemoteAddr)

	token := s.authService.GenerateToken()
	s.authService.AddToken(token)
	http.SetCookie(w, &http.Cookie{Name: params.AuthCookieName, Value: token, HttpOnly: true, Expires: time.Now().Add(params.MaxTokenLife)})
}

// POST /api/logout – deletes auth cookie
//
// Params: -
//
// Response: -
//
func (s Server) logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(params.AuthCookieName)
	if err != nil {
		return
	}

	s.logger.Warnf("%s logged out\n", r.RemoteAddr)

	token := c.Value
	s.authService.DeleteToken(token)
	// Delete cookie
	http.SetCookie(w, &http.Cookie{Name: params.AuthCookieName, Expires: time.Unix(0, 0)})
}
