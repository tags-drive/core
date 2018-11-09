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
	{"/api/files/recover", "POST", recoverFile, true},
	{"/api/files/recover", "OPTIONS", setDebugHeaders, true},
	{"/api/files/tags", "PUT", changeFileTags, true},
	{"/api/files/tags", "OPTIONS", setDebugHeaders, true},
	{"/api/files/name", "PUT", changeFilename, true},
	{"/api/files/name", "OPTIONS", setDebugHeaders, true},
	{"/api/files/description", "PUT", changeFileDescription, true},
	{"/api/files/description", "OPTIONS", setDebugHeaders, true},
	{"/api/files", "DELETE", deleteFile, true},
	{"/api/files", "OPTIONS", setDebugHeaders, true},
	{"/api/files/recent", "GET", returnRecentFiles, true},
	// tags
	{"/api/tags", "GET", returnTags, true},
	{"/api/tags", "POST", addTag, true},
	{"/api/tags", "PUT", changeTag, true},
	{"/api/tags", "OPTIONS", setDebugHeaders, true},
	{"/api/tags", "DELETE", deleteTag, true},
}
