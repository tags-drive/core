package tags

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"
	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files"
)

type jsonTagStorage struct {
	tags  Tags
	mutex *sync.RWMutex
}

func (jts jsonTagStorage) init() error {
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
	return jts.decode(f)
}

func (jts jsonTagStorage) write() {
	jts.mutex.RLock()
	defer jts.mutex.RUnlock()

	f, err := os.OpenFile(params.TagsFile, os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("Can't open file %s: %s\n", params.TagsFile, err)
		return
	}

	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(jts.tags)

	f.Close()
}

func (jts *jsonTagStorage) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&jts.tags)
}

func (jts jsonTagStorage) getAll() Tags {
	jts.mutex.RLock()
	defer jts.mutex.RUnlock()

	return jts.tags
}

func (jts *jsonTagStorage) addTag(tag Tag) {
	jts.mutex.Lock()

	// Get max ID (max)
	nextID := 0
	for id := range jts.tags {
		if nextID < id {
			nextID = id
		}
	}
	nextID++
	tag.ID = nextID
	jts.tags[nextID] = tag

	jts.mutex.Unlock()

	jts.write()
}

func (jts *jsonTagStorage) updateTag(id int, newName, newColor string) {
	jts.mutex.Lock()

	if _, ok := jts.tags[id]; !ok {
		jts.mutex.Unlock()
		return
	}

	tag := jts.tags[id]

	if newName != "" {
		tag.Name = newName
	}

	if newColor != "" {
		if newColor[0] != '#' {
			newColor = "#" + newColor
		}
		tag.Color = newColor
	}

	jts.tags[id] = tag

	jts.mutex.Unlock()

	jts.write()
}

func (jts *jsonTagStorage) deleteTag(id int) {
	jts.mutex.Lock()
	// We can skip files.DeleteTag(id), if tag doesn't exist
	if _, ok := jts.tags[id]; !ok {
		jts.mutex.Unlock()
		return
	}

	delete(jts.tags, id)
	jts.mutex.Unlock()

	jts.write()

	// Remove links to deleted tag
	files.DeleteTag(id)
}
