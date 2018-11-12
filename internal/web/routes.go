package web

import "net/http"

type route struct {
	path     string
	methods  string
	handler  http.HandlerFunc
	needAuth bool
}

var routes = []route{
	{"/", "GET", index, true},
	// auth
	{"/login", "GET", login, false},
	{"/login", "POST", authentication, false},
	{"/logout", "POST", logout, true},
	// files
	{"/api/files", "GET", returnFiles, true},
	{"/api/files/new", "GET", returnFilesNew, true},
	{"/api/files/download", "GET", downloadFiles, true},
	{"/api/files", "POST", upload, true},
	{"/api/files/recover", "POST", recoverFile, true},
	{"/api/files/tags", "PUT", changeFileTags, true},
	{"/api/files/name", "PUT", changeFilename, true},
	{"/api/files/description", "PUT", changeFileDescription, true},
	{"/api/files", "DELETE", deleteFile, true},
	{"/api/files/recent", "GET", returnRecentFiles, true},
	// tags
	{"/api/tags", "GET", returnTags, true},
	{"/api/tags", "POST", addTag, true},
	{"/api/tags", "PUT", changeTag, true},
	{"/api/tags", "DELETE", deleteTag, true},
}

var debugRoutes = []route{
	{"/login", "OPTIONS", setDebugHeaders, false},
	{"/logout", "OPTIONS", setDebugHeaders, false},
	{"/api/files", "OPTIONS", setDebugHeaders, false},
	{"/api/files/new", "OPTIONS", setDebugHeaders, true}, // TODO: remove
	{"/api/files/recover", "OPTIONS", setDebugHeaders, false},
	{"/api/files/tags", "OPTIONS", setDebugHeaders, false},
	{"/api/files/name", "OPTIONS", setDebugHeaders, false},
	{"/api/files/description", "OPTIONS", setDebugHeaders, false},
	{"/api/tags", "OPTIONS", setDebugHeaders, false},
}
