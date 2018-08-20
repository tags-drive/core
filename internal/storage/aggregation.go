package storage

type Mode int

const (
	// AndMode returns files, which have all tags (a && b && x)
	AndMode Mode = iota
	// OrMode returns files, which have at least ont tag (a || b || x)
	OrMode
	//NotMode return files, which have not passed tags
	NotMode
)

// isGoodFile checks if file has (or hasn't) passed tags
//
// We can use nested loop, because number of tags is small
func isGoodFile(m Mode, fileTags, passedTags []string) (res bool) {
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

// getFiles returns slice of FileInfo with passed tags. If tags is an empty slice, function will return all files
func (fs filesData) getFiles(m Mode, tags []string) []FileInfo {
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
func GetWithTags(m Mode, tags []string) []FileInfo {
	return allFiles.getFiles(m, tags)
}

// GetAll returns all files
func GetAll() []FileInfo {
	// We can pass any Mode
	return allFiles.getFiles(AndMode, []string{})
}
