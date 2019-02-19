package web

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tags-drive/core/internal/params"
)

const (
	indexPath = "./web/index.html"
	loginPath = "./web/login.html"
)

// GET /
//
func (s Server) index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(indexPath)
	if err != nil {
		s.processError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		s.logger.Errorf("can't io.Copy() %s: %s\n", f.Name(), err)
	}
	f.Close()
}

// GET /login
//
func (s Server) login(w http.ResponseWriter, r *http.Request) {
	// Redirect to / if user is authorized
	c, err := r.Cookie(params.AuthCookieName)
	if err == nil && s.authService.CheckToken(c.Value) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	f, err := os.Open(loginPath)
	if err != nil {
		s.processError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		s.logger.Errorf("can't io.Copy() %s: %s\n", f.Name(), err)
	}
	f.Close()
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

// extensionHandler servers extensions
func (s Server) extensionHandler(dir http.Dir) http.Handler {
	const blankFilename = "_blank.png"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := r.URL.Path
		f, err := dir.Open(ext + ".png")
		if err != nil {
			// return blank icon
			f, err = dir.Open(blankFilename)
			if err != nil {
				return
			}
			_, err = io.Copy(w, f)
			if err != nil {
				s.logger.Errorf("can't io.Copy() %s.png: %s\n", ext, err)
			}
			f.Close()
			return
		}

		io.Copy(w, f)
		if err != nil {
			s.logger.Errorf("can't io.Copy() %s.png: %s\n", ext, err)
		}
		f.Close()
	})
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

func mock(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Mock"))
}
