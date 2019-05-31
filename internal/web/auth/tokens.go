package auth

import (
	"os"
	"time"

	"github.com/tags-drive/core/internal/utils"
)

func (a Auth) write() {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	f, err := os.OpenFile(a.config.TokensJSONFile, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		a.logger.Errorf("can't open file %s: %s\n", a.config.TokensJSONFile, err)
		return
	}
	defer f.Close()

	err = utils.Encode(f, a.tokens, a.config.Encrypt, a.config.PassPhrase)
	if err != nil {
		a.logger.Warnf("can't write '%s': %s\n", a.config.TokensJSONFile, err)
	}
}

func (a *Auth) add(token string) {
	a.mutex.Lock()
	a.tokens = append(a.tokens, tokenStruct{Token: token, Expires: time.Now().Add(a.config.MaxTokenLife)})
	a.mutex.Unlock()

	a.write()
}

func (a *Auth) delete(token string) {
	a.mutex.Lock()

	tokenIndex := -1
	for i, tok := range a.tokens {
		if tok.Token == token {
			tokenIndex = i
			break
		}
	}
	if tokenIndex == -1 {
		a.mutex.Unlock()
		return
	}

	a.tokens = append(a.tokens[:tokenIndex], a.tokens[tokenIndex+1:]...)

	a.mutex.Unlock()

	a.write()
}

func (a Auth) check(token string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, tok := range a.tokens {
		if tok.Token == token {
			return true
		}
	}

	return false
}

// expire removes expired tokens
func (a *Auth) expire() {
	a.mutex.Lock()

	freshTokens := []tokenStruct{}
	now := time.Now()
	for _, tok := range a.tokens {
		if now.Before(tok.Expires) {
			freshTokens = append(freshTokens, tok)
		} else {
			a.logger.Debugf("token \"%s\" expired\n", tok.Token)
		}
	}

	a.tokens = freshTokens

	a.mutex.Unlock()

	a.write()
}
