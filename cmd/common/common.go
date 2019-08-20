// Package common contains paths to files, folder and etc. These paths are used among
// all applications in /cmd.
package common

const (
	// VarFolder is the main folder. All files are kept here.
	// DatFolder and ResizedImagesFolder must be subfolders of this directory.
	VarFolder           = "./var"
	DataFolder          = "./var/data"
	ResizedImagesFolder = "./var/data/resized"
	//
	DataBucket          = "var-data"
	ResizedImagesBucket = "var-data-resized"

	FilesJSONFile       = "./var/files.json"        // for files
	TagsJSONFile        = "./var/tags.json"         // for tags
	AuthTokensJSONFile  = "./var/auth_tokens.json"  // for auth tokens
	ShareTokensJSONFile = "./var/share_tokens.json" // for share tokens
)
