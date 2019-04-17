package auth

import "time"

type Config struct {
	Debug bool

	TokensJSONFile string
	Encrypt        bool
	PassPhrase     [32]byte

	MaxTokenLife time.Duration
}

// AuthServiceInterface provides methods for auth users
type AuthServiceInterface interface {
	// Start starts all background services
	StartBackgroundServices()

	// GenerateToken generates a new token. GenerateToken doesn't add new token, just return it!
	GenerateToken() string

	// AddToken adds passed token into storage
	AddToken(token string)

	// CheckToken returns true if token is in storage
	CheckToken(token string) bool

	// DeleteToken deletes token from a storage
	DeleteToken(token string)

	// Shutdown gracefully shutdown FileStorage
	Shutdown() error
}
