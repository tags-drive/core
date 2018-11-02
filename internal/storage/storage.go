package storage

import (
	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
)

// Init calls files.Init() and tags.Init()
func Init() error {
	err := files.Init()
	if err != nil {
		return err
	}

	return tags.Init()
}
