package auth

import (
	"os"
	"sync"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/utils"
)

const maxTokenSize = 30

type AuthService struct {
	config Config

	tokens []tokenStruct // we can use array instead of map because number of tokens is small and O(n) is enough
	mutex  *sync.RWMutex

	// this channel signals that AuthService.Shutdown() function was called
	shutdowned chan struct{}

	logger *clog.Logger
}

type tokenStruct struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expire"`
}

// NewAuthService create a new AuthService and inits tokens
func NewAuthService(cnf Config, lg *clog.Logger) (*AuthService, error) {
	service := &AuthService{
		config:     cnf,
		mutex:      new(sync.RWMutex),
		logger:     lg,
		shutdowned: make(chan struct{}),
	}

	if f, err := os.Open(service.config.TokensJSONFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "can't open file %s", cnf.TokensJSONFile)
		}

		// Have to create a new file
		err := service.createNewFile()
		if err != nil {
			return nil, err
		}
	} else {
		err = utils.Decode(f, &service.tokens, cnf.Encrypt, cnf.PassPhrase)
		if err != nil {
			return nil, errors.Wrap(err, "can't decode allToken.tokens")
		}

		f.Close()
	}

	return service, nil
}

func (a AuthService) createNewFile() error {
	a.logger.Debugf("file %s doesn't exist. Need to create a new file\n", a.config.TokensJSONFile)

	f, err := os.OpenFile(a.config.TokensJSONFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	defer f.Close()

	return utils.Encode(f, a.tokens, a.config.Encrypt, a.config.PassPhrase)
}

// Start starts all background services
func (a *AuthService) StartBackgroundServices() {
	// Start expiration function
	go func() {
		// Check tokens right now
		a.logger.Debugln("check expired tokens")
		a.expire()

		ticker := time.NewTicker(time.Hour * 6)
		for {
			select {
			case <-ticker.C:
				a.logger.Debugln("check expired tokens")
				a.expire()
			case <-a.shutdowned:
				ticker.Stop()
				return
			}
		}
	}()
}

// expire removes expired tokens
func (a *AuthService) expire() {
	a.mutex.Lock()
	defer func() {
		a.mutex.Unlock()
		a.write()
	}()

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
}

func (a AuthService) write() {
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

// GenerateToken generates a new token. GenerateToken doesn't add new token, just return it!
func (a AuthService) GenerateToken() string {
	return utils.GenerateRandomString(maxTokenSize)
}

// AddToken adds passed token into storage
func (a *AuthService) AddToken(token string) {
	a.mutex.Lock()
	defer func() {
		a.mutex.Unlock()
		a.write()
	}()

	a.tokens = append(a.tokens, tokenStruct{Token: token, Expires: time.Now().Add(a.config.MaxTokenLife)})
}

// DeleteToken deletes token from a storage
func (a *AuthService) DeleteToken(token string) {
	a.mutex.Lock()
	defer func() {
		a.mutex.Unlock()
		a.write()
	}()

	tokenIndex := -1
	for i, tok := range a.tokens {
		if tok.Token == token {
			tokenIndex = i
			break
		}
	}
	if tokenIndex == -1 {
		return
	}

	a.tokens = append(a.tokens[:tokenIndex], a.tokens[tokenIndex+1:]...)
}

// CheckToken returns true if token is in storage
func (a AuthService) CheckToken(token string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, tok := range a.tokens {
		if tok.Token == token {
			return true
		}
	}

	return false
}

// Shutdown gracefully shutdown FileStorage
func (a *AuthService) Shutdown() error {
	// Wait for all locks
	a.mutex.Lock()
	a.mutex.Unlock()

	close(a.shutdowned)

	return nil
}
