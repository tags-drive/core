package migrator

import (
	"io"
)

type readCloserWrapper struct {
	r io.Reader
}

func (rc readCloserWrapper) Read(p []byte) (n int, err error) {
	return rc.r.Read(p)
}

func (rc readCloserWrapper) Close() error {
	return nil
}
