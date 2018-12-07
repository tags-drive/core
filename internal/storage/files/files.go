package files

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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
	ErrFileIsNotExist   = errors.New("the file doesn't exist")
	ErrAlreadyExist     = errors.New("file already exists")
	ErrFileDeletedAgain = errors.New("file can't be deleted again")
)

// FileInfo contains the information about a file
type FileInfo struct {
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

type storage interface {
	init() error

	// getFile returns a file with passed filename
	getFile(filename string) (FileInfo, error)

	// getFiles returns files
	//     expr - parsed logical expression
	//     search - string, which filename has to contain (lower case)
	getFiles(expr, search string) (files []FileInfo)

	// add adds a file
	addFile(info FileInfo) error

	// renameFile renames a file
	renameFile(oldName string, newName string) error

	// updateFileTags updates tags of a file
	updateFileTags(filename string, changedTagsID []int) error

	// updateFileDescription update description of a file
	updateFileDescription(filename string, newDesc string) error

	// deleteFile marks file deleted and sets TimeToDelete
	// File can't be deleted several times (function should return ErrFileDeletedAgain)
	deleteFile(filename string) error

	// deleteFileForce deletes file
	deleteFileForce(filename string) error

	// recover removes file from Trash
	recover(filename string)

	// deleteTagFromFiles deletes a tag (it's called when user deletes a tag)
	deleteTagFromFiles(tagID int)

	// getExpiredDeletedFiles returns names of files with expired TimeToDelete
	getExpiredDeletedFiles() []string
}

// FileStorage exposes methods for interactions with files
type FileStorage struct {
	storage storage
}

// Init inits fs.storage
func (fs *FileStorage) Init() error {
	switch params.StorageType {
	case params.JSONStorage:
		fs.storage = &jsonFileStorage{
			info:  make(map[string]FileInfo),
			mutex: new(sync.RWMutex),
		}
	default:
		// Default storage is jsonFileStorage
		fs.storage = &jsonFileStorage{
			info:  make(map[string]FileInfo),
			mutex: new(sync.RWMutex),
		}
	}

	err := fs.storage.init()
	if err != nil {
		return errors.Wrapf(err, "can't init storage")
	}

	go fs.scheduleDeleting()

	return nil
}

func (fs FileStorage) Get(expr string, s SortMode, search string) ([]FileInfo, error) {
	parsedExpr, err := aggregation.ParseLogicalExpr(expr)
	if err != nil {
		return []FileInfo{}, err
	}

	search = strings.ToLower(search)
	files := fs.storage.getFiles(parsedExpr, search)
	sortFiles(s, files)
	return files, nil
}

func (fs FileStorage) GetRecent(number int) []FileInfo {
	files, _ := fs.Get("", SortByTimeDesc, "")
	if len(files) > number {
		files = files[:number]
	}
	return files
}

func (fs FileStorage) Archive(files []string) (body io.Reader, err error) {
	buff := bytes.NewBuffer([]byte(""))

	zipWriter := zip.NewWriter(buff)
	defer zipWriter.Close()

	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			log.Errorf("Can't load file \"%s\"\n", filename)
			continue
		}
		stat, err := f.Stat()
		if err != nil {
			log.Errorf("Can't load file \"%s\"\n", filename)
			continue
		}

		header, _ := zip.FileInfoHeader(stat)
		header.Method = zip.Deflate

		wr, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Errorf("Can't load file \"%s\"\n", filename)
			f.Close()
			continue
		}

		if params.Encrypt {
			_, err = sio.Decrypt(wr, f, sio.Config{Key: params.Key[:]})
		} else {
			_, err = io.Copy(wr, f)
		}

		if err != nil {
			log.Errorf("Can't load file \"%s\"\n", filename)
		}

		f.Close()
	}

	return buff, nil
}

func (fs FileStorage) Upload(f *multipart.FileHeader, tags []int) error {
	// At first, check does file exist
	if f, err := os.Open(params.DataFolder + "/" + f.Filename); !os.IsNotExist(err) {
		f.Close()
		return ErrAlreadyExist
	}

	// Uploading
	file, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "can't open a file")
	}
	defer file.Close()

	ext := filepath.Ext(f.Filename)
	info := FileInfo{
		Filename: f.Filename,
		Size:     f.Size,
		AddTime:  time.Now(),
		Origin:   params.DataFolder + "/" + f.Filename,
		Tags:     tags}

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		// Need to save original and resized image
		info.Type = typeImage
		info.Preview = params.ResizedImagesFolder + "/" + f.Filename
		img, err := resizing.Decode(file)
		if err != nil {
			return err
		}

		// Save an original image
		r, err := resizing.Encode(img, ext)
		if err != nil {
			return err
		}
		err = copyToFile(r, info.Origin)
		if err != nil {
			return err
		}
		
		// Save a resized image
		// We can ignore errors and only log them because the main file was already saved
		img = resizing.Resize(img)
		r, err = resizing.Encode(img, ext)
		if err != nil {
			log.Errorf("Can't encode a resized image %s: %s\n", info.Filename, err)
			break
		}
		err = copyToFile(r, info.Preview)
		if err != nil {
			log.Errorf("Can't save a resized image %s: %s\n", info.Filename, err)
		}
	default:
		// Save a file
		info.Type = typeFile
		err := copyToFile(file, info.Origin)
		if err != nil {
			return err
		}
	}

	return fs.storage.addFile(info)
}

// copyToFile copies data from src to new created file
func copyToFile(src io.Reader, path string) error {
	if f, err := os.Open(path); !os.IsNotExist(err) {
		f.Close()
		return ErrAlreadyExist
	}

	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}

	// Write file
	if params.Encrypt {
		_, err = sio.Encrypt(newFile, src, sio.Config{Key: params.Key[:]})
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
// If there was an error during renaming file on a disk, it tries to recover previous filename
func (fs FileStorage) Rename(oldName, newName string) error {
	err := fs.storage.renameFile(oldName, newName)
	if err != nil {
		return errors.Wrapf(err, "can't rename file in a storage\"%s\"", oldName)
	}

	err = os.Rename(params.DataFolder+"/"+oldName, params.DataFolder+"/"+newName)
	if err != nil {
		// Try to recover
		e := fs.storage.renameFile(oldName, newName)
		if e == nil {
			// Success. Return the first error
			return errors.Wrapf(err, "can't rename file \"%s\" on a disk; previous name was recovered", oldName)
		}

		// Return both errors
		return fmt.Errorf("can't rename file on a disk: %s; can't recover previous filename: %s", err, e)
	}

	return nil
}

func (fs FileStorage) ChangeTags(filename string, tags []int) error {
	return fs.storage.updateFileTags(filename, tags)
}

func (fs FileStorage) DeleteTagFromFiles(tagID int) {
	fs.storage.deleteTagFromFiles(tagID)
}

func (fs FileStorage) ChangeDescription(filename, newDescription string) error {
	return fs.storage.updateFileDescription(filename, newDescription)
}

func (fs FileStorage) Delete(filename string) error {
	return fs.storage.deleteFile(filename)
}

func (fs FileStorage) DeleteForce(filename string) error {
	file, err := fs.storage.getFile(filename)
	if err != nil {
		return err
	}

	err = fs.storage.deleteFileForce(filename)
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
			log.Errorf("Can't delete a resized image %s: %s", file.Filename, err)
		}
	}

	return nil
}

// scheduleDeleting deletes files with expired TimeToDelete
// It has to be run in goroutine
func (fs FileStorage) scheduleDeleting() {
	ticker := time.NewTicker(time.Hour * 12)

	for ; true; <-ticker.C {
		log.Infoln("Delete old files")

		var err error
		for _, filename := range fs.storage.getExpiredDeletedFiles() {
			err = fs.DeleteForce(filename)
			if err != nil {
				log.Errorf("Can't delete file \"%s\": %s\n", filename, err)
			} else {
				log.Infof("File \"%s\" was successfully deleted\n", filename)
			}
		}
	}
}

func (fs FileStorage) Recover(filename string) {
	fs.storage.recover(filename)
}
