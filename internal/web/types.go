package web

type Config struct {
	Debug bool

	DataFolder string

	Port  string
	IsTLS bool

	Login          string
	Password       string
	SkipLogin      string
	AuthCookieName string
	MaxTokenLife   string
	TokensJSONFile string

	Encrypt    bool
	PassPhrase [32]byte

	Version string
}

// ServerInterface provides methods for interactions web server
type ServerInterface interface {
	Start() error

	// Shutdown gracefully shutdowns server
	Shutdown() error
}
