package files

import (
	"sort"

	"github.com/fvbommel/util/sortorder"
)

func sortFiles(s FilesSortMode, files []File) {
	switch s {
	case SortByNameAsc:
		sort.Slice(files, func(i, j int) bool {
			return sortorder.NaturalLess(files[i].Filename, files[j].Filename)
		})
	case SortByNameDesc:
		sort.Slice(files, func(i, j int) bool {
			return !sortorder.NaturalLess(files[i].Filename, files[j].Filename)
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
