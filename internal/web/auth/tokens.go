package auth

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"

	"github.com/tags-drive/core/internal/params"
)

type tokenStruct struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expire"`
}

type tokens struct {
	tokens []tokenStruct // we can use []tokenStruct instead of map, because number of tokens is small and O(n) also isn't huge
	mutex  *sync.RWMutex
}

func (t tokens) write() {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	f, err := os.OpenFile(params.TokensFile, os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("Can't open file %s: %s\n", params.TokensFile, err)
		return
	}

	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(t.tokens)

	f.Close()
}

func (t *tokens) add(token string) {
	t.mutex.Lock()
	t.tokens = append(t.tokens, tokenStruct{Token: token, Expires: time.Now().Add(params.MaxTokenLife)})
	t.mutex.Unlock()

	t.write()
}

func (t *tokens) delete(token string) {
	t.mutex.Lock()

	tokenIndex := -1
	for i, tok := range t.tokens {
		if tok.Token == token {
			tokenIndex = i
			break
		}
	}
	if tokenIndex == -1 {
		t.mutex.Unlock()
		return
	}

	t.tokens = append(t.tokens[:tokenIndex], t.tokens[tokenIndex+1:]...)

	t.mutex.Unlock()

	t.write()
}

func (t tokens) check(token string) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, tok := range t.tokens {
		if tok.Token == token {
			return true
		}
	}

	return false
}

// expire removes expired tokens
func (t *tokens) expire() {
	t.mutex.Lock()

	var freshTokens []tokenStruct
	now := time.Now()
	for _, tok := range t.tokens {
		if now.Before(tok.Expires) {
			freshTokens = append(freshTokens, tok)
		}
	}

	t.tokens = freshTokens

	t.mutex.Unlock()

	t.write()
}
