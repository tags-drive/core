package files

import (
	"sort"

	"github.com/tags-drive/core/cmd"
)

func sortFiles(s cmd.FilesSortMode, files []cmd.FileInfo) {
	switch s {
	case cmd.SortByNameAsc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Filename < files[j].Filename
		})
	case cmd.SortByNameDesc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Filename > files[j].Filename
		})
	case cmd.SortByTimeAsc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].AddTime.Unix() < files[j].AddTime.Unix()
		})
	case cmd.SortByTimeDesc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].AddTime.Unix() > files[j].AddTime.Unix()
		})
	case cmd.SortBySizeAsc:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size < files[j].Size
		})
	case cmd.SortBySizeDecs:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size > files[j].Size
		})
	}
}
