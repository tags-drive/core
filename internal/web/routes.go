package web

import "net/http"

var routes = []struct {
	path     string
	methods  string
	handler  http.HandlerFunc
	needAuth bool
}{
	{"/", "GET", index, true},
	// auth
	{"/login", "GET", login, false},
	{"/login", "POST", authentication, false},
	{"/logout", "POST", logout, true},
	// files
	{"/api/files", "GET", returnFiles, true},
	{"/api/files", "POST", upload, true},
	{"/api/files", "PUT", renameFile, true},
	{"/api/files", "DELETE", deleteFile, true},
	{"/api/files/recent", "GET", returnRecentFiles, true},
	// tags
	{"/api/tags", "PUT", changeTags, true},
}
