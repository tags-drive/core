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

// TagStorage exposes methods for interactions with files
type TagStorage struct {
	storage storage
}

// Init inits ts.storage
func (ts *TagStorage) Init() error {
	switch params.StorageType {
	case params.JSONStorage:
		ts.storage = &jsonTagStorage{
			tags:  make(Tags),
			mutex: new(sync.RWMutex),
		}
	default:
		// Default storage is jsonTagStorage
		ts.storage = &jsonTagStorage{
			tags:  make(Tags),
			mutex: new(sync.RWMutex),
		}
	}

	err := ts.storage.init()
	if err != nil {
		return errors.Wrapf(err, "can't decode data")
	}

	return nil
}

func (ts TagStorage) GetAll() Tags {
	return ts.storage.getAll()
}

func (ts TagStorage) Add(name, color string) {
	t := Tag{Name: name, Color: color}
	ts.storage.addTag(t)
}

func (ts TagStorage) Delete(id int) {
	ts.storage.deleteTag(id)
}

func (ts TagStorage) Change(id int, newName, newColor string) {
	ts.storage.updateTag(id, newName, newColor)
}
