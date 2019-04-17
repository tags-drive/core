package auth

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/minio/sio"
	"github.com/tags-drive/core/internal/params"
)

func (a Auth) write() {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	f, err := os.OpenFile(params.TokensJSONFile, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		a.logger.Errorf("can't open file %s: %s\n", params.TokensJSONFile, err)
		return
	}
	defer f.Close()

	if !params.Encrypt {
		// Encode directly into the file
		enc := json.NewEncoder(f)
		if params.Debug {
			enc.SetIndent("", "  ")
		}
		err := enc.Encode(a.tokens)
		if err != nil {
			a.logger.Warnf("can't write '%s': %s\n", params.TokensJSONFile, err)
		}

		return
	}

	// Encode into buffer
	buff := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buff)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(a.tokens)

	// Write into the file (params.Encrypt is true, if we are here)
	_, err = sio.Encrypt(f, buff, sio.Config{Key: params.PassPhrase[:]})

	if err != nil {
		a.logger.Warnf("can't write '%s': %s\n", params.TokensJSONFile, err)
	}
}

func (a *Auth) add(token string) {
	a.mutex.Lock()
	a.tokens = append(a.tokens, tokenStruct{Token: token, Expires: time.Now().Add(params.MaxTokenLife)})
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

	var freshTokens []tokenStruct
	now := time.Now()
	for _, tok := range a.tokens {
		if now.Before(tok.Expires) {
			freshTokens = append(freshTokens, tok)
		} else {
			a.logger.Infof("token \"%s\" expired\n", tok.Token)
		}
	}

	a.tokens = freshTokens

	a.mutex.Unlock()

	a.write()
}
