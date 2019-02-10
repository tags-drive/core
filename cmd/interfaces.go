package cmd

import (
	"io"
	"mime/multipart"
)

// Server provides methods for interactions web server
type Server interface {
	Start() error

	// Shutdown gracefully shutdown server
	Shutdown() error
}

// FileStorageInterface provides methods for interactions with files
type FileStorageInterface interface {
	// Get returns all "good" sorted files
	//
	// If expr isn't valid, Get returns ErrBadExpessionSyntax
	// count must be greater than 0, else all files will be returned ([offset:])
	Get(expr string, s FilesSortMode, search string, offset, count int) ([]FileInfo, error)
	// GetFile returns a file with passed id
	GetFile(id int) (FileInfo, error)
	// GetRecent returns the last uploaded files
	GetRecent(number int) []FileInfo
	// ArchiveFiles archives passed files and returns io.Reader with archive
	Archive(fileIDs []int) (io.Reader, error)

	// UploadFile uploads a new file
	Upload(file *multipart.FileHeader, tags []int) error

	// Rename renames a file
	Rename(fileID int, newName string) (updatedFile FileInfo, err error)
	// ChangeTags changes the tags
	ChangeTags(fileID int, tags []int) (updatedFile FileInfo, err error)
	// ChangeDescription changes the description
	ChangeDescription(fileID int, newDescription string) (updatedFile FileInfo, err error)

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

// TagStorageInterface provides methods for interactions with tags
type TagStorageInterface interface {
	// Get return tag with passed id. If a tag doesn't exist, it returns Tag{}, false
	Get(id int) (Tag, bool)

	// GetAll returns all tags
	GetAll() Tags

	// Add adds a new tag with passed name and color
	Add(name, color string)

	// Change changes a tag with passed id.
	// If pass empty newName (or newColor), field Name (or Color) won't be changed.
	Change(id int, newName, newColor string)

	// Delete deletes a tag with passed id
	Delete(id int)

	// Check checks is there tag with passed id
	Check(id int) bool

	// Shutdown gracefully shutdown TagStorage
	Shutdown() error
}

// AuthService provides methods for auth users
type AuthService interface {
	// GenerateToken generates a new token. GenerateToken doesn't add new token, just return it!
	GenerateToken() string

	// AddToken adds passed token into storage
	AddToken(token string)

	// CheckToken returns true if token is in storage
	CheckToken(token string) bool

	// DeleteToken deletes token from a storage
	DeleteToken(token string)

	// Shutdown gracefully shutdown FileStorage
	Shutdown() error
}
