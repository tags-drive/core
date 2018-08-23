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

type Tag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type tagsStruct struct {
	tags  []Tag
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

func (t tagsStruct) getAll() []Tag {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tags
}

func (t *tagsStruct) add(tag Tag) {
	t.mutex.Lock()
	t.tags = append(t.tags, tag)
	t.mutex.Unlock()

	t.write()
}

func (t *tagsStruct) delete(name string) {
	t.mutex.Lock()
	index := -1
	for i, tag := range t.tags {
		if tag.Name == name {
			index = i
			break
		}
	}
	if index == -1 {
		t.mutex.Unlock()
		return
	}

	t.tags = append(t.tags[0:index], t.tags[index+1:]...)
	t.mutex.Unlock()

	t.write()
}

func (t *tagsStruct) change(tag, newName, newColor string) {
	t.mutex.Lock()

	index := -1
	for i := range t.tags {
		if t.tags[i].Name == tag {
			index = i
			break
		}
	}
	if index == -1 {
		t.mutex.Unlock()
		return
	}

	if newName != "" {
		t.tags[index].Name = newName
	}

	if newColor != "" {
		t.tags[index].Color = newColor
	}

	t.mutex.Unlock()

	t.write()
}

var allTags = tagsStruct{mutex: new(sync.RWMutex)}

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

func GetAllTags() []Tag {
	return allTags.getAll()
}

// GetTags returns tags with passed names
func GetTags(names []string) []Tag {
	tags := allTags.getAll()
	var res []Tag
	for _, t := range tags {
		for _, n := range names {
			if t.Name == n {
				res = append(res, t)
				break
			}
		}
	}

	return res
}

func AddTag(t Tag) {
	allTags.add(t)
}

func DeleteTag(name string) {
	allTags.delete(name)
}

// Change changes a tag with Name name.
// If pass empty newName (or newColor), field Name (or Color) won't be changed.
func Change(name, newName, newColor string) {
	allTags.change(name, newName, newColor)
}
