package share

import (
	clog "github.com/ShoshinNikita/log/v2"
)

func NewShareStorage(cnf Config, lg *clog.Logger) (ShareStorageInterface, error) {
	storage := newJsonShareStorage(cnf, lg)

	err := storage.init()
	if err != nil {
		return nil, err
	}

	return storage, nil
}
