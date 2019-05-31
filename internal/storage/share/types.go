package share

import (
	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
)

type Config struct {
	ShareTokenJSONFile string

	Encrypt    bool
	PassPhrase [32]byte
}

type ShareStorageInterface interface {
	// CreateToken creates new token with access to passed files
	CreateToken(fileIDs []int) (token string)

	// CheckToken checks if a token exists
	CheckToken(token string) bool

	// CheckFile checks if a token grants access to a file
	CheckFile(token string, id int) bool

	// DeleteFile deletes all refs to a file
	DeleteFile(id int)

	// FilterFiles filters files according to token share permissions
	FilterFiles(token string, files []files.File) ([]files.File, error)

	// TODO
	// FilterTags filters tags according to token share permissions
	FilterTags(token string, tags tags.Tags) (tags.Tags, error)

	Shutdown() error
}
