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

type FileStorage interface {
	GetFiles(ids ...int) []files.File
}

type ShareStorageInterface interface {
	// GetAllTokens returns all tokens with shared files ids
	GetAllTokens() map[string][]int

	// GetFilesIDs returns files shared by a passed token
	GetFilesIDs(token string) (filesIDs []int, err error)

	// CheckToken checks if a token exists
	CheckToken(token string) bool

	// CreateToken creates new token with access to passed files
	CreateToken(filesIDs []int) (token string)

	// DeleteToken delete a share token
	DeleteToken(token string)

	// CheckFile checks if a token grants access to a file
	CheckFile(token string, id int) bool

	// DeleteFile deletes all refs to a file
	DeleteFile(id int)

	// FilterFiles filters files according to token share permissions
	FilterFiles(token string, files []files.File) ([]files.File, error)

	// FilterTags filters tags according to token share permissions
	FilterTags(token string, tags tags.Tags) (tags.Tags, error)

	Shutdown() error
}
