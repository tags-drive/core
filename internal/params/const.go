package params

import "time"

const Version = "v0.6.0"

// Folders
const (
	// DataFolder is a folder where all files are kept
	DataFolder = "./data"
	// ResizedImagesFolder is a folder where all resized images are kept
	ResizedImagesFolder = "./data/resized"
	// FilesJSONFile is a json file with files information
	FilesJSONFile = "./configs/files.json"
	// TokensJSONFile is a json file with list of tokens
	TokensJSONFile = "./configs/tokens.json"
	// TagsJSONFile is a json file with list of tags (with name and color)
	TagsJSONFile = "./configs/tags.json"
)

// Web const vars
const (
	// MaxTokenLife defines the max lifetime of a token (2 months)
	MaxTokenLife = time.Hour * 24 * 60
	// AuthCookieName defines name of cookie that contains token
	AuthCookieName = "auth"
)

// Storage types
const (
	// JSONStorage is used for StorageType
	JSONStorage = "json"
)
