package files

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

// jsonStorage implements files.storage interface.
// It is a map (filename: FileInfo) with RWMutex
type jsonStorage struct {
	info  map[string]FileInfo
	mutex *sync.RWMutex
}

func (js jsonStorage) init() error {
	// Create folders
	err := os.MkdirAll(params.DataFolder, 0600)
	if err != nil {
		return errors.Wrapf(err, "can't create a folder %s", params.DataFolder)
	}

	err = os.MkdirAll(params.ResizedImagesFolder, 0600)
	if err != nil {
		return errors.Wrapf(err, "can't create a folder %s", params.ResizedImagesFolder)
	}

	f, err := os.OpenFile(params.Files, os.O_RDWR, 0600)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			log.Infof("File %s doesn't exist. Need to create a new file\n", params.Files)
			f, err = os.OpenFile(params.Files, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return errors.Wrap(err, "can't create a new file")
			}
			// Write empty structure
			f.Write([]byte("{}"))
			// Can exit because we don't need to decode files from the file
			f.Close()
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.Files)
	}

	defer f.Close()

	return js.decode(f)
}

// write writes js.info into params.TagsFile
func (js jsonStorage) write() {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	f, err := os.OpenFile(params.Files, os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("Can't open file %s: %s\n", params.Files, err)
		return
	}

	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(js.info)

	f.Close()
}

// decode decodes js.info
func (js *jsonStorage) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&js.info)
}

func (js jsonStorage) getFile(filename string) (FileInfo, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	f, ok := js.info[filename]
	if !ok {
		return FileInfo{}, ErrFileIsNotExist
	}
	return f, nil
}

// getFiles returns slice of FileInfo with passed tags. If tags is an empty slice, function will return all files
func (js jsonStorage) getFiles(m TagMode, tags []int, search string) (files []FileInfo) {
	js.mutex.RLock()
	if len(tags) == 0 {
		files = make([]FileInfo, len(js.info))
		i := 0
		for _, v := range js.info {
			files[i] = v
			i++
		}
	} else {
		for _, v := range js.info {
			if isGoodFile(m, v.Tags, tags) {
				files = append(files, v)
			}
		}
	}

	js.mutex.RUnlock()

	if search == "" {
		return files
	}

	// Need to remove files with incorrect name
	var goodFiles []FileInfo
	for _, f := range files {
		if strings.Contains(f.Filename, search) {
			goodFiles = append(goodFiles, f)
		}
	}

	return goodFiles
}

// addFile adds an element into js.info and call js.write()
func (js *jsonStorage) addFile(info FileInfo) error {
	js.mutex.Lock()

	if _, ok := js.info[info.Filename]; ok {
		js.mutex.Unlock()
		return ErrAlreadyExist
	}

	info.Tags = []int{} // https://github.com/ShoshinNikita/tags-drive/issues/19
	js.info[info.Filename] = info
	js.mutex.Unlock()

	js.write()

	return nil
}

// renameFile renames a file
func (js *jsonStorage) renameFile(oldName string, newName string) error {
	js.mutex.Lock()
	if _, ok := js.info[oldName]; !ok {
		js.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Check does file with new name exist
	if _, ok := js.info[newName]; ok {
		js.mutex.Unlock()
		return ErrAlreadyExist
	}

	// Update map
	f := js.info[oldName]
	delete(js.info, oldName)
	f.Filename = newName
	f.Origin = params.DataFolder + "/" + newName
	js.info[newName] = f

	js.mutex.Unlock()

	js.write()

	return nil
}

// deleteFile deletes an element (from structure) and call js.write()
func (js *jsonStorage) deleteFile(filename string) error {
	js.mutex.Lock()

	if _, ok := js.info[filename]; !ok {
		js.mutex.Unlock()
		return ErrFileIsNotExist
	}

	delete(js.info, filename)

	js.mutex.Unlock()

	js.write()

	return nil
}

func (js *jsonStorage) deleteTagFromFiles(tagID int) {
	js.mutex.Lock()

	for filename, f := range js.info {
		index := -1
		for i := range f.Tags {
			if f.Tags[i] == tagID {
				index = i
				break
			}
		}
		if index == -1 {
			continue
		}
		// Erase tag
		f.Tags = append(f.Tags[0:index], f.Tags[index+1:]...)

		js.info[filename] = f
	}

	js.mutex.Unlock()

	js.write()
}

func (js *jsonStorage) updateFileTags(filename string, changedTagsID []int) error {
	js.mutex.Lock()

	if _, ok := js.info[filename]; !ok {
		js.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Update map
	f := js.info[filename]
	f.Tags = changedTagsID
	js.info[filename] = f

	js.mutex.Unlock()

	js.write()

	return nil
}

func (js *jsonStorage) updateFileDescription(filename string, newDesc string) error {
	js.mutex.Lock()

	if _, ok := js.info[filename]; !ok {
		js.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Update map
	f := js.info[filename]
	f.Description = newDesc
	js.info[filename] = f

	js.mutex.Unlock()

	js.write()

	return nil
}
