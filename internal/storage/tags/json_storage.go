package tags

import (
	"os"
	"sync"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/utils"
)

type jsonTagStorage struct {
	config Config

	tags  Tags
	mutex *sync.RWMutex

	logger *clog.Logger
}

func newJsonTagStorage(cnf Config, lg *clog.Logger) *jsonTagStorage {
	return &jsonTagStorage{
		config: cnf,
		tags:   make(Tags),
		mutex:  new(sync.RWMutex),
		logger: lg,
	}
}

func (jts *jsonTagStorage) init() error {
	f, err := os.OpenFile(jts.config.TagsJSONFile, os.O_RDWR, 0666)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			return jts.createNewFile()
		}

		return errors.Wrapf(err, "can't open file %s", jts.config.TagsJSONFile)
	}
	defer f.Close()

	return utils.Decode(f, &jts.tags, jts.config.Encrypt, jts.config.PassPhrase)
}

func (jts *jsonTagStorage) createNewFile() error {
	jts.logger.Debugf("file %s doesn't exist. Need to create a new file\n", jts.config.TagsJSONFile)

	// Just create a new file
	f, err := os.OpenFile(jts.config.TagsJSONFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	f.Close()

	// Write empty tag map
	jts.write()

	return nil
}

func (jts jsonTagStorage) write() {
	jts.mutex.RLock()
	defer jts.mutex.RUnlock()

	f, err := os.OpenFile(jts.config.TagsJSONFile, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		jts.logger.Errorf("can't open file %s: %s\n", jts.config.TagsJSONFile, err)
		return
	}
	defer f.Close()

	err = utils.Encode(f, jts.tags, jts.config.Encrypt, jts.config.PassPhrase)
	if err != nil {
		jts.logger.Warnf("can't write '%s': %s", jts.config.TagsJSONFile, err)
	}
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

func (jts *jsonTagStorage) updateTag(id int, newName, newColor string) (Tag, error) {
	jts.mutex.Lock()

	if _, ok := jts.tags[id]; !ok {
		jts.mutex.Unlock()
		return Tag{}, errors.New("tag doesn't exist")
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

func (jts *jsonTagStorage) updateGroup(id int, newGroup string) (Tag, error) {
	jts.mutex.Lock()

	if _, ok := jts.tags[id]; !ok {
		jts.mutex.Unlock()
		return Tag{}, errors.New("tag doesn't exist")
	}

	tag := jts.tags[id]

	tag.Group = newGroup

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
