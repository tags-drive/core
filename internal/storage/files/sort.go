package files

import (
	"sort"
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
