// Package params provides global vars getted from environment
package params

import (
	"os"
	"strings"
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
)

func init() {
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

	IsTLS = func() bool {
		value := os.Getenv("TLS")
		if strings.ToLower(value) == "false" {
			return false
		}
		return true
	}()

	Login = func() (login string) {
		login = os.Getenv("LOGIN")
		if login == "" {
			return "user"
		}
		return
	}()

	Password = func() (pswrd string) {
		pswrd = os.Getenv("PSWRD")
		if pswrd == "" {
			return "qwerty"
		}
		return
	}()

	Debug = func() bool {
		value := os.Getenv("DBG")
		if strings.ToLower(value) == "true" {
			return true
		}
		return false
	}()
}
