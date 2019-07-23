package files

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/minio/sio"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/storage/files/aggregation"
	"github.com/tags-drive/core/internal/storage/files/extensions"
	"github.com/tags-drive/core/internal/storage/files/resizing"
	"github.com/tags-drive/core/internal/utils"
)

const (
	originURLPrefix  = "data/"
	previewURLPrefix = "data/resized/"
)

// Errors
var (
	ErrFileIsNotExist    = errors.New("the file doesn't exist")
	ErrAlreadyExist      = errors.New("file already exists")
	ErrFileDeletedAgain  = errors.New("file can't be deleted again")
	ErrOffsetOutOfBounds = errors.New("offset is out of bounds")
	ErrEmptyNewName      = errors.New("new name can't be empty")
)

// FileStorage exposes methods for interactions with files
type FileStorage struct {
	config Config

	storage metadataStorage
	logger  *clog.Logger
}

// NewFileStorage creates new FileStorage
func NewFileStorage(cnf Config, lg *clog.Logger) (*FileStorage, error) {
	// Create folders
	err := os.MkdirAll(cnf.DataFolder, 0666)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create a folder %s", cnf.DataFolder)
	}

	err = os.MkdirAll(cnf.ResizedImagesFolder, 0666)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create a folder %s", cnf.ResizedImagesFolder)
	}

	var st metadataStorage

	switch cnf.StorageType {
	case "json":
		fallthrough
	default:
		st = newJsonFileStorage(cnf, lg)
	}

	fs := &FileStorage{
		config:  cnf,
		storage: st,
		logger:  lg,
	}

	err = fs.storage.init()
	if err != nil {
		return nil, errors.Wrapf(err, "can't init files storage")
	}

	return fs, nil
}

// StartBackgroundServices starts all background services
func (fs FileStorage) StartBackgroundServices() {
	go fs.scheduleDeleting()
}

// Get returns all "good" sorted files
//
// If cnf.Expr isn't valid, Get returns ErrBadExpessionSyntax
func (fs FileStorage) Get(cnf GetFilesConfig) ([]File, error) {
	var (
		offset = cnf.Offset
		count  = cnf.Count
	)

	parsedExpr, err := aggregation.ParseLogicalExpr(cnf.Expr)
	if err != nil {
		return []File{}, err
	}

	search := strings.ToLower(cnf.Search)
	files := fs.storage.getFiles(parsedExpr, search, cnf.IsRegexp)
	if len(files) == 0 && offset == 0 {
		// We don't return error, when there're no files and offset isn't set
		return []File{}, nil
	}

	// Filter files at first
	if cnf.Filter != nil {
		files, err = cnf.Filter(files)
		if err != nil {
			return []File{}, err
		}
	}

	if offset >= len(files) {
		return []File{}, ErrOffsetOutOfBounds
	}

	sortFiles(cnf.SortMode, files)

	if count == 0 || offset+count > len(files) {
		count = len(files) - offset
	}

	return files[offset : offset+count], nil
}

// GetFile returns a file with passed id
func (fs FileStorage) GetFile(id int) (File, error) {
	return fs.storage.getFile(id)
}

// CheckFile checks if file with passed id exists
func (fs FileStorage) CheckFile(id int) bool {
	return fs.storage.checkFile(id)
}

// GetFiles returns files with passed ids
func (fs FileStorage) GetFiles(ids ...int) []File {
	return fs.storage.getFilesWithIDs(ids...)
}

// GetRecent returns the last uploaded files
func (fs FileStorage) GetRecent(number int) []File {
	cnf := GetFilesConfig{SortMode: SortByTimeDesc, Count: number}
	files, _ := fs.Get(cnf)
	return files
}

// Archive archives passed files and returns io.Reader with archive
func (fs FileStorage) Archive(ids []int) (body io.Reader, err error) {
	// Max size of an archive in memory is 20MB
	buff := utils.NewBuffer(20 << 20)
	defer buff.Finish()

	zipWriter := zip.NewWriter(buff)
	defer zipWriter.Close()

	for _, id := range ids {
		fileInfo, err := fs.storage.getFile(id)
		if err != nil {
			// Skip non-existent file
			continue
		}

		path := fs.config.DataFolder + "/" + strconv.FormatInt(int64(id), 10)
		f, err := os.Open(path)
		if err != nil {
			fs.logger.Errorf("can't load file \"%s\"\n", fileInfo.Filename)
			continue
		}
		stat, err := f.Stat()
		if err != nil {
			fs.logger.Errorf("can't load file \"%s\"\n", fileInfo.Filename)
			continue
		}

		header, _ := zip.FileInfoHeader(stat)
		header.Name = fileInfo.Filename // Set right filename
		header.Method = zip.Deflate

		wr, err := zipWriter.CreateHeader(header)
		if err != nil {
			fs.logger.Errorf("can't load file \"%s\"\n", fileInfo.Filename)
			f.Close()
			continue
		}

		if fs.config.Encrypt {
			_, err = sio.Decrypt(wr, f, sio.Config{Key: fs.config.PassPhrase[:]})
		} else {
			_, err = io.Copy(wr, f)
		}

		if err != nil {
			fs.logger.Errorf("can't load file \"%s\"\n", fileInfo.Filename)
		}

		f.Close()
	}

	return buff, nil
}

// Upload uploads a new file
func (fs FileStorage) Upload(f *multipart.FileHeader, tags []int) (err error) {
	file, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "can't open a file")
	}
	defer file.Close()

	ext := filepath.Ext(f.Filename)
	fileType := extensions.GetExt(ext)

	newFileID := fs.storage.addFile(f.Filename, fileType, tags, f.Size, time.Now())

	// If we will get a major error, we will have to panic to delete record in file storage
	defer func() {
		if r := recover(); r != nil {
			// Remove record in storage
			// We can only log this error
			e := fs.storage.deleteFileForce(newFileID)
			if e != nil {
				fs.logger.Errorf("can't delete record in file storage after error in Upload function: %s\n", e)
			}

			e, ok := r.(error)
			if !ok {
				err = errors.New("unexpected error")
				return
			}

			err = e
		}
	}()

	originPath := fs.config.DataFolder + "/" + strconv.FormatInt(int64(newFileID), 10)

	// Save file
	switch fileType.FileType {
	case extensions.FileTypeImage:
		// Create 2 io.Reader from file
		imageReader := new(bytes.Buffer)
		fileReader := io.TeeReader(file, imageReader)

		// Save an original image
		err = fs.copyToFile(fileReader, originPath)
		if err != nil {
			panic(err)
		}

		// After saving the original file we can ignore errors and only log them.

		// Convert imageReader into image.Image
		previewPath := fs.config.ResizedImagesFolder + "/" + strconv.Itoa(newFileID)
		img, err := resizing.Decode(imageReader)
		if err != nil {
			fs.logger.Errorf("can't decode an image %s: %s\n", f.Filename, err)
			break
		}

		// Save a resized image
		img = resizing.Resize(img)
		r, err := resizing.Encode(img, ext)
		if err != nil {
			fs.logger.Errorf("can't encode a resized image %s: %s\n", f.Filename, err)
			break
		}
		err = fs.copyToFile(r, previewPath)
		if err != nil {
			fs.logger.Errorf("can't save a resized image %s: %s\n", f.Filename, err)
		}
	default:
		// Save a file
		err := fs.copyToFile(file, originPath)
		if err != nil {
			panic(err)
		}
	}

	// TODO: does it really help?
	// resizing.Decode() allocates a lot of memory. GC doesn't keep up to free it
	// when there are a lot of Upload() calls. Calling runtime.GC() can
	// decrease max memory usage by 1.5 times with very small performance drop.
	runtime.GC()

	return nil
}

// copyToFile copies data from src to new created file
func (fs FileStorage) copyToFile(src io.Reader, path string) error {
	// We trunc file, if it already exists
	newFile, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}

	// Write file
	if fs.config.Encrypt {
		_, err = sio.Encrypt(newFile, src, sio.Config{Key: fs.config.PassPhrase[:]})
	} else {
		_, err = io.Copy(newFile, src)
	}

	newFile.Close()

	if err != nil {
		// Deleting of the bad file
		os.Remove(path)
		return errors.Wrap(err, "can't copy a new file")
	}

	return nil
}

// Rename renames a file
func (fs FileStorage) Rename(id int, newName string) (File, error) {
	file, err := fs.storage.renameFile(id, newName)
	if err != nil {
		return File{}, errors.Wrap(err, "can't rename file in a storage")
	}

	// We don't rename a file on a disk, because it is saved by id
	return file, nil
}

// ChangeTags changes the tags
func (fs FileStorage) ChangeTags(id int, tags []int) (File, error) {
	return fs.storage.updateFileTags(id, tags)
}

// ChangeDescription changes the description
func (fs FileStorage) ChangeDescription(id int, newDescription string) (File, error) {
	return fs.storage.updateFileDescription(id, newDescription)
}

// Delete "moves" a file into Trash
func (fs FileStorage) Delete(id int) error {
	return fs.storage.deleteFile(id)
}

// Recover "removes" file from Trash
func (fs FileStorage) Recover(id int) {
	fs.storage.recover(id)
}

// DeleteForce deletes file from storage and from disk
func (fs FileStorage) DeleteForce(id int) error {
	file, err := fs.storage.getFile(id)
	if err != nil {
		return err
	}

	err = fs.storage.deleteFileForce(id)
	if err != nil {
		return err
	}

	// Delete the original file
	filepath := fs.config.DataFolder + "/" + strconv.Itoa(file.ID)
	err = os.Remove(filepath)
	if err != nil {
		return err
	}

	if file.Preview != "" {
		// Delete the resized image

		filepath = fs.config.ResizedImagesFolder + "/" + strconv.Itoa(file.ID)
		err = os.Remove(filepath)
		if err != nil {
			// Only log error
			fs.logger.Errorf("can't delete a resized image %s: %s\n", file.Filename, err)
		}
	}

	return nil
}

// AddTagsToFiles adds a tag to files
func (fs FileStorage) AddTagsToFiles(filesIDs, tagsIDs []int) {
	fs.storage.addTagsToFiles(filesIDs, tagsIDs)
}

// RemoveTagsFromFiles removes tags from files
func (fs FileStorage) RemoveTagsFromFiles(filesIDs, tagsIDs []int) {
	fs.storage.removeTagsFromFiles(filesIDs, tagsIDs)
}

// RemoveTagFromAllFiles deletes a tag from all files
func (fs FileStorage) RemoveTagFromAllFiles(tagID int) {
	fs.storage.removeTagFromAllFiles(tagID)
}

// scheduleDeleting deletes files with expired TimeToDelete
// It has to be run in goroutine
func (fs FileStorage) scheduleDeleting() {
	ticker := time.NewTicker(time.Hour * 12)

	for ; true; <-ticker.C {
		fs.logger.Debugln("delete old files")

		var err error
		for _, id := range fs.storage.getExpiredDeletedFiles() {
			file, _ := fs.storage.getFile(id)
			err = fs.DeleteForce(id)
			if err != nil {
				fs.logger.Errorf("can't remove file \"%s\" from trash: %s\n", file.Filename, err)
			} else {
				fs.logger.Debugf("file \"%s\" was successfully deleted\n", file.Filename)
			}
		}
	}
}

// Shutdown gracefully shutdown FileStorage
func (fs FileStorage) Shutdown() error {
	return fs.storage.shutdown()
}
