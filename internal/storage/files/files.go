package files

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/minio/sio"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
	"github.com/tags-drive/core/internal/storage/files/resizing"
)

const (
	typeImage = "image"
	typeFile  = "file"
	//
	timeBeforeDeleting = time.Hour * 24 * 7 // 7 days. File is deleted from storage and from disk after this time since user deleted the file
)

// Errors
var (
	ErrFileIsNotExist    = errors.New("the file doesn't exist")
	ErrAlreadyExist      = errors.New("file already exists")
	ErrFileDeletedAgain  = errors.New("file can't be deleted again")
	ErrOffsetOutOfBounds = errors.New("offset is out of bounds")
)

// FileInfo contains the information about a file
type FileInfo struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Type     string `json:"type"`              // typeImage or typeFile
	Origin   string `json:"origin"`            // Origin is a path to a file (params.DataFolder/filename)
	Preview  string `json:"preview,omitempty"` // Preview is a path to a resized image (only if Type == TypeImage)
	//
	Tags        []int     `json:"tags"`
	Description string    `json:"description"`
	Size        int64     `json:"size"`
	AddTime     time.Time `json:"addTime"`
	//
	Deleted      bool      `json:"deleted"`
	TimeToDelete time.Time `json:"timeToDelete"`
}

// storage for files metadata
type storage interface {
	init() error

	// getFile returns a file with passed filename
	getFile(id int) (FileInfo, error)

	// getFiles returns files
	//     expr - parsed logical expression
	//     search - string, which filename has to contain (lower case)
	getFiles(expr, search string) (files []FileInfo)

	// add adds a file
	addFile(filename, fileType string, tags []int, size int64, addTime time.Time) (id int)

	// renameFile renames a file
	renameFile(id int, newName string) error

	// updateFileTags updates tags of a file
	updateFileTags(id int, changedTagsID []int) error

	// updateFileDescription update description of a file
	updateFileDescription(id int, newDesc string) error

	// deleteFile marks file deleted and sets TimeToDelete
	// File can't be deleted several times (function should return ErrFileDeletedAgain)
	deleteFile(id int) error

	// deleteFileForce deletes file
	deleteFileForce(id int) error

	// recover removes file from Trash
	recover(id int)

	// deleteTagFromFiles deletes a tag (it's called when user deletes a tag)
	deleteTagFromFiles(tagID int)

	// getExpiredDeletedFiles returns names of files with expired TimeToDelete
	getExpiredDeletedFiles() []int

	shutdown() error
}

// FileStorage exposes methods for interactions with files
type FileStorage struct {
	storage storage
	logger  *log.Logger
}

// NewFileStorage creates new FileStorage
func NewFileStorage(lg *log.Logger) (*FileStorage, error) {
	var st storage

	switch params.StorageType {
	case params.JSONStorage:
		st = &jsonFileStorage{
			maxID:  0,
			files:  make(map[int]FileInfo),
			mutex:  new(sync.RWMutex),
			logger: lg,
		}
	default:
		st = &jsonFileStorage{
			maxID:  0,
			files:  make(map[int]FileInfo),
			mutex:  new(sync.RWMutex),
			logger: lg,
		}
	}

	fs := &FileStorage{
		storage: st,
		logger:  lg,
	}

	err := fs.storage.init()
	if err != nil {
		return nil, errors.Wrapf(err, "can't init files storage")
	}

	go fs.scheduleDeleting()

	return fs, nil
}

func (fs FileStorage) Get(expr string, s SortMode, search string, offset, count int) ([]FileInfo, error) {
	parsedExpr, err := aggregation.ParseLogicalExpr(expr)
	if err != nil {
		return []FileInfo{}, err
	}

	search = strings.ToLower(search)
	files := fs.storage.getFiles(parsedExpr, search)
	if len(files) == 0 && offset == 0 {
		// We don't return error, when there're no files and offset isn't set
		return []FileInfo{}, nil
	}

	if offset >= len(files) {
		return []FileInfo{}, ErrOffsetOutOfBounds
	}

	sortFiles(s, files)

	if count == 0 || offset+count > len(files) {
		count = len(files) - offset
	}

	return files[offset : offset+count], nil
}

func (fs FileStorage) GetFile(id int) (FileInfo, error) {
	return fs.storage.getFile(id)
}

func (fs FileStorage) GetRecent(number int) []FileInfo {
	files, _ := fs.Get("", SortByTimeDesc, "", 0, number)
	return files
}

func (fs FileStorage) Archive(ids []int) (body io.Reader, err error) {
	buff := bytes.NewBuffer([]byte(""))

	zipWriter := zip.NewWriter(buff)
	defer zipWriter.Close()

	for _, id := range ids {
		fileInfo, err := fs.storage.getFile(id)
		if err != nil {
			// Skip non-existent file
			continue
		}

		path := params.DataFolder + "/" + strconv.FormatInt(int64(id), 10)
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

		if params.Encrypt {
			_, err = sio.Decrypt(wr, f, sio.Config{Key: params.PassPhrase[:]})
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

func (fs FileStorage) Upload(f *multipart.FileHeader, tags []int) error {
	file, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "can't open a file")
	}
	defer file.Close()

	ext := filepath.Ext(f.Filename)
	var fileType string

	// Define file type
	switch ext {
	case ".jpg", ".JPG", ".jpeg", ".JPEG", ".png", ".PNG", ".gif", ".GIF", "webp", "WEBP":
		fileType = typeImage
	default:
		// Save a file
		fileType = typeFile
	}

	newFileID := fs.storage.addFile(f.Filename, fileType, tags, f.Size, time.Now())

	originPath := params.DataFolder + "/" + strconv.FormatInt(int64(newFileID), 10)

	// Save file
	switch fileType {
	case typeImage:
		previewPath := params.ResizedImagesFolder + "/" + strconv.FormatInt(int64(newFileID), 10)
		img, err := resizing.Decode(file)
		if err != nil {
			return err
		}

		// Save an original image
		r, err := resizing.Encode(img, ext)
		if err != nil {
			return err
		}
		err = copyToFile(r, originPath)
		if err != nil {
			return err
		}

		// Save a resized image
		// We can ignore errors and only log them because the main file was already saved
		img = resizing.Resize(img)
		r, err = resizing.Encode(img, ext)
		if err != nil {
			fs.logger.Errorf("can't encode a resized image %s: %s\n", f.Filename, err)
			break
		}
		err = copyToFile(r, previewPath)
		if err != nil {
			fs.logger.Errorf("can't save a resized image %s: %s\n", f.Filename, err)
		}
	default:
		// Save a file
		err := copyToFile(file, originPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// copyToFile copies data from src to new created file
func copyToFile(src io.Reader, path string) error {
	// We trunc file, if it already exists
	newFile, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}

	// Write file
	if params.Encrypt {
		_, err = sio.Encrypt(newFile, src, sio.Config{Key: params.PassPhrase[:]})
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
func (fs FileStorage) Rename(id int, newName string) error {
	err := fs.storage.renameFile(id, newName)
	if err != nil {
		return errors.Wrap(err, "can't rename file in a storage")
	}

	// We don't rename a file on disk, because id is constant
	return nil
}

func (fs FileStorage) ChangeTags(id int, tags []int) error {
	return fs.storage.updateFileTags(id, tags)
}

func (fs FileStorage) ChangeDescription(id int, newDescription string) error {
	return fs.storage.updateFileDescription(id, newDescription)
}

func (fs FileStorage) Delete(id int) error {
	return fs.storage.deleteFile(id)
}

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
	err = os.Remove(file.Origin)
	if err != nil {
		return err
	}

	if file.Preview != "" {
		// Delete the resized image
		err = os.Remove(file.Preview)
		if err != nil {
			// Only log error
			fs.logger.Errorf("can't delete a resized image %s: %s", file.Filename, err)
		}
	}

	return nil
}

func (fs FileStorage) DeleteTagFromFiles(tagID int) {
	fs.storage.deleteTagFromFiles(tagID)
}

// scheduleDeleting deletes files with expired TimeToDelete
// It has to be run in goroutine
func (fs FileStorage) scheduleDeleting() {
	ticker := time.NewTicker(time.Hour * 12)

	for ; true; <-ticker.C {
		fs.logger.Infoln("delete old files")

		var err error
		for _, id := range fs.storage.getExpiredDeletedFiles() {
			file, _ := fs.storage.getFile(id)
			err = fs.DeleteForce(id)
			if err != nil {
				fs.logger.Errorf("can't delete file \"%s\": %s\n", file.Filename, err)
			} else {
				fs.logger.Infof("file \"%s\" was successfully deleted\n", file.Filename)
			}
		}
	}
}

func (fs FileStorage) Recover(id int) {
	fs.storage.recover(id)
}

func (fs FileStorage) Shutdown() error {
	return fs.storage.shutdown()
}
