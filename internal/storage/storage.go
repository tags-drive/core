package storage

import (
	"github.com/ShoshinNikita/tags-drive/internal/storage/files"
	"github.com/ShoshinNikita/tags-drive/internal/storage/tags"
)

// Init calls files.Init() and tags.Init()
func Init() error {
	err := files.Init()
	if err != nil {
		return err
	}

	return tags.Init()
}
