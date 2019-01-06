package files

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"
	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
)

// jsonFileStorage implements files.storage interface.
// It is a map (id: FileInfo) with RWMutex
type jsonFileStorage struct {
	// maxID is max id of current files. It is computed in init() method
	maxID int
	files map[int]FileInfo
	mutex *sync.RWMutex

	logger *log.Logger
}

func (jfs *jsonFileStorage) init() error {
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
			jfs.logger.Infof("File %s doesn't exist. Need to create a new file\n", params.Files)
			f, err = os.OpenFile(params.Files, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return errors.Wrap(err, "can't create a new file")
			}
			// Write empty structure
			f.Write([]byte("{}"))
			// Can exit because we don't need to decode files from the file
			f.Close()
			// We don't have to compute maxID, because there're no any files
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.Files)
	}

	defer f.Close()

	err = jfs.decode(f)
	if err != nil {
		return errors.Wrap(err, "can't decode file")
	}

	// Compute maxID
	for id := range jfs.files {
		if id > jfs.maxID {
			jfs.maxID = id
		}
	}

	return nil
}

// write writes js.info into params.TagsFile
func (jfs jsonFileStorage) write() {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	f, err := os.OpenFile(params.Files, os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		jfs.logger.Errorf("Can't open file %s: %s\n", params.Files, err)
		return
	}

	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(jfs.files)

	f.Close()
}

// decode decodes js.info
func (jfs *jsonFileStorage) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&jfs.files)
}

// checkFile return true if file with passed filename exists
func (jfs jsonFileStorage) checkFile(id int) bool {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	_, ok := jfs.files[id]
	return ok
}

func (jfs jsonFileStorage) getFile(id int) (FileInfo, error) {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	f, ok := jfs.files[id]
	if !ok {
		return FileInfo{}, ErrFileIsNotExist
	}
	return f, nil
}

// getFiles returns slice of FileInfo. If parsedExpr == "", it returns all files
func (jfs jsonFileStorage) getFiles(parsedExpr, search string) (files []FileInfo) {
	jfs.mutex.RLock()

	for _, v := range jfs.files {
		if aggregation.IsGoodFile(parsedExpr, v.Tags) {
			files = append(files, v)
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

// addFile adds an element into js.files and call js.write()
// It also defines FileInfo.Origin and FileInfo.Preview (if file is image) as
// `params.DataFolder + "/" + id` and `params.ResizedImagesFolder + "/" + id`
func (jfs *jsonFileStorage) addFile(filename, fileType string, tags []int, size int64, addTime time.Time) (id int) {
	fileInfo := FileInfo{Filename: filename,
		Type:    fileType,
		Tags:    tags,
		Size:    size,
		AddTime: addTime,
	}

	// We need a special var for thread safety
	fileID := 0

	if fileInfo.Tags == nil {
		fileInfo.Tags = []int{} // https://github.com/tags-drive/core/issues/19
	}

	jfs.mutex.Lock()

	// Set id
	jfs.maxID++
	fileID = jfs.maxID
	fileInfo.ID = fileID

	fileInfo.Origin = params.DataFolder + "/" + strconv.FormatInt(int64(fileID), 10)
	if fileType == typeImage {
		fileInfo.Preview = params.ResizedImagesFolder + "/" + strconv.FormatInt(int64(fileID), 10)
	}

	jfs.files[jfs.maxID] = fileInfo

	jfs.mutex.Unlock()

	jfs.write()

	return fileID
}

// renameFile renames a file
func (jfs *jsonFileStorage) renameFile(id int, newName string) error {
	if !jfs.checkFile(id) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.files[id]
	delete(jfs.files, id)
	f.Filename = newName
	// We don't have to change Origin, because we don't change filename
	jfs.files[id] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

func (jfs *jsonFileStorage) updateFileTags(id int, changedTagsID []int) error {
	if !jfs.checkFile(id) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.files[id]
	if changedTagsID == nil {
		changedTagsID = []int{} // https://github.com/tags-drive/core/issues/19
	}
	f.Tags = changedTagsID
	jfs.files[id] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

func (jfs *jsonFileStorage) updateFileDescription(id int, newDesc string) error {
	if !jfs.checkFile(id) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.files[id]
	f.Description = newDesc
	jfs.files[id] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

// deleteFile sets Deleted = true and update TimeToDelete
func (jfs *jsonFileStorage) deleteFile(id int) error {
	if !jfs.checkFile(id) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	f := jfs.files[id]
	if f.Deleted {
		jfs.mutex.Unlock()
		return ErrFileDeletedAgain
	}

	f.Deleted = true
	f.TimeToDelete = time.Now().Add(timeBeforeDeleting)
	jfs.files[id] = f

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

// deleteFile deletes an element (from structure) and call js.write()
func (jfs *jsonFileStorage) deleteFileForce(id int) error {
	if !jfs.checkFile(id) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	delete(jfs.files, id)

	jfs.mutex.Unlock()

	jfs.write()

	return nil
}

// recover sets Deleted = false
func (jfs *jsonFileStorage) recover(id int) {
	if !jfs.checkFile(id) {
		return
	}

	jfs.mutex.Lock()

	if !jfs.files[id].Deleted {
		return
	}

	f := jfs.files[id]
	f.Deleted = false
	f.TimeToDelete = time.Time{}
	jfs.files[id] = f

	jfs.mutex.Unlock()

	jfs.write()
}

func (jfs *jsonFileStorage) deleteTagFromFiles(tagID int) {
	jfs.mutex.Lock()

	for id, f := range jfs.files {
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

		jfs.files[id] = f
	}

	jfs.mutex.Unlock()

	jfs.write()
}

// getExpiredDeletedFiles returns ids of files with expired TimeToDelete
func (jfs *jsonFileStorage) getExpiredDeletedFiles() []int {
	jfs.mutex.RLock()

	var filesForDeleting []int
	now := time.Now()
	for id, file := range jfs.files {
		if file.Deleted && file.TimeToDelete.Before(now) {
			filesForDeleting = append(filesForDeleting, id)
		}
	}

	jfs.mutex.RUnlock()

	return filesForDeleting
}
