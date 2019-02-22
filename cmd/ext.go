package cmd

// Ext is a struct which contains type of the original file and type for preview
type Ext struct {
	Ext         string   `json:"ext"`
	FileType    FileType `json:"fileType"`
	Supported   bool     `json:"supported"`
	PreviewType FileType `json:"previewType"`
}

type FileType string

// File types for Ext
const (
	FileTypeArchive  FileType = "archive"
	FileTypeAudio    FileType = "audio"
	FileTypeImage    FileType = "image"
	FileTypeLanguage FileType = "lang"
	FileTypeText     FileType = "text"
	FileTypeVideo    FileType = "video"

	MediaTypeAudioMP3  FileType = "audio/mpeg"
	MediaTypeAudioOGG  FileType = "audio/ogg"
	MediaTypeAudioWAV  FileType = "audio/wav"
	MediaTypeVideoMP4  FileType = "video/mp4"
	MediaTypeVideoWebM FileType = "video/webm"

	TypeUnsupported FileType = "unsupported"
)
