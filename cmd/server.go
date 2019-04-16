package cmd

// ServerInterface provides methods for interactions web server
type ServerInterface interface {
	Start() error

	// Shutdown gracefully shutdowns server
	Shutdown() error
}

// AuthService provides methods for auth users
type AuthService interface {
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
