package share

import (
	clog "github.com/ShoshinNikita/log/v2"
	"github.com/pkg/errors"
)

var (
	ErrInvalidToken = errors.New("invalid share token")
)

func NewShareStorage(cnf Config, fs FileStorage, lg *clog.Logger) (ShareStorageInterface, error) {
	storage := newJsonShareStorage(cnf, fs, lg)

	err := storage.init()
	if err != nil {
		return nil, err
	}

	return storage, nil
}
