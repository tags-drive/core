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

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

// FileInfo contains the information about a file
type FileInfo struct {
	Filename    string    `json:"filename"`
	Type        string    `json:"type"`
	Origin      string    `json:"origin"` // Origin is a path to a file (params.DataFolder/filename)
	Description string    `json:"description"`
	Size        int64     `json:"size"`
	Tags        []int     `json:"tags"`
	AddTime     time.Time `json:"addTime"`

	// Only if Type == TypeImage
	Preview string `json:"preview,omitempty"` // Preview is a path to a resized image
}

// jsonStorage implements files.storage interface.
// It is a map (filename: FileInfo) with RWMutex
type jsonStorage struct {
	info  map[string]FileInfo
	mutex *sync.RWMutex
}

func (fs jsonStorage) init() error {
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

	return fs.decode(f)
}

// write writes fs.info into params.TagsFile
func (fs jsonStorage) write() {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	f, err := os.OpenFile(params.Files, os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("Can't open file %s: %s\n", params.Files, err)
		return
	}

	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(fs.info)

	f.Close()
}

// decode decodes fs.info
func (fs *jsonStorage) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&fs.info)
}

/* Files */

func (fs jsonStorage) getFile(filename string) (FileInfo, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	f, ok := fs.info[filename]
	if !ok {
		return FileInfo{}, ErrFileIsNotExist
	}
	return f, nil
}

// getFiles returns slice of FileInfo with passed tags. If tags is an empty slice, function will return all files
func (fs jsonStorage) getFiles(m TagMode, tags []int, search string) (files []FileInfo) {
	fs.mutex.RLock()
	if len(tags) == 0 {
		files = make([]FileInfo, len(fs.info))
		i := 0
		for _, v := range fs.info {
			files[i] = v
			i++
		}
	} else {
		for _, v := range fs.info {
			if isGoodFile(m, v.Tags, tags) {
				files = append(files, v)
			}
		}
	}

	fs.mutex.RUnlock()

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

// addFile adds an element into fs.info and call fs.write()
func (fs *jsonStorage) addFile(info FileInfo) error {
	fs.mutex.Lock()

	if _, ok := fs.info[info.Filename]; ok {
		fs.mutex.Unlock()
		return ErrAlreadyExist
	}

	info.Tags = []int{} // https://github.com/ShoshinNikita/tags-drive/issues/19
	fs.info[info.Filename] = info
	fs.mutex.Unlock()

	fs.write()

	return nil
}

// renameFile renames a file
func (fs *jsonStorage) renameFile(oldName string, newName string) error {
	fs.mutex.Lock()
	if _, ok := fs.info[oldName]; !ok {
		fs.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Check does file with new name exist
	if _, ok := fs.info[newName]; ok {
		fs.mutex.Unlock()
		return ErrAlreadyExist
	}

	// Update map
	f := fs.info[oldName]
	delete(fs.info, oldName)
	f.Filename = newName
	f.Origin = params.DataFolder + "/" + newName
	fs.info[newName] = f

	// We have to unlock mutex after renaming, in order to user can't get invalid file
	err := os.Rename(params.DataFolder+"/"+oldName, params.DataFolder+"/"+newName)
	fs.mutex.Unlock()
	if err != nil {
		return err
	}

	fs.write()

	return nil
}

// deleteFile deletes an element (from structure) and call fs.write()
func (fs *jsonStorage) deleteFile(filename string) error {
	fs.mutex.Lock()

	if _, ok := fs.info[filename]; !ok {
		fs.mutex.Unlock()
		return ErrFileIsNotExist
	}

	delete(fs.info, filename)

	fs.mutex.Unlock()

	fs.write()

	return nil
}

func (fs *jsonStorage) deleteTagFromFiles(tagID int) {
	fs.mutex.Lock()

	for filename, f := range fs.info {
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

		fs.info[filename] = f
	}

	fs.mutex.Unlock()

	fs.write()
}

func (fs *jsonStorage) updateFileTags(filename string, changedTagsID []int) error {
	fs.mutex.Lock()

	if _, ok := fs.info[filename]; !ok {
		fs.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Update map
	f := fs.info[filename]
	f.Tags = changedTagsID
	fs.info[filename] = f

	fs.mutex.Unlock()

	fs.write()

	return nil
}

func (fs *jsonStorage) updateFileDescription(filename string, newDesc string) error {
	fs.mutex.Lock()

	if _, ok := fs.info[filename]; !ok {
		fs.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Update map
	f := fs.info[filename]
	f.Description = newDesc
	fs.info[filename] = f

	fs.mutex.Unlock()

	fs.write()

	return nil
}
