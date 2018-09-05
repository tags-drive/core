package files

import (
	"encoding/json"
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

var allFiles = filesData{
	info:  make(map[string]FileInfo),
	mutex: new(sync.RWMutex),
}

// Init reads params.Files and decode its data
func Init() error {
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
			log.Infof("File %s doesn't exist. Need to create a new file\n", params.Files)
			f, err = os.OpenFile(params.Files, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return errors.Wrap(err, "can't create a new file")
			}
			// Write empty structure
			json.NewEncoder(f).Encode(allFiles)
			// Can exit because we don't need to decode files from the file
			f.Close()
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.Files)
	}

	defer f.Close()
	err = allFiles.decode(f)
	if err != nil {
		return errors.Wrapf(err, "can't decode data")
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

	return allFiles.add(info)
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

// writeFile writes file on a disk. It encrypts (or doesn't encrypt) the file according to params.Encrypt
func writeFile(dst io.Writer, src io.Reader) (int64, error) {
	if params.Encrypt {
		return sio.Encrypt(dst, src, sio.Config{Key: params.Key[:]})
	}

	return io.Copy(dst, src)
}

// DeleteFile deletes file from structure and from disk
func DeleteFile(filename string) error {
	file, err := allFiles.get(filename)
	if err != nil {
		return err
	}

	err = allFiles.delete(filename)
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
	return allFiles.rename(oldName, newName)
}

// ChangeTags changes the tags
func ChangeTags(filename string, tags []string) error {
	return allFiles.changeTags(filename, tags)
}

// ChangeDescription changes the description of a file
func ChangeDescription(filename, newDescription string) error {
	return allFiles.changeDescription(filename, newDescription)
}
