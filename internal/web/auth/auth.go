package auth

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"
	"github.com/tags-drive/core/internal/params"
)

// DefaultTokenSize can be used for function GenerateToken()
const DefaultTokenSize = 30

var allTokens = tokens{mutex: new(sync.RWMutex)}

// Init inits allTokens and run in goroutine function for token expiring
func Init() error {
	f, err := os.Open(params.TokensFile)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			log.Infof("File %s doesn't exist. Need to create a new file\n", params.TokensFile)
			f, err = os.OpenFile(params.TokensFile, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return errors.Wrap(err, "can't create a new file")
			}
			// Write empty structure
			json.NewEncoder(f).Encode(allTokens.tokens)
			// Can exit because we don't need to decode files from the file
			f.Close()
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.TokensFile)
	}

	err = json.NewDecoder(f).Decode(&allTokens.tokens)
	if err != nil {
		return errors.Wrap(err, "can't decode allToken.tokens")
	}

	// Expiration function
	go func() {
		ticker := time.NewTicker(time.Hour * 6)
		for ; true; <-ticker.C {
			log.Infoln("Check expired tokens")
			allTokens.expire()
		}
	}()

	return nil
}

// GenerateToken generates token with passed size
func GenerateToken() string {
	return generate(DefaultTokenSize)
}

// AddToken adds new generated token
func AddToken(token string) {
	allTokens.add(token)
}

// DeleteToken deletes token
func DeleteToken(token string) {
	allTokens.delete(token)
}

// CheckToken return true, if there's a passed token
func CheckToken(token string) bool {
	return allTokens.check(token)
}
