package files

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

func (fs filesData) get(filename string) (FileInfo, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	f, ok := fs.info[filename]
	if !ok {
		return FileInfo{}, ErrFileIsNotExist
	}
	return f, nil
}

// add adds an element into fs.info and call fs.write()
func (fs *filesData) add(info FileInfo) error {
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

// rename renames a file
func (fs *filesData) rename(oldName string, newName string) error {
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

func (fs *filesData) deleteTag(tagID int) {
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

func (fs *filesData) changeTags(filename string, changedTagsID []int) error {
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

func (fs *filesData) changeDescription(filename string, newDesc string) error {
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
