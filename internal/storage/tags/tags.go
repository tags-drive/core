package tags

import (
	"encoding/json"
	"io"
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

type tagsStruct struct {
	tags  Tags
	mutex *sync.RWMutex
}

func (t tagsStruct) write() {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	f, err := os.OpenFile(params.TagsFile, os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("Can't open file %s: %s\n", params.TagsFile, err)
		return
	}

	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(t.tags)

	f.Close()
}

func (t *tagsStruct) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&t.tags)
}

func (t tagsStruct) getAll() Tags {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tags
}

func (t *tagsStruct) add(tag Tag) {
	t.mutex.Lock()

	// Get max ID (max)
	nextID := 0
	for id := range t.tags {
		if nextID < id {
			nextID = id
		}
	}
	nextID++
	tag.ID = nextID
	t.tags[nextID] = tag
	t.mutex.Unlock()

	t.write()
}

func (t *tagsStruct) deleteTag(id int) {
	t.mutex.Lock()
	delete(t.tags, id)
	t.mutex.Unlock()

	t.write()
}

func (t *tagsStruct) change(id int, newName, newColor string) {
	t.mutex.Lock()

	if _, ok := t.tags[id]; !ok {
		t.mutex.Unlock()
		return
	}

	tag := t.tags[id]

	if newName != "" {
		tag.Name = newName
	}

	if newColor != "" {
		if newColor[0] != '#' {
			newColor = "#" + newColor
		}
		tag.Color = newColor
	}

	t.tags[id] = tag

	t.mutex.Unlock()

	t.write()
}

var allTags = tagsStruct{tags: make(Tags), mutex: new(sync.RWMutex)}

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
	allTags.add(t)
}

func DeleteTag(id int) {
	allTags.deleteTag(id)
}

// Change changes a tag with passed id.
// If pass empty newName (or newColor), field Name (or Color) won't be changed.
func Change(id int, newName, newColor string) {
	allTags.change(id, newName, newColor)
}
