package storage

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

// FileInfo contains the information about a file
type FileInfo struct {
	Filename string    `json:"filename"`
	Size     int64     `json:"size"`
	Tags     []string  `json:"tags"`
	AddTime  time.Time `json:"add_time"`
}

// filesData is a map (filename: FileInfo) with RWMutex
// files.json keeps only filesData.info
type filesData struct {
	info  map[string]FileInfo
	mutex *sync.RWMutex
}

/* Persistent */

// write writes fs.info into params.TagsFile
func (fs filesData) write() {
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
func (fs *filesData) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&fs.info)
}

/* Files */

// add adds an element into fs.info and call fs.write()
func (fs *filesData) add(info FileInfo) error {
	fs.mutex.Lock()

	if _, ok := fs.info[info.Filename]; ok {
		fs.mutex.Unlock()
		return ErrAlreadyExist
	}

	fs.info[info.Filename] = info
	fs.mutex.Unlock()

	fs.write()

	return nil
}

// rename renames a file
func (fs *filesData) rename(oldName string, newName string) error {
	fs.mutex.Lock()
	if _, ok := fs.info[oldName]; !ok {
		fs.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Update map
	f := fs.info[oldName]
	delete(fs.info, oldName)
	f.Filename = newName
	f.AddTime = time.Now()
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

// delete deletes an element (from structure) and call fs.write()
func (fs *filesData) delete(filename string) error {
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

/* Tags */

func (fs *filesData) changeTags(filename string, tags []string) error {
	fs.mutex.Lock()

	if _, ok := fs.info[filename]; !ok {
		fs.mutex.Unlock()
		return ErrFileIsNotExist
	}

	// Update map
	f := fs.info[filename]
	f.Tags = tags
	fs.info[filename] = f

	fs.mutex.Unlock()

	fs.write()

	return nil
}
