package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tags-drive/internal/params"
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
		return errors.Wrapf(err, "can't open file")
	}

	path := params.DataFolder + "/" + f.Filename
	// TODO
	if _, err := os.Open(path); !os.IsNotExist(err) {
		return fmt.Errorf("file %s already exists", path)
	}

	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrapf(err, "can't create a new file %s\n", path)
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		// Deleting of the bad file
		os.Remove(path)
		return errors.Wrap(err, "Can't copy a new file")
	}

	// Adding into global list //
	return allFiles.add(FileInfo{Filename: f.Filename, Size: f.Size, AddTime: time.Now(), Tags: tags})
}

// DeleteFile deletes file from structure and from disk
func DeleteFile(filename string) error {
	err := allFiles.delete(filename)
	if err != nil {
		return err
	}

	return os.Remove(params.DataFolder + "/" + filename)
}
