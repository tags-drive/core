package utils

import (
	"math/rand"
	"time"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomString generates random string of a length n
func GenerateRandomString(n int) string {
	s := make([]byte, n)
	length := len(alphabet)

	for i := 0; i < n; i++ {
		s[i] = alphabet[rand.Intn(length)]
	}

	return string(s)
}
