package storage_test

import (
	"testing"
	"time"

	"github.com/ShoshinNikita/tags-drive/internal/storage"
)

func TestIsGoodFile(t *testing.T) {
	tests := []struct {
		m     storage.TagMode
		fTags []string
		pTags []string
		res   bool
	}{
		{storage.ModeAnd, []string{"a", "b", "c"}, []string{"a", "c"}, true},
		{storage.ModeAnd, []string{"a", "b", "c"}, []string{"a", "e"}, false},
		{storage.ModeOr, []string{"a", "b", "c"}, []string{"a", "e"}, true},
		{storage.ModeOr, []string{"a", "b", "c"}, []string{"f", "e"}, false},
		{storage.ModeNot, []string{"p", "b", "c"}, []string{"a", "e"}, true},
		{storage.ModeNot, []string{"a", "b", "c"}, []string{"a", "e"}, false},
		// Empty file tags
		{storage.ModeAnd, []string{}, []string{"a", "e"}, false},
		{storage.ModeOr, []string{}, []string{"a", "e"}, false},
		{storage.ModeNot, []string{}, []string{"a", "e"}, true},
		// Empty passed tags
		{storage.ModeAnd, []string{"a", "b", "c"}, []string{}, true},
		{storage.ModeOr, []string{"a", "b", "c"}, []string{}, true},
		{storage.ModeNot, []string{"a", "b", "c"}, []string{}, true},
	}

	for i, tt := range tests {
		res := storage.IsGoodFile(tt.m, tt.fTags, tt.pTags)
		if res != tt.res {
			t.Errorf("Test #%d Want: %v Got %v", i, tt.res, res)
		}
	}

}

func TestSortFiles(t *testing.T) {
	getTime := func(s string) time.Time {
		t, _ := time.Parse("01-02-2006 15:04:05", s)
		return t
	}

	isEqual := func(a, b []storage.FileInfo) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i].Filename != b[i].Filename ||
				a[i].AddTime != b[i].AddTime {
				return false
			}
		}

		return true
	}

	tests := []struct {
		s     storage.SortMode
		files []storage.FileInfo
		res   []storage.FileInfo
	}{
		{storage.SortByNameAsc,
			[]storage.FileInfo{
				storage.FileInfo{Filename: "abc"},
				storage.FileInfo{Filename: "cbd"},
				storage.FileInfo{Filename: "aaa"},
				storage.FileInfo{Filename: "fer"},
			},
			[]storage.FileInfo{
				storage.FileInfo{Filename: "aaa"},
				storage.FileInfo{Filename: "abc"},
				storage.FileInfo{Filename: "cbd"},
				storage.FileInfo{Filename: "fer"},
			},
		},
		{storage.SortByNameDesc,
			[]storage.FileInfo{
				storage.FileInfo{Filename: "abc"},
				storage.FileInfo{Filename: "cbd"},
				storage.FileInfo{Filename: "aaa"},
				storage.FileInfo{Filename: "fer"},
			},
			[]storage.FileInfo{
				storage.FileInfo{Filename: "fer"},
				storage.FileInfo{Filename: "cbd"},
				storage.FileInfo{Filename: "abc"},
				storage.FileInfo{Filename: "aaa"},
			},
		},
		{storage.SortByTimeAsc,
			[]storage.FileInfo{
				storage.FileInfo{AddTime: getTime("05-05-2018 15:45:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:22:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:16:35")},
				storage.FileInfo{AddTime: getTime("05-04-2018 15:22:35")},
			},
			[]storage.FileInfo{
				storage.FileInfo{AddTime: getTime("05-04-2018 15:22:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:16:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:22:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:45:35")},
			},
		},
		{storage.SortByTimeDesc,
			[]storage.FileInfo{
				storage.FileInfo{AddTime: getTime("05-05-2018 15:45:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:22:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:16:35")},
				storage.FileInfo{AddTime: getTime("05-04-2018 15:22:35")},
			},
			[]storage.FileInfo{
				storage.FileInfo{AddTime: getTime("05-05-2018 15:45:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:22:35")},
				storage.FileInfo{AddTime: getTime("05-05-2018 15:16:35")},
				storage.FileInfo{AddTime: getTime("05-04-2018 15:22:35")},
			},
		},
	}

	for i, tt := range tests {
		storage.SortFiles(tt.s, tt.files)
		if !isEqual(tt.files, tt.res) {
			t.Errorf("Test #%d Want: %v Got: %v", i, tt.res, tt.files)
		}
	}
}
