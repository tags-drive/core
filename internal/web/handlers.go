package web

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
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
