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
	{"/api/files/download", "GET", downloadFiles, true},
	{"/api/files", "POST", upload, true},
	{"/api/files", "PUT", changeFile, true},
	{"/api/files", "DELETE", deleteFile, true},
	{"/api/files/recent", "GET", returnRecentFiles, true},
	// tags
	{"/api/tags", "GET", returnTags, true},
	{"/api/tags", "POST", addTag, true},
	{"/api/tags", "PUT", changeTag, true},
	{"/api/tags", "DELETE", deleteTag, true},
}
