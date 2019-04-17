package files

import (
	"os"
	"testing"
	"time"

	clog "github.com/ShoshinNikita/log/v2"

	"github.com/tags-drive/core/internal/storage/files/extensions"
)

func areArraysEqualInt(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	m := make(map[int]struct{}, len(a))
	// Add elements
	for _, v := range a {
		m[v] = struct{}{}
	}
	// Check
	for _, v := range b {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}

func areArraysEqualString(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	m := make(map[string]struct{}, len(a))
	// Add elements
	for _, v := range a {
		m[v] = struct{}{}
	}
	// Check
	for _, v := range b {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}

func newStorage() *jsonFileStorage {
	return newJsonFileStorage(clog.NewProdLogger())
}

// addDefaultFiles adds next files into storage:
// | id | name | tags      |
// | -- | ---- | --------- |
// | 1  | 1    | [1, 2, 3] |
// | 2  | 2    | [1, 2, 7] |
// | 3  | 3    | [3]       |
// | 4  | 4    | [2, 3]    |
// | 5  | 5    | [4, 5, 6] |
// | 6  | 6    | []        |
//
func addDefaultFiles(storage *jsonFileStorage) {
	files := []struct {
		filename string
		tags     []int
	}{
		{"1", []int{1, 2, 3}},
		{"2", []int{1, 2, 7}},
		{"3", []int{3}},
		{"4", []int{2, 3}},
		{"5", []int{4, 5, 6}},
		{"6", []int{}},
	}

	now := time.Now()

	for _, f := range files {
		storage.addFile(f.filename, extensions.Ext{}, f.tags, 0, now)
	}
}

// removeConfigFile remove "configs/files.json"
func removeConfigFile() {
	os.Remove("configs/files.json")
}

func TestMain(m *testing.M) {
	// Create "configs"folder
	err := os.Mkdir("configs", 0666)
	if err != nil {
		clog.Fatalln("can't create folder \"configs\":", err)
	}

	code := m.Run()

	// Remove folder "configs"
	err = os.RemoveAll("configs")
	if err != nil {
		clog.Errorln("can't remove folder \"configs\":", err)
	}
	// Remove "data" folder created by storage.init()
	err = os.RemoveAll("data")
	if err != nil {
		clog.Errorln("can't remove folder \"configs\":", err)
	}

	os.Exit(code)
}

func TestAddTagsToFiles(t *testing.T) {
	storage := newStorage()
	storage.init()
	addDefaultFiles(storage)

	tests := []struct {
		files  []int
		tags   []int
		result [][]int // (index + 1) – id of a file
	}{
		{
			files: []int{1, 2, 3, 4, 5, 6},
			tags:  []int{1},
			result: [][]int{
				{1, 2, 3},    // 1
				{1, 2, 7},    // 2
				{1, 3},       // 3
				{1, 2, 3},    // 4
				{1, 4, 5, 6}, // 5
				{1},          // 6
			},
		},
		{
			files: []int{4, 6},
			tags:  []int{5},
			result: [][]int{
				{1, 2, 3},    // 1
				{1, 2, 7},    // 2
				{1, 3},       // 3
				{1, 2, 3, 5}, // 4
				{1, 4, 5, 6}, // 5
				{1, 5},       // 6
			},
		},
		{
			files: []int{1},
			tags:  []int{7, 8, 9},
			result: [][]int{
				{1, 2, 3, 7, 8, 9}, // 1
				{1, 2, 7},          // 2
				{1, 3},             // 3
				{1, 2, 3, 5},       // 4
				{1, 4, 5, 6},       // 5
				{1, 5},             // 6
			},
		},
		{
			files: []int{1, 3, 5},
			tags:  []int{20, 30, 40},
			result: [][]int{
				{1, 2, 3, 7, 8, 9, 20, 30, 40}, // 1
				{1, 2, 7},                      // 2
				{1, 3, 20, 30, 40},             // 3
				{1, 2, 3, 5},                   // 4
				{1, 4, 5, 6, 20, 30, 40},       // 5
				{1, 5},                         // 6
			},
		},
		{
			files: []int{1, 2, 3, 4, 5, 6},
			tags:  []int{1, 2, 3, 4},
			result: [][]int{
				{1, 2, 3, 4, 7, 8, 9, 20, 30, 40}, // 1
				{1, 2, 3, 4, 7},                   // 2
				{1, 2, 3, 4, 20, 30, 40},          // 3
				{1, 2, 3, 4, 5},                   // 4
				{1, 2, 3, 4, 5, 6, 20, 30, 40},    // 5
				{1, 2, 3, 4, 5},                   // 6
			},
		},
	}

	for i, tt := range tests {
		storage.addTagsToFiles(tt.files, tt.tags)
		for id, res := range tt.result {
			f, ok := storage.files[id+1]
			if !ok {
				t.Errorf("Test #%d: file with id %d doesn't exist", i+1, id+1)
				break
			}

			if !areArraysEqualInt(f.Tags, res) {
				t.Errorf("Test #%d: file with id %d: Want: %v Got: %v", i+1, id+1, res, f.Tags)
			}
		}
	}

	removeConfigFile()
}

func TestRemoveTagsFromFiles(t *testing.T) {
	storage := newStorage()
	storage.init()
	addDefaultFiles(storage)

	tests := []struct {
		files  []int
		tags   []int
		result [][]int // (index + 1) – id of a file
	}{
		{
			files: []int{1, 2, 3, 4, 5, 6},
			tags:  []int{1},
			result: [][]int{
				{2, 3},    // 1
				{2, 7},    // 2
				{3},       // 3
				{2, 3},    // 4
				{4, 5, 6}, // 5
				{},        // 6
			},
		},
		{
			files: []int{6},
			tags:  []int{1, 2, 3},
			result: [][]int{
				{2, 3},    // 1
				{2, 7},    // 2
				{3},       // 3
				{2, 3},    // 4
				{4, 5, 6}, // 5
				{},        // 6
			},
		},
		{
			files: []int{3, 4},
			tags:  []int{1, 2, 3},
			result: [][]int{
				{2, 3},    // 1
				{2, 7},    // 2
				{},        // 3
				{},        // 4
				{4, 5, 6}, // 5
				{},        // 6
			},
		},
		{
			files: []int{1, 2, 3, 4, 5, 6},
			tags:  []int{2, 3, 4, 5, 6, 7},
			result: [][]int{
				{}, // 1
				{}, // 2
				{}, // 3
				{}, // 4
				{}, // 5
				{}, // 6
			},
		},
	}

	for i, tt := range tests {
		storage.removeTagsFromFiles(tt.files, tt.tags)
		for id, res := range tt.result {
			f, ok := storage.files[id+1]
			if !ok {
				t.Errorf("Test #%d: file with id %d doesn't exist", i+1, id+1)
				break
			}

			if !areArraysEqualInt(f.Tags, res) {
				t.Errorf("Test #%d: file with id %d: Want: %v Got: %v", i+1, id+1, res, f.Tags)
			}
		}
	}

	removeConfigFile()
}

func TestGetFiles(t *testing.T) {
	storage := newStorage()
	storage.init()

	files := []struct {
		filename string
	}{
		{filename: "123.jpg"},
		{filename: "456.jpeg"},
		{filename: "1.png"},
		{filename: "test.exe"},
		{filename: "123.png"},
		{filename: "text.txt"},
	}

	now := time.Now()
	for _, f := range files {
		storage.addFile(f.filename, extensions.Ext{}, []int{}, 0, now)
	}

	requests := []struct {
		search   string
		isRegexp bool
		result   []string // filenames
	}{
		{search: "", isRegexp: false, result: []string{"123.jpg", "456.jpeg", "1.png", "test.exe", "123.png", "text.txt"}},
		{search: "123\\.jpg", isRegexp: false, result: []string{}},
		{search: "123", isRegexp: false, result: []string{"123.jpg", "123.png"}},
		//
		{search: "", isRegexp: true, result: []string{"123.jpg", "456.jpeg", "1.png", "test.exe", "123.png", "text.txt"}},
		{search: "123", isRegexp: true, result: []string{"123.jpg", "123.png"}},
		{search: "123\\.jpg", isRegexp: true, result: []string{"123.jpg"}},
		{search: "123\\.(jpg|png)", isRegexp: true, result: []string{"123.jpg", "123.png"}},
		{search: "123\\.(jpg|png|txt)", isRegexp: true, result: []string{"123.jpg", "123.png"}},
		{search: "^1.*g$", isRegexp: true, result: []string{"123.jpg", "123.png", "1.png"}},
		{search: "\\.[^png]+$", isRegexp: true, result: []string{"test.exe", "text.txt"}},
	}

	for i, r := range requests {
		files := storage.getFiles("", r.search, r.isRegexp)
		if len(files) != len(r.result) {
			t.Errorf("Test #%d \nWant: %v\nGet: %v\n\n", i, r.result, files)
		} else {
			a1 := make([]string, len(files))
			a2 := make([]string, len(files))

			for i := range files {
				a1[i] = r.result[i]
				a2[i] = files[i].Filename
			}

			if !areArraysEqualString(a1, a2) {
				t.Errorf("Test #%d \nWant: %v\nGet: %v\n\n", i, r.result, files)
			}
		}
	}

	removeConfigFile()
}
