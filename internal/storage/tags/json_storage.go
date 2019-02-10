package tags

import (
	"io"
	"os"
	"sync"

	clog "github.com/ShoshinNikita/log/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
)

type jsonTagStorage struct {
	tags  cmd.Tags
	mutex *sync.RWMutex

	logger *clog.Logger
	json   jsoniter.API
}

func newJsonTagStorage(lg *clog.Logger) *jsonTagStorage {
	return &jsonTagStorage{
		tags:   make(cmd.Tags),
		mutex:  new(sync.RWMutex),
		logger: lg,
		json:   jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (jts *jsonTagStorage) init() error {
	f, err := os.OpenFile(params.TagsFile, os.O_RDWR, 0666)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			jts.logger.Infof("file %s doesn't exist. Need to create a new file\n", params.TagsFile)
			f, err = os.OpenFile(params.TagsFile, os.O_CREATE|os.O_RDWR, 0666)
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

	f, err := os.OpenFile(params.TagsFile, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		jts.logger.Errorf("can't open file %s: %s\n", params.TagsFile, err)
		return
	}

	enc := jts.json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(jts.tags)

	f.Close()
}

func (jts *jsonTagStorage) decode(r io.Reader) error {
	return jts.json.NewDecoder(r).Decode(&jts.tags)
}

func (jts jsonTagStorage) getAll() cmd.Tags {
	jts.mutex.RLock()
	defer jts.mutex.RUnlock()

	return jts.tags
}

func (jts *jsonTagStorage) addTag(tag cmd.Tag) {
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

func (jts *jsonTagStorage) updateTag(id int, newName, newColor string) (cmd.Tag, error) {
	jts.mutex.Lock()

	if _, ok := jts.tags[id]; !ok {
		jts.mutex.Unlock()
		return cmd.Tag{}, errors.New("tag doesn't exist")
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

	return tag, nil
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
}

func (jts jsonTagStorage) check(id int) bool {
	jts.mutex.RLock()
	defer jts.mutex.RUnlock()

	_, ok := jts.tags[id]
	return ok
}

func (jts jsonTagStorage) shutdown() error {
	// We have not to do any special operations because we update json file on every change.
	// Also there are no any requests because server is already down. But it's better to check the mutex
	// just in case.

	jts.mutex.Lock()
	jts.mutex.Unlock()

	// There will be no any new requests.

	return nil
}
