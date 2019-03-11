package cmd

import "time"

// FileInfo contains the information about a file
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
