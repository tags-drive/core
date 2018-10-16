package tags

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/ShoshinNikita/log"
	"github.com/ShoshinNikita/tags-drive/internal/params"
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

var allTags = jsonStorage{tags: make(Tags), mutex: new(sync.RWMutex)}

// Init reads params.TagsFiles and decode its data
func Init() error {
	f, err := os.OpenFile(params.TagsFile, os.O_RDWR, 0600)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			log.Infof("File %s doesn't exist. Need to create a new file\n", params.TagsFile)
			f, err = os.OpenFile(params.TagsFile, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return errors.Wrap(err, "can't create a new file")
			}
			// Write empty structure
			json.NewEncoder(f).Encode(allTags.tags)
			// Can exit because we don't need to decode files from the file
			f.Close()
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.TagsFile)
	}

	defer f.Close()
	err = allTags.decode(f)
	if err != nil {
		return errors.Wrapf(err, "can't decode data")
	}

	return nil
}

func GetAllTags() Tags {
	return allTags.getAll()
}

func AddTag(t Tag) {
	allTags.addTag(t)
}

func DeleteTag(id int) {
	allTags.deleteTag(id)
}

// Change changes a tag with passed id.
// If pass empty newName (or newColor), field Name (or Color) won't be changed.
func Change(id int, newName, newColor string) {
	allTags.updateTag(id, newName, newColor)
}
