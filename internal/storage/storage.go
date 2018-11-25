package storage

import (
	"io"
	"mime/multipart"

	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
)

// FileStorageInterface provides methods for interactions with files
type FileStorageInterface interface {
	Init() error
	//
	Get(expr string, s files.SortMode, search string) []files.FileInfo
	GetRecent(number int) []files.FileInfo
	Archive(filenames []string) (io.Reader, error)
	//
	Upload(*multipart.FileHeader) error
	//
	Rename(oldName, newName string) error
	ChangeTags(filename string, tags []int) error
	ChangeDescription(filename, newDescription string) error
	//
	Delete(filename string) error
	DeleteForce(filename string) error
	Recover(filename string)
	DeleteTagFromFiles(tagID int)
}

var Files FileStorageInterface

// TagStorageInterface provides methods for interactions with tags
type TagStorageInterface interface {
	Init() error
	//
	GetAll() tags.Tags
	Add(tags.Tag)
	Change(id int, newName, newColor string)
	Delete(id int)
}

var Tags TagStorageInterface

// Init calls files.Init() and tags.Init()
func Init() error {
	err := Files.Init()
	if err != nil {
		return err
	}

	return Tags.Init()
}
