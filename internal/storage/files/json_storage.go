package files

import (
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
)

// jsonFileStorage implements files.storage interface.
// It is a map (id: cmd.FileInfo) with RWMutex
type jsonFileStorage struct {
	// maxID is max id of current files. It is computed in init() method
	maxID int
	files map[int]cmd.FileInfo
	mutex *sync.RWMutex

	logger *clog.Logger
	json   jsoniter.API
}

func newJsonFileStorage(lg *clog.Logger) *jsonFileStorage {
	return &jsonFileStorage{
		maxID:  0,
		files:  make(map[int]cmd.FileInfo),
		mutex:  new(sync.RWMutex),
		logger: lg,
		json:   jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (jfs *jsonFileStorage) init() error {
	// Create folders
	err := os.MkdirAll(params.DataFolder, 0666)
	if err != nil {
		return errors.Wrapf(err, "can't create a folder %s", params.DataFolder)
	}

	err = os.MkdirAll(params.ResizedImagesFolder, 0666)
	if err != nil {
		return errors.Wrapf(err, "can't create a folder %s", params.ResizedImagesFolder)
	}

	f, err := os.OpenFile(params.Files, os.O_RDWR, 0666)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			jfs.logger.Infof("file %s doesn't exist. Need to create a new file\n", params.Files)
			f, err = os.OpenFile(params.Files, os.O_CREATE|os.O_RDWR, 0666)
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

	f, err := os.OpenFile(params.Files, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		jfs.logger.Errorf("can't open file %s: %s\n", params.Files, err)
		return
	}

	enc := jfs.json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(jfs.files)

	f.Close()
}

// decode decodes js.info
func (jfs *jsonFileStorage) decode(r io.Reader) error {
	return jfs.json.NewDecoder(r).Decode(&jfs.files)
}

// checkFile return true if file with passed filename exists
func (jfs jsonFileStorage) checkFile(id int) bool {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	_, ok := jfs.files[id]
	return ok
}

func (jfs jsonFileStorage) getFile(id int) (cmd.FileInfo, error) {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	f, ok := jfs.files[id]
	if !ok {
		return cmd.FileInfo{}, ErrFileIsNotExist
	}
	return f, nil
}

// getFiles returns slice of cmd.FileInfo. If parsedExpr == "", it returns all files
func (jfs jsonFileStorage) getFiles(parsedExpr, search string) (files []cmd.FileInfo) {
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
	var goodFiles []cmd.FileInfo
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Filename), search) {
			goodFiles = append(goodFiles, f)
		}
	}

	return goodFiles
}

// addFile adds an element into js.files and call js.write()
// It also defines cmd.FileInfo.Origin and cmd.FileInfo.Preview (if file is image) as
// `params.DataFolder + "/" + id` and `params.ResizedImagesFolder + "/" + id`
func (jfs *jsonFileStorage) addFile(filename string, fileType cmd.Ext, tags []int, size int64, addTime time.Time) (id int) {
	fileInfo := cmd.FileInfo{Filename: filename,
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
	if fileType.FileType == cmd.FileTypeImage {
		fileInfo.Preview = params.ResizedImagesFolder + "/" + strconv.FormatInt(int64(fileID), 10)
	}

	jfs.files[jfs.maxID] = fileInfo

	jfs.mutex.Unlock()

	jfs.write()

	return fileID
}

// renameFile renames a file
func (jfs *jsonFileStorage) renameFile(id int, newName string) (cmd.FileInfo, error) {
	if !jfs.checkFile(id) {
		return cmd.FileInfo{}, ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.files[id]
	f.Filename = newName
	jfs.files[id] = f

	jfs.mutex.Unlock()

	jfs.write()

	return f, nil
}

func (jfs *jsonFileStorage) updateFileTags(id int, changedTagsID []int) (cmd.FileInfo, error) {
	if !jfs.checkFile(id) {
		return cmd.FileInfo{}, ErrFileIsNotExist
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

	return f, nil
}

func (jfs *jsonFileStorage) updateFileDescription(id int, newDesc string) (cmd.FileInfo, error) {
	if !jfs.checkFile(id) {
		return cmd.FileInfo{}, ErrFileIsNotExist
	}

	jfs.mutex.Lock()

	// Update map
	f := jfs.files[id]
	f.Description = newDesc
	jfs.files[id] = f

	jfs.mutex.Unlock()

	jfs.write()

	return f, nil
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

func (jfs *jsonFileStorage) addTagsToFiles(filesIDs, tagsID []int) {
	merge := func(a, b []int) []int {
		t := make(map[int]struct{}, len(a)+len(b))
		for i := range a {
			t[a[i]] = struct{}{}
		}
		for i := range b {
			t[b[i]] = struct{}{}
		}

		res := make([]int, 0, len(a)+len(b))
		for k := range t {
			res = append(res, k)
		}

		return res
	}

	goodID := func(id int) bool {
		for i := range filesIDs {
			if filesIDs[i] == id {
				return true
			}
		}

		return false
	}

	jfs.mutex.Lock()

	for id, f := range jfs.files {
		if !goodID(id) {
			continue
		}

		f.Tags = merge(f.Tags, tagsID)

		jfs.files[id] = f
	}

	jfs.mutex.Unlock()

	jfs.write()
}

func (jfs *jsonFileStorage) removeTagsFromFiles(filesIDs, tagsID []int) {
	exclude := func(a, b []int) []int {
		t := make(map[int]bool, len(a)+len(b))
		for i := range a {
			t[a[i]] = true
		}
		for i := range b {
			t[b[i]] = false
		}

		res := make([]int, 0, len(a)+len(b))
		for k, v := range t {
			if v {
				res = append(res, k)
			}
		}

		return res
	}

	goodID := func(id int) bool {
		for i := range filesIDs {
			if filesIDs[i] == id {
				return true
			}
		}

		return false
	}

	jfs.mutex.Lock()

	for id, f := range jfs.files {
		if !goodID(id) {
			continue
		}

		f.Tags = exclude(f.Tags, tagsID)

		jfs.files[id] = f
	}

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

func (jfs jsonFileStorage) shutdown() error {
	// We have not to do any special operations because we update json file on every change.
	// Also there are no any requests because server is already down. But it's better to check the mutex
	// just in case.

	jfs.mutex.Lock()
	jfs.mutex.Unlock()

	// There will be no any new requests.

	return nil
}
