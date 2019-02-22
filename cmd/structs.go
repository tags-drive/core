package cmd

import "time"

// Tag contains the information about a tag
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Tags is a map of Tag
type Tags map[int]Tag

// Ext is a struct which contains type of the original file and type for preview
type Ext struct {
	Ext         string   `json:"ext"`
	FileType    FileType `json:"fileType"`
	Supported   bool     `json:"supported"`
	PreviewType FileType `json:"previewType"`
}

type FileType string

// FileInfo contains the information about a file
type FileInfo struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Type     string `json:"type"`              // typeImage or typeFile
	Origin   string `json:"origin"`            // Origin is a path to a file (params.DataFolder/filename)
	Preview  string `json:"preview,omitempty"` // Preview is a path to a resized image (only if Type == TypeImage)
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
