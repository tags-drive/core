package web

import (
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
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
