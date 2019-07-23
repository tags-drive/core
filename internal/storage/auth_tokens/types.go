package auth

import "time"

type Config struct {
	Debug bool

	TokensJSONFile string
	Encrypt        bool
	PassPhrase     [32]byte

	MaxTokenLife time.Duration
}
