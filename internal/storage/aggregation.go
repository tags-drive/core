package storage

import (
	"sort"
	"time"
)

type TagMode int

const (
	// AndMode returns files, which have all tags (a && b && x)
	AndMode TagMode = iota
	// OrMode returns files, which have at least ont tag (a || b || x)
	OrMode
	//NotMode return files, which have not passed tags
	NotMode
)

type SortMode int

const (
	NameAscMode SortMode = iota
	NameDescMode
	TimeAscMode
	TimeDescMode

	// TODO add SizeAscMode and SizeDescMode
)

// isGoodFile checks if file has (or hasn't) passed tags
//
// We can use nested loop, because number of tags is small
func isGoodFile(m TagMode, fileTags, passedTags []string) (res bool) {
	if len(passedTags) == 0 {
		return true
	}

	switch m {
	case AndMode:
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
	case OrMode:
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
	case NotMode:
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
	case NameAscMode:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Filename < files[j].Filename
		})
	case NameDescMode:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Filename > files[j].Filename
		})
	case TimeAscMode:
		sort.Slice(files, func(i, j int) bool {
			return files[i].AddTime.Before(files[j].AddTime) // before == <
		})
	case TimeDescMode:
		sort.Slice(files, func(i, j int) bool {
			return files[i].AddTime.After(files[j].AddTime) // after == >
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
	files := allFiles.getFiles(AndMode, []string{})
	sortFiles(s, files)
	return files
}

// GetRecent returns the last uploaded files
//
// Func uses GetAll(TimeDescMode)
func GetRecent(number int, maxAge time.Time) []FileInfo {
	files := GetAll(TimeDescMode)
	if len(files) > number {
		files = files[:number]
	}
	return files
}
