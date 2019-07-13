package web

import (
	"fmt"
	"html/template"
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
		state := &requestState{}

		if checkAuth(r) || s.config.SkipLogin {
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

		ctx := storeRequestState(r.Context(), state)
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

func (s Server) openGraphMiddleware(h http.Handler) http.Handler {
	type ogData struct {
		SiteURL  string
		ImageURL string
	}

	const (
		imagePath = "/static/icons/tag-1024x1024.png"
		ogPage    = `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<title>Tags Drive</title>

				<meta name="description" content="Open source self-hosted cloud drive with tags" />
				<meta name="robots" content="noindex, nofollow" />

				<meta property="og:title" content="Tags Drive" />
				<meta property="og:type" content="website" />
				<meta property="og:description" content="Open source self-hosted cloud drive with tags" />
				<meta property="og:url" content="{{.SiteURL}}" />
				<meta property="og:image" content="{{.ImageURL}}" />
				<meta property="og:image:width" content="1024" />
				<meta property="og:image:height" content="1024" />
				<meta property="og:image:type" content="image/png" />
				<meta property="og:image:alt" content="tag" />
			</head>
			</html>`
	)

	ogPageTemplate := template.Must(template.New("openGraphTemplate").Parse(ogPage))

	isCrawler := func(userAgent string) bool {
		// List of popular crawlers
		crawlersUserAgents := [...]string{
			"TelegramBot",         // Telegram
			"Twitterbot",          // Twitter
			"facebookexternalhit", // Facebook
			"WhatsApp",            // WhatsApp
			"vkShare",             // VK
		}

		for i := range crawlersUserAgents {
			if strings.Contains(userAgent, crawlersUserAgents[i]) {
				return true
			}
		}

		return false
	}

	scheme := "http://"
	if s.config.IsTLS {
		scheme = "https://"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isCrawler(r.UserAgent()) {
			// Serve the request as usual
			h.ServeHTTP(w, r)
			return
		}

		data := ogData{
			SiteURL:  scheme + r.Host,
			ImageURL: scheme + r.Host + imagePath,
		}

		ogPageTemplate.Execute(w, data)
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
