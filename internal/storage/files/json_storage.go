package files

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/sio"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
	"github.com/tags-drive/core/internal/storage/files/extensions"
)

// saveInterval is used in saveOnDisk. It defines interval between calls of jfs.write()
const saveInterval = time.Second * 10

// jsonFileStorage implements files.storage interface.
// It is a map (id: FileInfo) with RWMutex
type jsonFileStorage struct {
	// maxID is max id of current files. It is computed in init() method
	maxID int
	files map[int]File
	mutex *sync.RWMutex

	logger *clog.Logger
	json   jsoniter.API

	shutdownChan chan struct{}
	// number of changes since last write() call
	changes *uint32
}

func newJsonFileStorage(lg *clog.Logger) *jsonFileStorage {
	changes := new(uint32)
	atomic.StoreUint32(changes, 0)

	return &jsonFileStorage{
		maxID:        0,
		files:        make(map[int]File),
		mutex:        new(sync.RWMutex),
		logger:       lg,
		json:         jsoniter.ConfigCompatibleWithStandardLibrary,
		shutdownChan: make(chan struct{}),
		changes:      changes,
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

	f, err := os.OpenFile(params.FilesJSONFile, os.O_RDWR, 0666)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			// We don't have to compute maxID, because there're no any files
			// Can exit because we don't need to decode files from the file
			err := jfs.createNewFile()
			if err != nil {
				return err
			}

			go jfs.saveOnDisk()
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.FilesJSONFile)
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

	go jfs.saveOnDisk()
	return nil
}

func (jfs jsonFileStorage) createNewFile() error {
	jfs.logger.Infof("file %s doesn't exist. Need to create a new file\n", params.FilesJSONFile)

	// Just create a new file
	f, err := os.OpenFile(params.FilesJSONFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	f.Close()

	// Write empty files map
	jfs.write()

	return nil
}

// saveOnDisk calls jfs.write() every saveInterval seconds. It must be ran in goroutine
// It finishes when jfs.shutdownChan is closed.
func (jfs *jsonFileStorage) saveOnDisk() {
	ticker := time.NewTicker(saveInterval)
	for {
		select {
		case <-ticker.C:
			changed := atomic.LoadUint32(jfs.changes) != 0
			// Reset jfs.changes
			atomic.StoreUint32(jfs.changes, 0)

			if changed {
				jfs.write()
			}

		case <-jfs.shutdownChan:
			ticker.Stop()
			return
		}
	}
}

// write writes js.info into params.Files
func (jfs jsonFileStorage) write() {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	f, err := os.OpenFile(params.FilesJSONFile, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		jfs.logger.Errorf("can't open file %s: %s\n", params.FilesJSONFile, err)
		return
	}
	defer f.Close()

	if !params.Encrypt {
		// Encode directly into the file
		enc := jfs.json.NewEncoder(f)
		if params.Debug {
			enc.SetIndent("", "  ")
		}
		err := enc.Encode(jfs.files)
		if err != nil {
			jfs.logger.Warnf("can't write '%s': %s\n", params.FilesJSONFile, err)
		}

		return
	}

	// Encode into buffer
	buff := bytes.NewBuffer([]byte{})
	enc := jfs.json.NewEncoder(buff)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(jfs.files)

	// Write into the file (params.Encrypt is true, if we are here)
	_, err = sio.Encrypt(f, buff, sio.Config{Key: params.PassPhrase[:]})

	if err != nil {
		jfs.logger.Warnf("can't write '%s': %s\n", params.FilesJSONFile, err)
	}
}

// decode decodes js.info
func (jfs *jsonFileStorage) decode(r io.Reader) error {
	if !params.Encrypt {
		return jfs.json.NewDecoder(r).Decode(&jfs.files)
	}

	// Have to decrypt at first
	buff := bytes.NewBuffer([]byte{})
	_, err := sio.Decrypt(buff, r, sio.Config{Key: params.PassPhrase[:]})
	if err != nil {
		return err
	}

	return jfs.json.NewDecoder(buff).Decode(&jfs.files)
}

// checkFile return true if file with passed filename exists
func (jfs jsonFileStorage) checkFile(id int) bool {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	_, ok := jfs.files[id]
	return ok
}

func (jfs jsonFileStorage) getFile(id int) (File, error) {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	f, ok := jfs.files[id]
	if !ok {
		return File{}, ErrFileIsNotExist
	}
	return f, nil
}

// getFiles returns slice of FileInfo. If parsedExpr == "", it returns all files
func (jfs jsonFileStorage) getFiles(parsedExpr aggregation.LogicalExpr, search string, isRegexp bool) (files []File) {
	jfs.mutex.RLock()

	files = make([]File, 0, len(jfs.files))
	for _, v := range jfs.files {
		if aggregation.IsGoodFile(parsedExpr, v.Tags) {
			files = append(files, v)
		}
	}

	jfs.mutex.RUnlock()

	if search == "" {
		return files
	}

	var reg *regexp.Regexp
	if isRegexp {
		// search must be valid regular expression
		reg = regexp.MustCompile(search)
	}

	// Need to remove files with incorrect name
	var goodFiles []File
	for i := range files {
		if isRegexp && reg.MatchString(files[i].Filename) {
			goodFiles = append(goodFiles, files[i])
		} else if strings.Contains(strings.ToLower(files[i].Filename), search) {
			goodFiles = append(goodFiles, files[i])
		}
	}

	return goodFiles
}

// addFile adds an element into js.files and call js.write()
// It also defines FileInfo.Origin and FileInfo.Preview (if file is image) as
// `params.DataFolder + "/" + id` and `params.ResizedImagesFolder + "/" + id`
func (jfs *jsonFileStorage) addFile(filename string, fileType extensions.Ext, tags []int, size int64, addTime time.Time) (id int) {
	fileInfo := File{Filename: filename,
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
	defer jfs.mutex.Unlock()

	// Set id
	jfs.maxID++
	fileID = jfs.maxID
	fileInfo.ID = fileID

	fileInfo.Origin = params.DataFolder + "/" + strconv.FormatInt(int64(fileID), 10)
	if fileType.FileType == extensions.FileTypeImage {
		fileInfo.Preview = params.ResizedImagesFolder + "/" + strconv.FormatInt(int64(fileID), 10)
	}

	jfs.files[jfs.maxID] = fileInfo

	atomic.AddUint32(jfs.changes, 1)

	return fileID
}

// renameFile renames a file
func (jfs *jsonFileStorage) renameFile(id int, newName string) (File, error) {
	if !jfs.checkFile(id) {
		return File{}, ErrFileIsNotExist
	}

	jfs.mutex.Lock()
	defer jfs.mutex.Unlock()

	f := jfs.files[id]
	f.Filename = newName
	jfs.files[id] = f

	atomic.AddUint32(jfs.changes, 1)

	return f, nil
}

func (jfs *jsonFileStorage) updateFileTags(id int, changedTagsID []int) (File, error) {
	if !jfs.checkFile(id) {
		return File{}, ErrFileIsNotExist
	}

	if changedTagsID == nil {
		changedTagsID = []int{} // https://github.com/tags-drive/core/issues/19
	}

	jfs.mutex.Lock()
	defer jfs.mutex.Unlock()

	f := jfs.files[id]
	f.Tags = changedTagsID
	jfs.files[id] = f

	atomic.AddUint32(jfs.changes, 1)

	return f, nil
}

func (jfs *jsonFileStorage) updateFileDescription(id int, newDesc string) (File, error) {
	if !jfs.checkFile(id) {
		return File{}, ErrFileIsNotExist
	}

	jfs.mutex.Lock()
	defer jfs.mutex.Unlock()

	f := jfs.files[id]
	f.Description = newDesc
	jfs.files[id] = f

	atomic.AddUint32(jfs.changes, 1)

	return f, nil
}

// deleteFile sets Deleted = true and update TimeToDelete
func (jfs *jsonFileStorage) deleteFile(id int) error {
	if !jfs.checkFile(id) {
		return ErrFileIsNotExist
	}

	deleteTime := time.Now().Add(timeBeforeDeleting)

	jfs.mutex.Lock()
	defer jfs.mutex.Unlock()

	f := jfs.files[id]
	if f.Deleted {
		return ErrFileDeletedAgain
	}

	f.Deleted = true
	f.TimeToDelete = deleteTime
	jfs.files[id] = f

	atomic.AddUint32(jfs.changes, 1)

	return nil
}

// deleteFile deletes an element (from structure) and call js.write()
func (jfs *jsonFileStorage) deleteFileForce(id int) error {
	if !jfs.checkFile(id) {
		return ErrFileIsNotExist
	}

	jfs.mutex.Lock()
	defer jfs.mutex.Unlock()

	delete(jfs.files, id)

	atomic.AddUint32(jfs.changes, 1)

	return nil
}

// recover sets Deleted = false
func (jfs *jsonFileStorage) recover(id int) {
	if !jfs.checkFile(id) {
		return
	}

	jfs.mutex.Lock()
	defer jfs.mutex.Unlock()

	if !jfs.files[id].Deleted {
		return
	}

	f := jfs.files[id]
	f.Deleted = false
	f.TimeToDelete = time.Time{}
	jfs.files[id] = f

	atomic.AddUint32(jfs.changes, 1)
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
	defer jfs.mutex.Unlock()

	for id, f := range jfs.files {
		if !goodID(id) {
			continue
		}

		f.Tags = merge(f.Tags, tagsID)

		jfs.files[id] = f
	}

	atomic.AddUint32(jfs.changes, 1)
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
	defer jfs.mutex.Unlock()

	for id, f := range jfs.files {
		if !goodID(id) {
			continue
		}

		f.Tags = exclude(f.Tags, tagsID)

		jfs.files[id] = f
	}

	atomic.AddUint32(jfs.changes, 1)
}

func (jfs *jsonFileStorage) deleteTagFromFiles(tagID int) {
	jfs.mutex.Lock()
	defer jfs.mutex.Unlock()

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

	atomic.AddUint32(jfs.changes, 1)
}

// getExpiredDeletedFiles returns ids of files with expired TimeToDelete
func (jfs *jsonFileStorage) getExpiredDeletedFiles() []int {
	jfs.mutex.RLock()
	defer jfs.mutex.RUnlock()

	var filesForDeleting []int
	now := time.Now()
	for id, file := range jfs.files {
		if file.Deleted && file.TimeToDelete.Before(now) {
			filesForDeleting = append(filesForDeleting, id)
		}
	}

	return filesForDeleting
}

func (jfs jsonFileStorage) shutdown() error {
	// Stop saveOnDisk goroutine
	close(jfs.shutdownChan)

	// Wait for all locks
	jfs.mutex.Lock()
	jfs.mutex.Unlock()

	// Write changes
	jfs.write()

	// There will be no any new requests.

	return nil
}
