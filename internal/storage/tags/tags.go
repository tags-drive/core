package tags

import (
	"sync"

	"github.com/ShoshinNikita/log"
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

	// check returns true, if there's tag with passed it, else - false
	check(id int) bool
}

// TagStorage exposes methods for interactions with files
type TagStorage struct {
	storage storage
	logger  *log.Logger
}

// NewTagStorage creates new FileStorage
func NewTagStorage(lg *log.Logger) (*TagStorage, error) {
	var st storage

	switch params.StorageType {
	case params.JSONStorage:
		st = &jsonTagStorage{
			tags:   make(Tags),
			mutex:  new(sync.RWMutex),
			logger: lg,
		}
	default:
		st = &jsonTagStorage{
			tags:   make(Tags),
			mutex:  new(sync.RWMutex),
			logger: lg,
		}
	}

	ts := &TagStorage{
		storage: st,
		logger:  lg,
	}

	err := ts.storage.init()
	if err != nil {
		return nil, errors.Wrapf(err, "can't init tags storage")
	}

	return ts, nil
}

func (ts TagStorage) GetAll() Tags {
	return ts.storage.getAll()
}

func (ts TagStorage) Add(name, color string) {
	t := Tag{Name: name, Color: color}
	ts.storage.addTag(t)
}

func (ts TagStorage) Change(id int, newName, newColor string) {
	ts.storage.updateTag(id, newName, newColor)
}

func (ts TagStorage) Delete(id int) {
	ts.storage.deleteTag(id)
}

func (ts TagStorage) Check(id int) bool {
	return ts.storage.check(id)
}

func (ts TagStorage) Shutdown() error {
	return nil
}
