package files

import (
	"io"
	"os"
	"time"

	"errors"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
	"github.com/tags-drive/core/internal/storage/files/extensions"
)

type Config struct {
	Debug bool

	VarFolder           string
	DataFolder          string
	ResizedImagesFolder string
	// A file is deleted from the storage and from a disk after this time since user add the file into the Trash
	TimeBeforeDeleting time.Duration

	MetadataStorageType string
	FilesJSONFile       string

	Encrypt    bool
	PassPhrase [32]byte
}

type FilterFilesFunction func([]File) ([]File, error)

type GetFilesConfig struct {
	Expr     string
	SortMode FilesSortMode
	Search   string
	IsRegexp bool
	Offset   int
	Count    int                 // count must be greater than 0, else all files will be returned ([offset:])
	Filter   FilterFilesFunction // can be nil
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

// metadataStorage is a storage for the files metadata
type metadataStorage interface {
	init() error

	// getFile returns a file with passed filename
	getFile(id int) (File, error)

	checkFile(id int) bool

	// getFiles returns files
	//     expr - parsed logical expression
	//     search - string, which filename has to contain (lower case)
	//     isRegexp - is expr a regular expression (if it is true, expr must be valid regular expression)
	getFiles(expr aggregation.LogicalExpr, search string, isRegexp bool) (files []File)

	getFilesWithIDs(ids ...int) []File

	// add adds a file
	addFile(filename string, fileType extensions.Ext, tags []int, size int64, addTime time.Time) (id int)

	// renameFile renames a file
	renameFile(id int, newName string) (File, error)

	// updateFileTags updates tags of a file
	updateFileTags(id int, changedTagsID []int) (File, error)

	// updateFileDescription update description of a file
	updateFileDescription(id int, newDesc string) (File, error)

	// deleteFile marks file deleted and sets TimeToDelete
	// File can't be deleted several times (function should return ErrFileDeletedAgain)
	deleteFile(id int) error

	// deleteFileForce deletes file
	deleteFileForce(id int) error

	// recover removes file from Trash
	recover(id int)

	// addTagsToFiles adds a tag to files
	addTagsToFiles(filesIDs, tagsID []int)

	// removeTagsFromFiles removes tags from selected files
	removeTagsFromFiles(filesIDs, tagsID []int)

	// deleteTagFromFiles deletes a tag
	removeTagFromAllFiles(tagID int)

	// getExpiredDeletedFiles returns names of files with expired TimeToDelete
	getExpiredDeletedFiles() []int

	shutdown() error
}

// binaryStorage is a storage for the files themselves
type binaryStorage interface {
	// GetFile writes a file into passed io.Writer
	GetFile(w io.Writer, fileID int, resized bool) error

	GetFileStats(fileID int) (os.FileInfo, error)

	SaveFile(r io.Reader, fileID int, resized bool) error

	DeleteFile(fileID int, resized bool) error
}

type binaryStorageMock struct{}

var mockError = errors.New("mock storage is used")

func (_ binaryStorageMock) GetFile(w io.Writer, fileID int, resized bool) error  { return mockError }
func (_ binaryStorageMock) GetFileStats(fileID int) (os.FileInfo, error)         { return nil, mockError }
func (_ binaryStorageMock) SaveFile(r io.Reader, fileID int, resized bool) error { return mockError }
func (_ binaryStorageMock) DeleteFile(fileID int, resized bool) error            { return mockError }
