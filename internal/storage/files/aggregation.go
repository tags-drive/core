package files

import (
	"sort"
	"strings"
)

type TagMode int

const (
	// ModeAnd returns files, which have all tags (a && b && x)
	ModeAnd TagMode = iota
	// ModeOr returns files, which have at least ont tag (a || b || x)
	ModeOr
	//ModeNot return files, which have not passed tags
	ModeNot
)

type SortMode int

const (
	SortByNameAsc SortMode = iota
	SortByNameDesc
	SortByTimeAsc
	SortByTimeDesc
	SortBySizeAsc
	SortBySizeDecs
)

// isGoodFile checks if file has passed tags
//
// We can use nested loop, because number of tags is small
func isGoodFile(m TagMode, fileTags, passedTags []int) (res bool) {
	if len(passedTags) == 0 {
		return true
	}

	switch m {
	case ModeAnd:
		if len(fileTags) == 0 {
			return false
		}
		for _, pt := range passedTags {
			has := false
			for _, ft := range fileTags {
				if pt == ft {
					has = true
					break
				}
			}
			if !has {
				return false
			}
		}
		return true
	case ModeOr:
		if len(fileTags) == 0 {
			return false
		}
		for _, pt := range passedTags {
			for _, ft := range fileTags {
				if pt == ft {
					return true
				}
			}
		}
		return false
	case ModeNot:
		if len(fileTags) == 0 {
			return true
		}
		for _, pt := range passedTags {
			for _, ft := range fileTags {
				if pt == ft {
					return false
				}
			}
		}
		return true
	}

	return false
}

func sortFiles(s SortMode, files []FileInfo) {
	switch s {
	case SortByNameAsc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Filename < files[j].Filename
		})
	case SortByNameDesc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Filename > files[j].Filename
		})
	case SortByTimeAsc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].AddTime.Unix() < files[j].AddTime.Unix()
		})
	case SortByTimeDesc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].AddTime.Unix() > files[j].AddTime.Unix()
		})
	case SortBySizeAsc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size < files[j].Size
		})
	case SortBySizeDecs:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size > files[j].Size
		})
	}
}

// getFiles returns slice of FileInfo with passed tags. If tags is an empty slice, function will return all files
func (fs filesData) getFiles(m TagMode, tags []int, search string) (files []FileInfo) {
	fs.mutex.RLock()
	if len(tags) == 0 {
		files = make([]FileInfo, len(fs.info))
		i := 0
		for _, v := range fs.info {
			files[i] = v
			i++
		}
	} else {
		for _, v := range fs.info {
			if isGoodFile(m, v.Tags, tags) {
				files = append(files, v)
			}
		}
	}

	fs.mutex.RUnlock()

	if search == "" {
		return files
	}

	// Need to remove files with incorrect name
	var goodFiles []FileInfo
	for _, f := range files {
		if strings.Contains(f.Filename, search) {
			goodFiles = append(goodFiles, f)
		}
	}

	return goodFiles
}

// Get returns all files with (or without) passed tags
// For more information, see AndMode, OrMode, NotMode
func Get(m TagMode, s SortMode, tags []int, search string) []FileInfo {
	files := allFiles.getFiles(m, tags, search)
	sortFiles(s, files)
	return files
}

// GetRecent returns the last uploaded files
//
// Func uses GetAll(TimeDescMode)
func GetRecent(number int) []FileInfo {
	files := Get(ModeAnd, SortByTimeDesc, []int{}, "")
	if len(files) > number {
		files = files[:number]
	}
	return files
}
