package limiter

type RateLimiterInterface interface {
	Take(remoteAddr string) bool
}
