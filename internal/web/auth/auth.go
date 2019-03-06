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
	"github.com/tags-drive/core/internal/params"
)

const maxTokenSize = 30

type Auth struct {
	tokens []tokenStruct // we can use array instead of map because number of tokens is small and O(n) is enough
	mutex  *sync.RWMutex

	// this channel signals that Auth.Shutdown() function was called
	shutdowned chan struct{}

	logger *clog.Logger
}

type tokenStruct struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expire"`
}

// NewAuthService create new Auth and inits tokens
func NewAuthService(lg *clog.Logger) (*Auth, error) {
	service := &Auth{
		mutex:      new(sync.RWMutex),
		logger:     lg,
		shutdowned: make(chan struct{}),
	}

	f, err := os.Open(params.TokensFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "can't open file %s", params.TokensFile)
		}

		// Have to create a new file
		err = service.createNewFile()
		if err != nil {
			return nil, err
		}
	} else {
		if !params.Encrypt {
			err = json.NewDecoder(f).Decode(&service.tokens)
		} else {
			// Have to decrypt at first
			buff := bytes.NewBuffer([]byte{})
			_, err := sio.Decrypt(buff, f, sio.Config{Key: params.PassPhrase[:]})
			if err == nil {
				err = json.NewDecoder(buff).Decode(&service.tokens)
			}
		}

		if err != nil {
			return nil, errors.Wrap(err, "can't decode allToken.tokens")
		}

		f.Close()
	}

	// Expiration function
	go func() {
		// Check tokens right now
		lg.Infoln("check expired tokens")
		service.expire()

		ticker := time.NewTicker(time.Hour * 6)
		for {
			select {
			case <-ticker.C:
				lg.Infoln("check expired tokens")
				service.expire()
			case <-service.shutdowned:
				ticker.Stop()
				return
			}
		}
	}()

	return service, nil
}

func (a Auth) createNewFile() error {
	a.logger.Infof("file %s doesn't exist. Need to create a new file\n", params.TokensFile)

	f, err := os.OpenFile(params.TokensFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	defer f.Close()

	// Write empty structure
	if !params.Encrypt {
		return json.NewEncoder(f).Encode(a.tokens)
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

	return err
}

// GenerateToken generates a new token
func (a Auth) GenerateToken() string {
	return generate(maxTokenSize)
}

// AddToken adds new generated token
func (a *Auth) AddToken(token string) {
	a.add(token)
}

// DeleteToken deletes token
func (a *Auth) DeleteToken(token string) {
	a.delete(token)
}

// CheckToken return true, if there's a passed token
func (a Auth) CheckToken(token string) bool {
	return a.check(token)
}

func (a *Auth) Shutdown() error {
	close(a.shutdowned)

	return nil
}
