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

type internalStorage interface {
	// GetAllTokens returns all tokens with shared files ids
	getAllTokens() map[string][]int

	// GetFilesIDs returns files shared by a passed token
	getFilesIDs(token string) (filesIDs []int, err error)

	// CheckToken checks if a token exists
	checkToken(token string) bool

	// CreateToken creates new token with access to passed files
	createToken(filesIDs []int) (token string)

	// DeleteToken delete a share token
	deleteToken(token string)

	// CheckFile checks if a token grants access to a file
	checkFile(token string, id int) bool

	// DeleteFile deletes all refs to a file
	deleteFile(id int)

	// FilterFiles filters files according to token share permissions
	filterFiles(token string, files []files.File) ([]files.File, error)

	// FilterTags filters tags according to token share permissions
	filterTags(token string, tags tags.Tags) (tags.Tags, error)

	shutdown() error
}
