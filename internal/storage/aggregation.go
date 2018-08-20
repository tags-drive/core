package storage

import (
	"sort"
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

// isGoodFile checks if file has (or hasn't) passed tags
//
// We can use nested loop, because number of tags is small
func isGoodFile(m TagMode, fileTags, passedTags []string) (res bool) {
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
func (fs filesData) getFiles(m TagMode, tags []string) []FileInfo {
	if len(tags) == 0 {
		files := make([]FileInfo, len(fs.info))
		i := 0
		for _, v := range fs.info {
			files[i] = v
			i++
		}
		return files
	}

	var files []FileInfo

	for _, v := range fs.info {
		if isGoodFile(m, v.Tags, tags) {
			files = append(files, v)
		}
	}

	return files
}

// GetWithTags returns all files with (or without) passed tags
// For more information, see AndMode, OrMode, NotMode
func GetWithTags(m TagMode, s SortMode, tags []string) []FileInfo {
	files := allFiles.getFiles(m, tags)
	sortFiles(s, files)
	return files
}

// GetAll returns all files
func GetAll(s SortMode) []FileInfo {
	// We can use any Mode
	files := allFiles.getFiles(ModeAnd, []string{})
	sortFiles(s, files)
	return files
}

// GetRecent returns the last uploaded files
//
// Func uses GetAll(TimeDescMode)
func GetRecent(number int) []FileInfo {
	files := GetAll(SortByTimeDesc)
	if len(files) > number {
		files = files[:number]
	}
	return files
}
