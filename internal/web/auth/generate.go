package auth

import (
	"math/rand"
	"time"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generate(n int) string {
	token := make([]byte, n)

	for i := 0; i < n; i++ {
		token[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(token)
}
