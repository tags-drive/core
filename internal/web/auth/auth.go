package auth

import (
	"bytes"
	"encoding/json"
	"os"
	"sync"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/minio/sio"
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

func (a AuthService) createNewFile() error {
	a.logger.Debugf("file %s doesn't exist. Need to create a new file\n", a.config.TokensJSONFile)

	f, err := os.OpenFile(a.config.TokensJSONFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	defer f.Close()

	// Write empty structure
	if !a.config.Encrypt {
		return json.NewEncoder(f).Encode(a.tokens)
	}

	// Encode into buffer
	buff := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buff)
	if a.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(a.tokens)

	// Write into the file (params.Encrypt is true, if we are here)
	_, err = sio.Encrypt(f, buff, sio.Config{Key: a.config.PassPhrase[:]})

	return err
}

// GenerateToken generates a new token. GenerateToken doesn't add new token, just return it!
func (a AuthService) GenerateToken() string {
	return utils.GenerateRandomString(maxTokenSize)
}

// AddToken adds passed token into storage
func (a *AuthService) AddToken(token string) {
	a.add(token)
}

// DeleteToken deletes token from a storage
func (a *AuthService) DeleteToken(token string) {
	a.delete(token)
}

// CheckToken returns true if token is in storage
func (a AuthService) CheckToken(token string) bool {
	return a.check(token)
}

// Shutdown gracefully shutdown FileStorage
func (a *AuthService) Shutdown() error {
	close(a.shutdowned)

	return nil
}
