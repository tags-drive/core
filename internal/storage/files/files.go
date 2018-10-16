package files

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/minio/sio"
	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tags-drive/internal/params"
	"github.com/ShoshinNikita/tags-drive/internal/storage/files/resizing"
)

const (
	typeImage = "image"
	typeFile  = "file"
)

// Errors
var (
	ErrFileIsNotExist = errors.New("the file doesn't exist")
	ErrAlreadyExist   = errors.New("file already exists")
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

type storage interface {
	init() error

	// getFile returns a file with passed filename
	getFile(filename string) (FileInfo, error)

	// getFiles returns files
	//     tagMode - mode of tags
	//     tags - list of needed tags
	//     search - string, which filename has to contain
	getFiles(m TagMode, tags []int, search string) (files []FileInfo)

	// add adds a file
	addFile(info FileInfo) error

	// renameFile renames a file
	renameFile(oldName string, newName string) error

	// updateFileTags updates tags of a file
	updateFileTags(filename string, changedTagsID []int) error

	// updateFileDescription update description of a file
	updateFileDescription(filename string, newDesc string) error

	// deleteFile deletes a file from list (not from disk)
	deleteFile(filename string) error

	// deleteTagFromFiles deletes a tag (it's called when user deletes a tag)
	deleteTagFromFiles(tagID int)
}

var fileStorage = struct {
	storage
}{}

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

	return nil
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
		err = upload(r, info.Origin)
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
		err = upload(r, info.Preview)
		if err != nil {
			log.Errorf("Can't save a resized image %s: %s\n", info.Filename, err)
		}
	default:
		// Save a file
		info.Type = typeFile
		upload(file, info.Origin)
	}

	return fileStorage.addFile(info)
}

func upload(src io.Reader, path string) error {
	if _, err := os.Open(path); !os.IsNotExist(err) {
		return ErrAlreadyExist
	}

	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	defer newFile.Close()

	_, err = writeFile(newFile, src)
	if err != nil {
		// Deleting of the bad file
		os.Remove(path)
		return errors.Wrap(err, "can't copy a new file")
	}

	return nil
}

// writeFile writes file into dst. It encrypts (or doesn't encrypt) the file according to params.Encrypt
func writeFile(dst io.Writer, src io.Reader) (int64, error) {
	if params.Encrypt {
		return sio.Encrypt(dst, src, sio.Config{Key: params.Key[:]})
	}

	return io.Copy(dst, src)
}

// DeleteFile deletes file from structure and from disk
func DeleteFile(filename string) error {
	file, err := fileStorage.getFile(filename)
	if err != nil {
		return err
	}

	err = fileStorage.deleteFile(filename)
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

// RenameFile renames a file
func RenameFile(oldName, newName string) error {
	// At first, rename file on disk
	err := os.Rename(params.DataFolder+"/"+oldName, params.DataFolder+"/"+newName)
	if err != nil {
		return err
	}

	return fileStorage.renameFile(oldName, newName)
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
			log.Errorf("Can't load file %s\n", filename)
			continue
		}
		stat, err := f.Stat()
		if err != nil {
			log.Errorf("Can't load file %s\n", filename)
			continue
		}

		header, err := zip.FileInfoHeader(stat)
		header.Method = zip.Deflate

		wr, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Errorf("Can't load file %s\n", filename)
			f.Close()
			continue
		}

		if params.Encrypt {
			_, err = sio.Decrypt(wr, f, sio.Config{Key: params.Key[:]})
		} else {
			_, err = io.Copy(wr, f)
		}

		if err != nil {
			log.Errorf("Can't load file %s\n", filename)
		}

		f.Close()
	}

	return buff, nil
}
