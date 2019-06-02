package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/minio/sio"
)

// authMiddleware checks if a user is authorized. If the user isn't and resource is shareable,
// it checks if "shareToken" passed and a token is valid.
func (s Server) authMiddleware(h http.Handler, shareable bool) http.Handler {
	checkAuth := func(r *http.Request) bool {
		c, err := r.Cookie(s.config.AuthCookieName)
		if err != nil {
			return false
		}

		token := c.Value
		return s.authService.CheckToken(token)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.SkipLogin {
			h.ServeHTTP(w, r)
			return
		}

		state := &requestState{}

		if checkAuth(r) {
			state.authorized = true
		}

		if shareable {
			shareToken := r.FormValue("shareToken")
			if shareToken != "" {
				state.shareToken = shareToken

				if !s.shareStorage.CheckToken(shareToken) {
					s.processError(w, "invalid share token", http.StatusBadRequest)
					return
				}

				// Limit access even when user is authorized
				state.shareAccess = true
			}
		}

		if !state.authorized && !state.shareAccess {
			if strings.HasPrefix(r.URL.String(), "/api/") || strings.HasPrefix(r.URL.String(), "/data/") {
				// Redirect won't help
				s.processError(w, "need auth", http.StatusUnauthorized)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := storeRequestState(context.Background(), state)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

func (s Server) decryptMiddleware(dir http.Dir) http.Handler {
	if !s.config.Encrypt {
		return http.FileServer(dir)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Path
		f, err := dir.Open(fileName)
		if err != nil {
			s.processError(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		_, err = sio.Decrypt(w, f, sio.Config{Key: s.config.PassPhrase[:]})
		if err != nil {
			s.processError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// cacheMiddleware sets "Cache-Control" header
func cacheMiddleware(h http.Handler, maxAge int64) http.Handler {
	maxAgeString := fmt.Sprintf("max-age=%d", maxAge)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", maxAgeString)
		w.Header().Add("Cache-Control", "private")
		h.ServeHTTP(w, r)
	})
}

// debugMiddleware logs requests and sets debug headers
func (s Server) debugMiddleware(h http.Handler) http.Handler {
	// Can be changed to debug
	const (
		// Log settings

		logDataRequests        = false
		logExtRequests         = false
		logStaticFilesRequests = false
		logFaviconRequests     = false

		// Builder of log records settings

		printForm = true
		// Note: request will be logged after ServeHTTP finish (it can take some time)
		printServingDuration = true
	)

	const (
		// time len + space (1) + [DBG] (5) + space (1) + method len (?) + space (1)
		indentionOffset = len(clog.DefaultTimeLayout) + 1 + 5 + 1 + 1

		// indentionOriginString helps not to allocate new memory with strings.Repeat()
		// It contains 50 spaces (must be enough forever)
		indentionOriginalString = "                                                  "
	)

	shouldLog := func(urlPath string) bool {
		// Check sorted by requests popularity

		// Don't log data requests
		if !logDataRequests && strings.HasPrefix(urlPath, "/data/") {
			return false
		}

		// Don't log extensions requests
		if !logExtRequests && strings.HasPrefix(urlPath, "/ext/") {
			return false
		}

		// Don't log static files
		if !logStaticFilesRequests && strings.HasPrefix(urlPath, "/static/") {
			return false
		}

		// Don't log favicon.ico
		if !logFaviconRequests && strings.HasSuffix(urlPath, "favicon.ico") {
			return false
		}

		return true
	}

	// buildLogRecord builds a string to log. The string contains information about a request.
	buildLogRecord := func(r *http.Request, duration time.Duration) string {
		builder := new(strings.Builder)

		// Add main info
		builder.WriteString(r.Method)
		builder.WriteByte(' ')
		builder.WriteString(r.URL.Path)
		builder.WriteByte('\n')

		indention := indentionOriginalString[:indentionOffset+len(r.Method)]

		// Add form
		r.ParseForm()
		if printForm && len(r.Form) > 0 {
			space := 0
			for k := range r.Form {
				if space < len(k) {
					space = len(k)
				}
			}
			space++

			builder.WriteString(indention)
			builder.WriteString("Form: \n")

			for k, v := range r.Form {
				p := strings.Repeat(" ", space-len(k))

				builder.WriteString(indention)
				builder.WriteByte(' ')
				builder.WriteByte('-')
				builder.WriteByte(' ')
				builder.WriteString(k)
				builder.WriteString(p)
				fmt.Fprint(builder, v)
				builder.WriteByte('\n')
			}
		}

		// Add duration
		if printServingDuration {
			builder.WriteString(indention)
			builder.WriteString("Duration: ")
			builder.WriteString(duration.String())
			builder.WriteByte('\n')
		}

		return builder.String()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldLog := shouldLog(r.URL.Path)

		if shouldLog && !printServingDuration {
			// Print immediately without duration
			s.logger.Debug(buildLogRecord(r, 0))
		}

		setDebugHeaders(w, r)

		now := time.Now()

		h.ServeHTTP(w, r)

		if shouldLog && printServingDuration {
			// Print with duration
			s.logger.Debug(buildLogRecord(r, time.Since(now)))
		}
	})
}
