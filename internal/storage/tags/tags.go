package tags

import (
	"sync"

	"github.com/tags-drive/core/internal/params"
	"github.com/pkg/errors"
)

const (
	// DefaultColor is a white color
	DefaultColor = "#ffffff"
)

var (
	ErrTagIsNotExist = errors.New("tag doesn't exist")
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

var tagStorage = struct {
	storage
}{}

// Init inits tagStorage
func Init() error {
	switch params.StorageType {
	case params.JSONStorage:
		tagStorage.storage = &jsonTagStorage{
			tags:  make(Tags),
			mutex: new(sync.RWMutex),
		}
	default:
		// Default storage is jsonTagStorage
		tagStorage.storage = &jsonTagStorage{
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

func GetAllTags() Tags {
	return tagStorage.getAll()
}

func AddTag(t Tag) {
	tagStorage.addTag(t)
}

func DeleteTag(id int) {
	tagStorage.deleteTag(id)
}

// Change changes a tag with passed id.
// If pass empty newName (or newColor), field Name (or Color) won't be changed.
func Change(id int, newName, newColor string) {
	tagStorage.updateTag(id, newName, newColor)
}
