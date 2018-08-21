package web

import (
	"io"
	"net/http"
	"os"
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

}

func auth(w http.ResponseWriter, r *http.Request) {

}

func checkAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
