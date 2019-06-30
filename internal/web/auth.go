package web

import (
	"net/http"
	"time"
)

// GET /api/user - lets check is user authorized
//
//
func (s Server) checkUser(w http.ResponseWriter, r *http.Request) {
	const response = `{"authorized":true}`

	if s.config.SkipLogin {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
		return
	}

	cookie, err := r.Cookie(s.config.AuthCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := cookie.Value
	if !s.authService.CheckToken(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

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

	var (
		login = r.FormValue("login")
		// password is already encrypted
		password = r.FormValue("password")
	)

	if password != s.config.Password || login != s.config.Login {
		if login != s.config.Login {
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

	http.SetCookie(w, &http.Cookie{
		Name:     s.config.AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(s.config.MaxTokenLife),
	})
}

// POST /api/logout – deletes auth cookie
//
// Params: -
//
// Response: -
//
func (s Server) logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(s.config.AuthCookieName)
	if err != nil {
		return
	}

	s.logger.Warnf("%s logged out\n", r.RemoteAddr)

	token := c.Value
	s.authService.DeleteToken(token)

	// Delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:    s.config.AuthCookieName,
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
}
