package web

import "net/http"

var routes = []struct {
	path     string
	methods  string
	handler  http.HandlerFunc
	needAuth bool
}{
	{"/", "GET", index, false}, // index should check is userdata correct itself
	{"/login", "GET", login, false},
	{"/login", "POST", auth, false},
	// files
	{"/api/files", "GET", returnFiles, true},
	{"/api/files", "POST", upload, true},
	{"/api/files", "PUT", renameFile, true},
	{"/api/files", "DELETE", deleteFile, true},
	{"/api/files/recent", "GET", returnRecentFiles, true},
	// tags
	{"/api/tags", "PUT", changeTags, true},
}
