package tags

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/ShoshinNikita/log"
	"github.com/ShoshinNikita/tags-drive/internal/params"
	"github.com/ShoshinNikita/tags-drive/internal/storage/files"
	"github.com/pkg/errors"
)

type jsonStorage struct {
	tags  Tags
	mutex *sync.RWMutex
}

func (t jsonStorage) init() error {
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
			f.Write([]byte("{}"))
			// Can exit because we don't need to decode tags from the file
			f.Close()
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.TagsFile)
	}

	defer f.Close()
	return t.decode(f)
}

func (t jsonStorage) write() {
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

func (t *jsonStorage) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&t.tags)
}

func (t jsonStorage) getAll() Tags {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tags
}

func (t *jsonStorage) addTag(tag Tag) {
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

func (t *jsonStorage) deleteTag(id int) {
	t.mutex.Lock()
	// We can skip files.DeleteTag(id), if tag doesn't exist
	if _, ok := t.tags[id]; !ok {
		t.mutex.Unlock()
		return
	}

	delete(t.tags, id)
	t.mutex.Unlock()

	t.write()

	files.DeleteTag(id)
}

func (t *jsonStorage) updateTag(id int, newName, newColor string) {
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
