package web

import (
	"io"
	"net/http"
	"os"
	"regexp"
)

const (
	indexPath  = "./web/index.html"
	loginPath  = "./web/login.html"
	mobilePage = "./web/mobile.html"
)

var isMobileDevice = regexp.MustCompile("Mobile|Android|iP(?:hone|od|ad)")

// GET /
//
func (s Server) index(w http.ResponseWriter, r *http.Request) {
	// Serve mobile devides (redirect to /mobile)
	userAgent := r.Header.Get("User-Agent")
	if isMobileDevice.MatchString(userAgent) {
		http.Redirect(w, r, "/mobile", http.StatusSeeOther)
		return
	}

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

// GET /mobile
//
func (s Server) mobile(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(mobilePage)
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

// GET /share
//
// Params:
//   - shareToken: share token
//
// Response: index or mobile page
//
func (s Server) share(w http.ResponseWriter, r *http.Request) {
	shareToken := r.FormValue("shareToken")
	if !s.shareStorage.CheckToken(shareToken) {
		s.processError(w, "invalid share token", http.StatusBadRequest)
		return
	}

	// Serve mobile devices (redirect to /mobile)
	userAgent := r.Header.Get("User-Agent")
	if isMobileDevice.MatchString(userAgent) {
		s.mobile(w, r)
		return
	}

	s.index(w, r)
}

// GET /login
//
func (s Server) login(w http.ResponseWriter, r *http.Request) {
	// Redirect to / if user is authorized
	c, err := r.Cookie(s.config.AuthCookieName)
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

// GET /version
//
// Response: backend version
//
func (s Server) backendVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.config.Version))
}

// extensionHandler servers extensions
func (s Server) extensionHandler(dir http.Dir) http.Handler {
	const blankFilename = "_blank.png"
	const iconExt = ".png"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := r.URL.Path

		if f, err := dir.Open(ext + iconExt); err == nil {
			// Return existing icon
			io.Copy(w, f)
			if err != nil {
				s.logger.Errorf("can't io.Copy() %s.png: %s\n", ext, err)
			}
			f.Close()
			return
		}

		// return blank icon
		f, err := dir.Open(blankFilename)
		if err != nil {
			return
		}
		_, err = io.Copy(w, f)
		if err != nil {
			s.logger.Errorf("can't io.Copy() %s.png: %s\n", ext, err)
		}
		f.Close()
		return
	})
}

func mock(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Mock"))
}
