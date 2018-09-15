// Package params provides global vars getted from environment
package params

import (
	"crypto/sha256"
	"os"
	"strings"
	"time"
)

const (
	// DataFolder is a folder, in which all files are kept
	DataFolder = "data"
	// ResizedImagesFolder is a folder, in which all resized images are kept
	ResizedImagesFolder = "data/resized"
	// Files is a json file with info about the files
	Files = "configs/files.json"
	// TokensFile is a json file with list of tokens
	TokensFile = "configs/tokens.json"
	// TagsFile is a json file with list of tags (with name and color)
	TagsFile = "configs/tags.json"
	// MaxTokenLife define the max lifetime of token (2 months)
	MaxTokenLife = time.Hour * 24 * 60
	// AuthCookieName defines name of cookie, which contains token
	AuthCookieName = "auth"
)

var (
	// Port for website
	Port string
	// IsTLS defines should the program use https
	IsTLS bool
	// Login for login
	Login string
	// Password for login
	Password string
	// Debug defines is debug mode
	Debug bool
	// Encrypt defines, should the program encrypt files. False by default
	Encrypt bool
	// Key is used for encrypting of files. Key is a sha256 sum of Password
	Key [32]byte
)

func init() {
	// Default - :80
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

	// Default - true
	IsTLS = func() bool {
		value := os.Getenv("TLS")
		if strings.ToLower(value) == "false" {
			return false
		}
		return true
	}()

	// Default - user
	Login = func() (login string) {
		login = os.Getenv("LOGIN")
		if login == "" {
			return "user"
		}
		return
	}()

	// Default - "qwerty"
	Password = func() (pswrd string) {
		pswrd = os.Getenv("PSWRD")
		if pswrd == "" {
			return "qwerty"
		}
		return
	}()

	// Default - false
	Debug = func() bool {
		value := os.Getenv("DBG")
		if strings.ToLower(value) == "true" {
			return true
		}
		return false
	}()

	// Default - false
	Encrypt = func() bool {
		enc := os.Getenv("ENCRYPT")
		if enc == "true" {
			return true
		}
		return false
	}()

	Key = sha256.Sum256([]byte(Password))
}
