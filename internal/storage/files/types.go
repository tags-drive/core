package files

import (
	"io"
	"mime/multipart"
	"time"

	"github.com/tags-drive/core/internal/storage/files/extensions"
)

type Config struct {
	Debug bool

	DataFolder          string
	ResizedImagesFolder string
	// A file is deleted from the storage and from a disk after this time since user add the file into the Trash
	TimeBeforeDeleting time.Duration

	StorageType   string
	FilesJSONFile string

	Encrypt    bool
	PassPhrase [32]byte
}

// FileStorageInterface provides methods for interactions with files
type FileStorageInterface interface {
	// Start starts all background services
	StartBackgroundServices()

	// Get returns all "good" sorted files
	//
	// If expr isn't valid, Get returns ErrBadExpessionSyntax
	// count must be greater than 0, else all files will be returned ([offset:])
	Get(expr string, s FilesSortMode, search string, isRegexp bool, offset, count int) ([]File, error)
	// GetFile returns a file with passed id
	GetFile(id int) (File, error)
	// GetRecent returns the last uploaded files
	GetRecent(number int) []File
	// ArchiveFiles archives passed files and returns io.Reader with archive
	Archive(fileIDs []int) (io.Reader, error)

	// UploadFile uploads a new file
	Upload(file *multipart.FileHeader, tags []int) error

	// Rename renames a file
	Rename(fileID int, newName string) (updatedFile File, err error)
	// ChangeTags changes the tags
	ChangeTags(fileID int, tags []int) (updatedFile File, err error)
	// ChangeDescription changes the description
	ChangeDescription(fileID int, newDescription string) (updatedFile File, err error)

	// Delete "move" a file into Trash
	Delete(fileID int) error
	// DeleteForce deletes file from storage and from disk
	DeleteForce(fileID int) error
	// Recover "removes" file from Trash
	Recover(fileID int)

	// AddTagsToFiles adds a tag to files
	AddTagsToFiles(filesIDs, tagsIDs []int)
	// RemoveTagsFromFiles
	RemoveTagsFromFiles(filesIDs, tagsIDs []int)

	// DeleteTagFromFiles deletes a tag from files
	DeleteTagFromFiles(tagID int)

	// Shutdown gracefully shutdown FileStorage
	Shutdown() error
}

// File contains the information about a file
type File struct {
	ID       int            `json:"id"`
	Filename string         `json:"filename"`
	Type     extensions.Ext `json:"type"`
	Origin   string         `json:"origin"`            // Origin is a URL address of a file
	Preview  string         `json:"preview,omitempty"` // Preview is a URL address of a resized image (only if Type.FileType == FileTypeImage)

	Tags        []int     `json:"tags"`
	Description string    `json:"description"`
	Size        int64     `json:"size"`
	AddTime     time.Time `json:"addTime"`

	Deleted      bool      `json:"deleted"`
	TimeToDelete time.Time `json:"timeToDelete"`
}

type FilesSortMode int

const (
	SortByNameAsc FilesSortMode = iota
	SortByNameDesc
	SortByTimeAsc
	SortByTimeDesc
	SortBySizeAsc
	SortBySizeDecs
)
