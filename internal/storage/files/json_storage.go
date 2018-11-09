package files

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/params"
)

// jsonFileStorage implements files.storage interface.
// It is a map (filename: FileInfo) with RWMutex
type jsonFileStorage struct {
	info  map[string]FileInfo
	mutex *sync.RWMutex
}

func (jfs jsonFileStorage) init() error {
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

	return jfs.decode(f)
}

// write writes js.info into params.TagsFile
func (jfs jsonFileStorage) write() {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	f, err := os.OpenFile(params.Files, os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("Can't open file %s: %s\n", params.Files, err)
		return
	}

	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(jfs.info)

	f.Close()
}

// decode decodes js.info
func (jfs *jsonFileStorage) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&jfs.info)
}

// checkFile return true if file exists
func (jfs jsonFileStorage) checkFile(filename string) bool {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	_, ok := jfs.info[filename]
	return ok
}

func (jfs jsonFileStorage) getFile(filename string) (FileInfo, error) {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	f, ok := jfs.info[filename]
	if !ok {
		return FileInfo{}, ErrFileIsNotExist
	}
	return f, nil
}

// getFiles returns slice of FileInfo with passed tags. If tags is an empty slice, function will return all files
func (jfs jsonFileStorage) getFiles(m TagMode, tags []int, search string) (files []FileInfo) {
	jfs.mutex.RLock()
	if len(tags) == 0 {
		files = make([]FileInfo, len(jfs.info))
		i := 0
		for _, v := range jfs.info {
			files[i] = v
			i++
		}
	} else {
		for _, v := range jfs.info {
			if isGoodFile(m, v.Tags, tags) {
				files = append(files, v)
			}
		}
	}

	jfs.mutex.RUnlock()

	if search == "" {
		return files
	}

	// Need to remove files with incorrect name
	var goodFiles []FileInfo
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Filename), search) {
			goodFiles = append(goodFiles, f)
		}
	}

	return goodFiles
}

// addFile adds an element into js.info and call js.write()
func (jfs *jsonFileStorage) addFile(info FileInfo) error {
	if jfs.checkFile(info.Filename) {
		return ErrAlreadyExist
	}

	jfs.mutex.Lock()

	info.Tags = []int{} // https://github.com/tags-drive/core/issues/19
	jfs.info[info.Filename] = info

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

// renameFile renames a file
func (jfs *jsonFileStorage) renameFile(oldName string, newName string) error {
	if !jfs.checkFile(oldName) {
		return ErrFileIsNotExist
	}

	// Check does file with new name exist
	if jfs.checkFile(newName) {
		return ErrAlreadyExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.info[oldName]
	delete(jfs.info, oldName)
	f.Filename = newName
	f.Origin = params.DataFolder + "/" + newName
	jfs.info[newName] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

func (jfs *jsonFileStorage) updateFileTags(filename string, changedTagsID []int) error {
	if !jfs.checkFile(filename) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.info[filename]
	f.Tags = changedTagsID
	jfs.info[filename] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

func (jfs *jsonFileStorage) updateFileDescription(filename string, newDesc string) error {
	if !jfs.checkFile(filename) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.info[filename]
	f.Description = newDesc
	jfs.info[filename] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

// deleteFile sets Deleted = true and update TimeToDelete
func (jfs *jsonFileStorage) deleteFile(filename string) error {
	if !jfs.checkFile(filename) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	f := jfs.info[filename]
	if f.Deleted {
		jfs.mutex.Unlock()
		return ErrFileDeletedAgain
	}

	f.Deleted = true
	f.TimeToDelete = time.Now().Add(timeBeforeDeleting)
	jfs.info[filename] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

// deleteFile deletes an element (from structure) and call js.write()
func (jfs *jsonFileStorage) deleteFileForce(filename string) error {
	if !jfs.checkFile(filename) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	delete(jfs.info, filename)

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

// recover sets Deleted = false
func (jfs *jsonFileStorage) recover(filename string) {
	if !jfs.checkFile(filename) {
		return
	}

	jfs.mutex.Lock()

	if !jfs.info[filename].Deleted {
		return
	}

	f := jfs.info[filename]
	f.Deleted = false
	f.TimeToDelete = time.Time{}
	jfs.info[filename] = f

	jfs.mutex.Unlock()

	jfs.write()
}

func (jfs *jsonFileStorage) deleteTagFromFiles(tagID int) {
	jfs.mutex.Lock()

	for filename, f := range jfs.info {
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

		jfs.info[filename] = f
	}

	jfs.mutex.Unlock()

	jfs.write()
}

// getExpiredDeletedFiles returns names of files with expired TimeToDelete
func (jfs *jsonFileStorage) getExpiredDeletedFiles() []string {
	jfs.mutex.RLock()

	var filesForDeleting []string
	now := time.Now()
	for _, file := range jfs.info {
		if file.Deleted && file.TimeToDelete.Before(now) {
			filesForDeleting = append(filesForDeleting, file.Filename)
		}
	}

	jfs.mutex.RUnlock()

	return filesForDeleting
}
