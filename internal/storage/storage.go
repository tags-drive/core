package storage

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"github.com/minio/sio"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

// Errors
var (
	ErrFileIsNotExist = errors.New("the file doesn't exist")
	ErrAlreadyExist   = errors.New("file already exists")
)

var allFiles = filesData{
	info:  make(map[string]FileInfo),
	mutex: new(sync.RWMutex),
}

// Init reads params.TagsFiles and decode its data
func Init() error {
	f, err := os.OpenFile(params.TagsFile, os.O_RDWR, 0600)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			log.Infof("File %s doesn't exist. Need to create a new file\n", params.TagsFile)
			f, err = os.OpenFile(params.TagsFile, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return errors.Wrap(err, "can't create a new file")
			}
			// Write empty structure
			json.NewEncoder(f).Encode(allFiles)
			// Can exit because we don't need to decode files from the file
			f.Close()
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.TagsFile)
	}

	defer f.Close()
	err = allFiles.decode(f)
	if err != nil {
		return errors.Wrapf(err, "can't decode data")
	}

	return nil
}

// UploadFile tries to upload a new file. If it was successful, the function calls Files.add()
func UploadFile(f *multipart.FileHeader, tags []string) error {
	// Uploading //
	file, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "can't open a file")
	}
	defer file.Close()

	path := params.DataFolder + "/" + f.Filename
	if _, err := os.Open(path); !os.IsNotExist(err) {
		return ErrAlreadyExist
	}

	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	defer newFile.Close()

	_, err = writeFile(newFile, file)
	if err != nil {
		// Deleting of the bad file
		os.Remove(path)
		return errors.Wrap(err, "can't copy a new file")
	}

	// Adding into global list //
	return allFiles.add(FileInfo{Filename: f.Filename, Size: f.Size, AddTime: time.Now(), Tags: tags})
}

// writeFile writes file on a disk. It encrypts (or doesn't encrypt) the file according to params.Encrypt
func writeFile(dst io.Writer, src io.Reader) (int64, error) {
	if params.Encrypt {
		return sio.Encrypt(dst, src, sio.Config{Key: params.Key[:]})
	}

	return io.Copy(dst, src)
}

// DeleteFile deletes file from structure and from disk
func DeleteFile(filename string) error {
	err := allFiles.delete(filename)
	if err != nil {
		return err
	}

	return os.Remove(params.DataFolder + "/" + filename)
}

// RenameFile renames a file
func RenameFile(oldName, newName string) error {
	return allFiles.rename(oldName, newName)
}

// ChangeTags changes the tags
func ChangeTags(filename string, tags []string) error {
	return allFiles.changeTags(filename, tags)
}
