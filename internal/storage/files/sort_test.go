package files

import (
	"testing"
	"time"
)

func TestSortFiles(t *testing.T) {
	getTime := func(s string) time.Time {
		tm, err := time.Parse("01-02-2006 15:04:05", s)
		if err != nil {
			t.Errorf("Bad time %s", s)
		}
		return tm
	}

	isEqual := func(a, b []File) bool {
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
		s     FilesSortMode
		files []File
		res   []File
	}{
		{SortByNameAsc,
			[]File{
				{Filename: "1"},
				{Filename: "100"},
				{Filename: "3"},
				{Filename: "2"},
				{Filename: "21"},
				{Filename: "20"},
			},
			[]File{
				{Filename: "1"},
				{Filename: "2"},
				{Filename: "3"},
				{Filename: "20"},
				{Filename: "21"},
				{Filename: "100"},
			},
		},
		{SortByNameAsc,
			[]File{
				{Filename: "abc"},
				{Filename: "cbd"},
				{Filename: "aaa"},
				{Filename: "fer"},
			},
			[]File{
				{Filename: "aaa"},
				{Filename: "abc"},
				{Filename: "cbd"},
				{Filename: "fer"},
			},
		},
		{SortByNameDesc,
			[]File{
				{Filename: "1"},
				{Filename: "100"},
				{Filename: "3"},
				{Filename: "2"},
				{Filename: "21"},
				{Filename: "20"},
			},
			[]File{
				{Filename: "100"},
				{Filename: "21"},
				{Filename: "20"},
				{Filename: "3"},
				{Filename: "2"},
				{Filename: "1"},
			},
		},
		{SortByNameDesc,
			[]File{
				{Filename: "abc"},
				{Filename: "cbd"},
				{Filename: "aaa"},
				{Filename: "fer"},
			},
			[]File{
				{Filename: "fer"},
				{Filename: "cbd"},
				{Filename: "abc"},
				{Filename: "aaa"},
			},
		},
		{SortByTimeAsc,
			[]File{
				{AddTime: getTime("05-05-2018 15:45:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-04-2018 15:22:35")},
			},
			[]File{
				{AddTime: getTime("05-04-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:45:35")},
			},
		},
		{SortByTimeDesc,
			[]File{
				{AddTime: getTime("05-05-2018 15:45:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-04-2018 15:22:35")},
			},
			[]File{
				{AddTime: getTime("05-05-2018 15:45:35")},
				{AddTime: getTime("05-05-2018 15:22:35")},
				{AddTime: getTime("05-05-2018 15:16:35")},
				{AddTime: getTime("05-04-2018 15:22:35")},
			},
		},
		{SortBySizeAsc,
			[]File{
				{Size: 15},
				{Size: 1515},
				{Size: 1885},
				{Size: 1365},
				{Size: 1551561651},
			},
			[]File{
				{Size: 15},
				{Size: 1365},
				{Size: 1515},
				{Size: 1885},
				{Size: 1551561651},
			},
		},
		{SortBySizeDecs,
			[]File{
				{Size: 15},
				{Size: 1515},
				{Size: 1885},
				{Size: 1365},
				{Size: 1551561651},
			},
			[]File{
				{Size: 1551561651},
				{Size: 1885},
				{Size: 1515},
				{Size: 1365},
				{Size: 15},
			},
		},
	}

	for i, tt := range tests {
		sortFiles(tt.s, tt.files)
		if !isEqual(tt.files, tt.res) {
			t.Errorf("Test #%d Want: %v Got: %v", i, tt.res, tt.files)
		}
	}
}
