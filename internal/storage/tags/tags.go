package tags

import (
	clog "github.com/ShoshinNikita/log/v2"
	"github.com/pkg/errors"
)

// storage is an internal storage for tags metadata
type storage interface {
	init() error

	// getAll returns all tags
	getAll() Tags

	// addTag adds a new tag
	addTag(tag Tag)

	// updateTag updates name and color of tag with id == tagID
	updateTag(id int, newName, newColor string) (Tag, error)

	// updateGroup updates only group of a tag
	updateGroup(id int, newGroup string) (Tag, error)

	// deleteTag deletes a tag
	deleteTag(id int)

	// check returns true, if there's tag with passed it, else - false
	check(id int) bool

	shutdown() error
}

// TagStorage exposes methods for interactions with files
type TagStorage struct {
	config Config

	storage storage
	logger  *clog.Logger
}

// NewTagStorage creates new FileStorage
func NewTagStorage(cnf Config, lg *clog.Logger) (*TagStorage, error) {
	var st storage

	switch cnf.StorageType {
	case "json":
		fallthrough
	default:
		st = newJsonTagStorage(cnf, lg)
	}

	ts := &TagStorage{
		config:  cnf,
		storage: st,
		logger:  lg,
	}

	err := ts.storage.init()
	if err != nil {
		return nil, errors.Wrapf(err, "can't init tags storage")
	}

	return ts, nil
}

// Get returns Tag. If a tag doesn't exist, it returns Tag{}, false
func (ts TagStorage) Get(id int) (Tag, bool) {
	allTags := ts.GetAll()
	tag, ok := allTags[id]
	return tag, ok
}

func (ts TagStorage) GetAll() Tags {
	return ts.storage.getAll()
}

func (ts TagStorage) Add(name, color string) {
	t := Tag{Name: name, Color: color}
	ts.storage.addTag(t)
}

func (ts TagStorage) UpdateTag(id int, newName, newColor string) (updatedTag Tag, err error) {
	return ts.storage.updateTag(id, newName, newColor)
}

func (ts TagStorage) UpdateGroup(id int, newGroup string) (updatedTag Tag, err error) {
	return ts.storage.updateGroup(id, newGroup)
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
