package web

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/web/auth"
)

func mock(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Mock"))
}

func setDebugHeaders(w http.ResponseWriter, r *http.Request) {
	if params.Debug {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./web/index.html")
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, f)
	f.Close()
}

func login(w http.ResponseWriter, r *http.Request) {
	// Redirect to / if user is authorized
	c, err := r.Cookie(params.AuthCookieName)
	if err == nil && auth.CheckToken(c.Value) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	f, err := os.Open("./web/login.html")
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setDebugHeaders(w, r)
	io.Copy(w, f)
	f.Close()
}

func logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(params.AuthCookieName)
	if err != nil {
		return
	}

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
			Error(w, "invalid login", http.StatusBadRequest)
		} else {
			Error(w, "invalid password", http.StatusBadRequest)
		}
		return
	}

	token := auth.GenerateToken()
	auth.AddToken(token)
	http.SetCookie(w, &http.Cookie{Name: params.AuthCookieName, Value: token, HttpOnly: true, Expires: time.Now().Add(params.MaxTokenLife)})
}

func extensionHandler(dir http.Dir) http.Handler {
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
			io.Copy(w, f)
			f.Close()
			return
		}

		io.Copy(w, f)
		f.Close()
	})
}
