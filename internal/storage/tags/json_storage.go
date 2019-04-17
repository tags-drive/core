package tags

import (
	"bytes"
	"io"
	"os"
	"sync"

	clog "github.com/ShoshinNikita/log/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/sio"
	"github.com/pkg/errors"
)

type jsonTagStorage struct {
	config Config

	tags  Tags
	mutex *sync.RWMutex

	logger *clog.Logger
	json   jsoniter.API
}

func newJsonTagStorage(cnf Config, lg *clog.Logger) *jsonTagStorage {
	return &jsonTagStorage{
		config: cnf,
		tags:   make(Tags),
		mutex:  new(sync.RWMutex),
		logger: lg,
		json:   jsoniter.ConfigCompatibleWithStandardLibrary,
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

	return jts.decode(f)
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

	if !jts.config.Encrypt {
		// Encode directly into the file
		enc := jts.json.NewEncoder(f)
		if jts.config.Debug {
			enc.SetIndent("", "  ")
		}
		err := enc.Encode(jts.tags)
		if err != nil {
			jts.logger.Warnf("can't write '%s': %s", jts.config.TagsJSONFile, err)
		}

		return
	}

	// Encode into buffer
	buff := bytes.NewBuffer([]byte{})
	enc := jts.json.NewEncoder(buff)
	if jts.config.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(jts.tags)

	// Write into the file (jts.config.Encrypt is true, if we are here)
	_, err = sio.Encrypt(f, buff, sio.Config{Key: jts.config.PassPhrase[:]})

	if err != nil {
		jts.logger.Warnf("can't write '%s': %s\n", jts.config.TagsJSONFile, err)
	}
}

func (jts *jsonTagStorage) decode(r io.Reader) error {
	if !jts.config.Encrypt {
		return jts.json.NewDecoder(r).Decode(&jts.tags)
	}

	// Have to decrypt at first
	buff := bytes.NewBuffer([]byte{})
	_, err := sio.Decrypt(buff, r, sio.Config{Key: jts.config.PassPhrase[:]})
	if err != nil {
		return err
	}

	return jts.json.NewDecoder(buff).Decode(&jts.tags)
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
