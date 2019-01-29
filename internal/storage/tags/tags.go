package tags

import (
	clog "github.com/ShoshinNikita/log/v2"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
)

// storage for tags metadata
type storage interface {
	init() error

	// getAll returns all tags
	getAll() cmd.Tags

	// addTag adds a new tag
	addTag(tag cmd.Tag)

	// updateTag updates name and color of tag with id == tagID
	updateTag(id int, newName, newColor string)

	// deleteTag deletes a tag
	deleteTag(id int)

	// check returns true, if there's tag with passed it, else - false
	check(id int) bool

	shutdown() error
}

// TagStorage exposes methods for interactions with files
type TagStorage struct {
	storage storage
	logger  *clog.Logger
}

// NewTagStorage creates new FileStorage
func NewTagStorage(lg *clog.Logger) (*TagStorage, error) {
	var st storage

	switch params.StorageType {
	case params.JSONStorage:
		st = newJsonTagStorage(lg)
	default:
		st = newJsonTagStorage(lg)
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

func (ts TagStorage) GetAll() cmd.Tags {
	return ts.storage.getAll()
}

func (ts TagStorage) Add(name, color string) {
	t := cmd.Tag{Name: name, Color: color}
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
	return ts.storage.shutdown()
}
