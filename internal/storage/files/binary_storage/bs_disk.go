package bs

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/minio/sio"
	"github.com/pkg/errors"
)

type DiskStorage struct {
	config DiskStorageConfig
}

type DiskStorageConfig struct {
	DataFolder          string
	ResizedImagesFolder string

	Encrypt    bool
	PassPhrase [32]byte
}

func NewDiskStorage(cnf DiskStorageConfig) (*DiskStorage, error) {
	// Create folders
	if err := os.MkdirAll(cnf.DataFolder, 0666); err != nil {
		return nil, errors.Wrap(err, "can't create DataFolder")
	}

	if err := os.MkdirAll(cnf.ResizedImagesFolder, 0666); err != nil {
		return nil, errors.Wrap(err, "can't create ResizedImagesFolder")
	}

	// Normalize paths
	if !strings.HasSuffix(cnf.DataFolder, "/") {
		cnf.DataFolder += "/"
	}
	if !strings.HasSuffix(cnf.ResizedImagesFolder, "/") {
		cnf.ResizedImagesFolder += "/"
	}

	return &DiskStorage{
		config: cnf,
	}, nil
}

func (ds DiskStorage) GetFile(w io.Writer, fileID int, resized bool) error {
	path := ds.getFilePath(fileID, resized)

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "can't open the file '%s'", path)
	}
	defer f.Close()

	copier := io.Copy
	if ds.config.Encrypt {
		copier = func(dst io.Writer, src io.Reader) (int64, error) {
			return sio.Decrypt(dst, src, sio.Config{Key: ds.config.PassPhrase[:]})
		}
	}

	_, err = copier(w, f)
	if err != nil {
		return errors.Wrapf(err, "can't copy the file '%s' into io.Writer", path)
	}

	return nil
}

func (ds DiskStorage) GetFileStats(fileID int) (os.FileInfo, error) {
	path := ds.getFilePath(fileID, false)
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open file '%s'", path)
	}

	// Can't use defer for file.Close()!
	stats, err := f.Stat()
	f.Close()
	return stats, err
}

func (ds DiskStorage) SaveFile(r io.Reader, fileID int, resized bool) error {
	path := ds.getFilePath(fileID, resized)

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "can't create a new file '%s'", path)
	}
	// TODO: delete file if an error occurred?
	defer f.Close()

	if ds.config.Encrypt {
		_, err := sio.Encrypt(f, r, sio.Config{Key: ds.config.PassPhrase[:]})
		return errors.Wrapf(err, "can't encrypt a file '%s'", err)
	}

	_, err = io.Copy(f, r)
	return errors.Wrapf(err, "can't copy io.Reader into a file '%s'", err)
}

func (ds DiskStorage) DeleteFile(fileID int, resized bool) error {
	path := ds.getFilePath(fileID, resized)
	return os.Remove(path)
}

func (ds DiskStorage) getFilePath(id int, resized bool) string {
	path := ds.config.DataFolder
	if resized {
		path = ds.config.ResizedImagesFolder
	}
	path += strconv.Itoa(id)
	return path
}
