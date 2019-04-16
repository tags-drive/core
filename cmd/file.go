package cmd

import (
	"io"
	"mime/multipart"
	"time"
)

// FileStorageInterface provides methods for interactions with files
type FileStorageInterface interface {
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
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Type     Ext    `json:"type"`
	Origin   string `json:"origin"`            // Origin is a path to a file (params.DataFolder/filename)
	Preview  string `json:"preview,omitempty"` // Preview is a path to a resized image (only if Type.FileType == FileTypeImage)
	//
	Tags        []int     `json:"tags"`
	Description string    `json:"description"`
	Size        int64     `json:"size"`
	AddTime     time.Time `json:"addTime"`
	//
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

// Ext is a struct which contains type of the original file and type for preview
type Ext struct {
	Ext         string      `json:"ext"`
	FileType    FileType    `json:"fileType"`
	Supported   bool        `json:"supported"`
	PreviewType PreviewType `json:"previewType"`
}

type FileType string

// File types
const (
	FileTypeUnsupported FileType = "unsupported"

	FileTypeArchive  FileType = "archive"
	FileTypeAudio    FileType = "audio"
	FileTypeImage    FileType = "image"
	FileTypeLanguage FileType = "lang"
	FileTypeText     FileType = "text"
	FileTypeVideo    FileType = "video"
)

type PreviewType string

// Preview types
const (
	PreviewTypeUnsupported PreviewType = ""

	// audio
	PreviewTypeAudioMP3 PreviewType = "audio/mpeg"
	PreviewTypeAudioOGG PreviewType = "audio/ogg"
	PreviewTypeAudioWAV PreviewType = "audio/wav"
	// image
	PreviewTypeImage PreviewType = "image"
	// video
	PreviewTypeVideoMP4  PreviewType = "video/mp4"
	PreviewTypeVideoWebM PreviewType = "video/webm"
	// text
	PreviewTypeText PreviewType = "text"
)
