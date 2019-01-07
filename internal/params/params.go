// Package params provides global vars getted from environment
package params

import (
	"crypto/sha256"
	"errors"
	"os"
	"strings"
)

var (
	// Port of the server
	Port string
	// IsTLS defines should the program use https
	IsTLS bool
	// Login is a user login
	Login string
	// Password is a user password
	Password string
	// Debug defines is debug mode
	Debug bool
	// SkipLogin let use Tags Drive without loginning (for Debug only)
	SkipLogin bool
	// Encrypt defines, should the program encrypt files. False by default
	Encrypt bool
	// PassPhrase is used to encrypt files. Key is a sha256 sum of env "PASS_PHRASE"
	PassPhrase [32]byte
	// StorageType is a type of storage
	StorageType string
)

// Parse parses env vars
func Parse() error {
	// Default is ":80"
	Port = func() string {
		p := os.Getenv("PORT")
		if p == "" {
			return ":80"
		}

		if !strings.HasPrefix(p, ":") {
			p = ":" + p
		}
		return p
	}()

	// Default is "true"
	IsTLS = func() bool {
		value := os.Getenv("TLS")
		return !(strings.ToLower(value) == "false")
	}()

	// Default is "user"
	Login = func() (login string) {
		login = os.Getenv("LOGIN")
		if login == "" {
			return "user"
		}
		return
	}()

	// Default is "qwerty"
	Password = func() (pswrd string) {
		pswrd = os.Getenv("PSWRD")
		if pswrd == "" {
			return "qwerty"
		}
		return
	}()

	// Default is "false"
	Debug = func() bool {
		value := os.Getenv("DBG")
		return strings.ToLower(value) == "true"
	}()

	// Default is "false"
	SkipLogin = func() bool {
		value := os.Getenv("SKIP_LOGIN")
		if Debug && strings.ToLower(value) == "true" {
			return true
		}
		return false
	}()

	// Default is "false"
	Encrypt = func() bool {
		enc := os.Getenv("ENCRYPT")
		return enc == "true"
	}()

	if Encrypt {
		phrase := os.Getenv("PASS_PHRASE")
		if phrase == "" {
			return errors.New("wrong env config: PASS_PHRASE can't be empty with ENCRYPT=true")
		}
		PassPhrase = sha256.Sum256([]byte(phrase))
	}

	StorageType = func() string {
		return JSONStorage
	}()

	return nil
}
