package tags

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/tags-drive/core/internal/params"
)

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Tags map[int]Tag

type storage interface {
	init() error

	// getAll returns all tags
	getAll() Tags

	// addTag adds a new tag
	addTag(tag Tag)

	// updateTag updates name and color of tag with id == tagID
	updateTag(id int, newName, newColor string)

	// deleteTag deletes a tag
	deleteTag(id int)
}

var tagStorage storage

// TagStorage exposes methods for interactions with files
type TagStorage struct{}

// Init inits tagStorage
func (ts TagStorage) Init() error {
	switch params.StorageType {
	case params.JSONStorage:
		tagStorage = &jsonTagStorage{
			tags:  make(Tags),
			mutex: new(sync.RWMutex),
		}
	default:
		// Default storage is jsonTagStorage
		tagStorage = &jsonTagStorage{
			tags:  make(Tags),
			mutex: new(sync.RWMutex),
		}
	}

	err := tagStorage.init()
	if err != nil {
		return errors.Wrapf(err, "can't decode data")
	}

	return nil
}

func (ts TagStorage) GetAll() Tags {
	return tagStorage.getAll()
}

func (ts TagStorage) Add(name, color string) {
	t := Tag{Name: name, Color: color}
	tagStorage.addTag(t)
}

func (ts TagStorage) Delete(id int) {
	tagStorage.deleteTag(id)
}

func (ts TagStorage) Change(id int, newName, newColor string) {
	tagStorage.updateTag(id, newName, newColor)
}
