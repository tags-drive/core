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

var fileStorage = struct{ storage }{}

// Init inits fileStorage
func Init() error {
	switch params.StorageType {
	case params.JSONStorage:
		fileStorage.storage = &jsonFileStorage{
			info:  make(map[string]FileInfo),
			mutex: new(sync.RWMutex),
		}
	default:
		// Default storage is jsonFileStorage
		fileStorage.storage = &jsonFileStorage{
			info:  make(map[string]FileInfo),
			mutex: new(sync.RWMutex),
		}
	}

	err := fileStorage.init()
	if err != nil {
		return errors.Wrapf(err, "can't init storage")
	}

	go scheduleDeleting()

	return nil
}

// Get returns all files with (or without) passed tags
func Get(parsedExpr string, s SortMode, search string) []FileInfo {
	search = strings.ToLower(search)
	files := fileStorage.getFiles(parsedExpr, search)
	sortFiles(s, files)
	return files
}

// GetRecent returns the last uploaded files
//
// Func uses Get("", SortByTimeDesc, "")
func GetRecent(number int) []FileInfo {
	files := Get("", SortByTimeDesc, "")
	if len(files) > number {
		files = files[:number]
	}
	return files
}

// UploadFile tries to upload a new file. If it was successful, the function calls Files.add()
func UploadFile(f *multipart.FileHeader) error {
	// Uploading //
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
		Origin:   params.DataFolder + "/" + f.Filename}

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
		copyToFile(file, info.Origin)
	}

	return fileStorage.addFile(info)
}

// copyToFile copies data from src to new created file
func copyToFile(src io.Reader, path string) error {
	if _, err := os.Open(path); !os.IsNotExist(err) {
		return ErrAlreadyExist
	}

	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	defer newFile.Close()

	// Write file
	if params.Encrypt {
		_, err = sio.Encrypt(newFile, src, sio.Config{Key: params.Key[:]})
	} else {
		_, err = io.Copy(newFile, src)
	}

	if err != nil {
		// Deleting of the bad file
		os.Remove(path)
		return errors.Wrap(err, "can't copy a new file")
	}

	return nil
}

// RenameFile renames a file
//
// If there was an error during renaming file on a disk, it tries to recover previous filename
//
func RenameFile(oldName, newName string) error {
	err := fileStorage.renameFile(oldName, newName)
	if err != nil {
		return errors.Wrapf(err, "can't rename file in a storage\"%s\"", oldName)
	}

	err = os.Rename(params.DataFolder+"/"+oldName, params.DataFolder+"/"+newName)
	if err != nil {
		// Try to recover
		e := fileStorage.renameFile(oldName, newName)
		if e == nil {
			// Success. Return the first error
			return errors.Wrapf(err, "can't rename file \"%s\" on a disk; previous name was recovered", oldName)
		}

		// Return both errors
		return fmt.Errorf("can't rename file on a disk: %s; can't recover previous filename: %s", err, e)
	}

	return nil
}

// ChangeTags changes the tags
func ChangeTags(filename string, tags []int) error {
	return fileStorage.updateFileTags(filename, tags)
}

// DeleteTag deletes a tag
func DeleteTag(tagID int) {
	fileStorage.deleteTagFromFiles(tagID)
}

// ChangeDescription changes the description of a file
func ChangeDescription(filename, newDescription string) error {
	return fileStorage.updateFileDescription(filename, newDescription)
}

// ArchiveFiles archives passed files and returns io.Reader
func ArchiveFiles(files []string) (body io.Reader, err error) {
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

// DeleteFile calls fileStorage.deleteFile
func DeleteFile(filename string) error {
	return fileStorage.deleteFile(filename)
}

// DeleteFileForce deletes file from structure and from disk
func DeleteFileForce(filename string) error {
	file, err := fileStorage.getFile(filename)
	if err != nil {
		return err
	}

	err = fileStorage.deleteFileForce(filename)
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
func scheduleDeleting() {
	ticker := time.NewTicker(time.Hour * 12)

	for ; true; <-ticker.C {
		log.Infoln("Delete old files")

		var err error
		for _, filename := range fileStorage.getExpiredDeletedFiles() {
			err = DeleteFileForce(filename)
			if err != nil {
				log.Errorf("can't delete file \"%s\": %s\n", filename, err)
			} else {
				log.Infof("file \"%s\" was successfully deleted\n", filename)
			}
		}
	}
}

// RecoverFile "removes" file from Trash
func RecoverFile(filename string) {
	fileStorage.recover(filename)
}
