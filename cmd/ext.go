package cmd

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
