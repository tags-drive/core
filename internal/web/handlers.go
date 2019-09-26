package web

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
		s.processError(w, "can't load index page", http.StatusInternalServerError, err)
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
		s.processError(w, "can't load mobile page", http.StatusInternalServerError, err)
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
	if !s.shareService.CheckToken(shareToken) {
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
		s.processError(w, "can't load login page", http.StatusInternalServerError, err)
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		s.logger.Errorf("can't io.Copy() %s: %s\n", f.Name(), err)
	}
	f.Close()
}

// GET /api/version
//
// Response: backend version
//
func (s Server) backendVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.config.Version))
}

// GET /api/ping
//
// Response: http.StatusOK (200)
//
func (s Server) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// GET /data/.../{id}
//
// Params:
//   - shareToken (optional): share token
//
func (s Server) serveData() (handler http.Handler) {
	getFileID := func(url string) (id int, ok bool) {
		var strID string
		for i := len(url) - 1; i >= 0; i-- {
			// Use all chars after last /
			if url[i] == '/' {
				strID = url[i+1:]
				break
			}
		}

		if strID == "" {
			return 0, false
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			return 0, false
		}

		return id, true
	}

	shouldUseCache := func(modTime time.Time, r *http.Request) bool {
		// Next code are taken from "net/http/fs.go", "checkIfModifiedSince" function

		ims := r.Header.Get("If-Modified-Since")
		if ims == "" {
			return false
		}

		t, err := http.ParseTime(ims)
		if err != nil {
			return false
		}

		// The Date-Modified header truncates sub-second precision, so
		// use mtime < t+1s instead of mtime <= t to check for unmodified.
		return modTime.Before(t.Add(1 * time.Second))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.EscapedPath()

		id, ok := getFileID(url)
		if !ok {
			s.processError(w, "invalid file id", http.StatusBadRequest)
			return
		}

		state, ok := getRequestState(r.Context())
		if !ok {
			s.processError(w, "share token doesn't grant access to this file", http.StatusForbidden)
			return
		}

		if state.shareAccess {
			// Have to check if a token grants access to the file
			if !s.shareService.CheckFile(state.shareToken, id) {
				s.processError(w, "share token doesn't grant access to this file", http.StatusForbidden)
				return
			}
		}

		if !s.fileStorage.CheckFile(id) {
			s.processError(w, "file doesn't exist", http.StatusNotFound)
			return
		}

		// Get a file to learn the modification time
		file, err := s.fileStorage.GetFile(id)
		if err != nil {
			s.processError(w, "can't load file", http.StatusInternalServerError, "can't get file from FileStorage:", err)
			return
		}

		// Always add "Last-Modified" header
		w.Header().Set("Last-Modified", file.AddTime.Format(http.TimeFormat))

		if shouldUseCache(file.AddTime, r) {
			// Response with http.StatusNotModified (304)

			// From "net/http/fs.go", "writeNotModified" function
			h := w.Header()
			delete(h, "Content-Type")
			delete(h, "Content-Length")
			w.WriteHeader(http.StatusNotModified)
			return
		}

		// Write a file
		resized := strings.Contains(url, "resized")
		err = s.fileStorage.CopyFile(w, id, resized)
		if err != nil {
			s.processError(w, "can't load file", http.StatusInternalServerError, err)
		}
	})
}

// extensionHandler servers extensions
func (s Server) extensionHandler(dir http.Dir) http.Handler {
	const blankFilename = "_blank.png"
	const iconExt = ".png"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := r.URL.Path

		if f, err := dir.Open(ext + iconExt); err == nil {
			// Return existing icon
			if _, err = io.Copy(w, f); err != nil {
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
