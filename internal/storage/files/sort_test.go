package files_test

import (
	"testing"
	"time"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/storage/files"
)

func TestSortFiles(t *testing.T) {
	getTime := func(s string) time.Time {
		tm, err := time.Parse("01-02-2006 15:04:05", s)
		if err != nil {
			t.Errorf("Bad time %s", s)
		}
		return tm
	}

	isEqual := func(a, b []cmd.File) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i].Filename != b[i].Filename ||
				a[i].AddTime != b[i].AddTime ||
				a[i].Size != b[i].Size {
				return false
			}
		}

		return true
	}

	tests := []struct {
		s     cmd.FilesSortMode
		files []cmd.File
		res   []cmd.File
	}{
		{cmd.SortByNameAsc,
			[]cmd.File{
				{Filename: "1"},
				{Filename: "100"},
				{Filename: "3"},
				{Filename: "2"},
				{Filename: "21"},
				{Filename: "20"},
			},
			[]cmd.File{
				{Filename: "1"},
				{Filename: "2"},
				{Filename: "3"},
				{Filename: "20"},
				{Filename: "21"},
				{Filename: "100"},
			},
		},
		{cmd.SortByNameAsc,
			[]cmd.File{
				{Filename: "abc"},
				{Filename: "cbd"},
				{Filename: "aaa"},
				{Filename: "fer"},
			},
			[]cmd.File{
				{Filename: "aaa"},
				{Filename: "abc"},
				{Filename: "cbd"},
				{Filename: "fer"},
			},
		},
		{cmd.SortByNameDesc,
			[]cmd.File{
				{Filename: "1"},
				{Filename: "100"},
				{Filename: "3"},
				{Filename: "2"},
				{Filename: "21"},
				{Filename: "20"},
			},
			[]cmd.File{
				{Filename: "100"},
				{Filename: "21"},
				{Filename: "20"},
				{Filename: "3"},
				{Filename: "2"},
				{Filename: "1"},
			},
		},
		{cmd.SortByNameDesc,
			[]cmd.File{
				{Filename: "abc"},
				{Filename: "cbd"},
				{Filename: "aaa"},
				{Filename: "fer"},
			},
			[]cmd.File{
				{Filename: "fer"},
				{Filename: "cbd"},
				{Filename: "abc"},
				{Filename: "aaa"},
			},
		},
		{cmd.SortByTimeAsc,
			[]cmd.File{
				{AddTime: getTime("05-05-2018 15:45:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-04-2018 15:22:35")},
			},
			[]cmd.File{
				{AddTime: getTime("05-04-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:45:35")},
			},
		},
		{cmd.SortByTimeDesc,
			[]cmd.File{
				{AddTime: getTime("05-05-2018 15:45:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-04-2018 15:22:35")},
			},
			[]cmd.File{
				{AddTime: getTime("05-05-2018 15:45:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-04-2018 15:22:35")},
			},
		},
		{cmd.SortBySizeAsc,
			[]cmd.File{
				{Size: 15},
				{Size: 1515},
				{Size: 1885},
				{Size: 1365},
				{Size: 1551561651},
			},
			[]cmd.File{
				{Size: 15},
				{Size: 1365},
				{Size: 1515},
				{Size: 1885},
				{Size: 1551561651},
			},
		},
		{cmd.SortBySizeDecs,
			[]cmd.File{
				{Size: 15},
				{Size: 1515},
				{Size: 1885},
				{Size: 1365},
				{Size: 1551561651},
			},
			[]cmd.File{
				{Size: 1551561651},
				{Size: 1885},
				{Size: 1515},
				{Size: 1365},
				{Size: 15},
			},
		},
	}

	for i, tt := range tests {
		files.SortFiles(tt.s, tt.files)
		if !isEqual(tt.files, tt.res) {
			t.Errorf("Test #%d Want: %v Got: %v", i, tt.res, tt.files)
		}
	}
}
