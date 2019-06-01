package files

import (
	"crypto/sha256"
	"os"
	"testing"
	"time"

	clog "github.com/ShoshinNikita/log/v2"

	"github.com/stretchr/testify/assert"
	"github.com/tags-drive/core/internal/storage/files/extensions"
)

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
		clog.Errorln("can't remove folder \"data\":", err)
	}

	os.Exit(code)
}

func TestCheckFile(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

	addDefaultFiles(storage)

	tests := []struct {
		id  int
		res bool
	}{
		{1, true},
		{2, true},
		{50, false},
		{-1, false},
	}

	for i, tt := range tests {
		res := storage.checkFile(tt.id)
		assert.Equalf(tt.res, res, "iteration #%d", i+1)
	}
}

func TestAddFile(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

	now := time.Now()

	tests := []struct {
		filename string
		tags     []int
		size     int64
		ext      extensions.Ext
		time     time.Time
		//
		res map[int]File
	}{
		{
			filename: "1.jpg",
			tags:     []int{1, 2, 3},
			size:     15,
			time:     now,
			ext:      extensions.Ext{},
			res: map[int]File{
				1: {
					ID:       1,
					Filename: "1.jpg",
					Size:     15,
					Tags:     []int{1, 2, 3},
					AddTime:  now,
					Type:     extensions.Ext{},
					Origin:   storage.config.DataFolder + "/" + "1",
				},
			},
		},
		{
			filename: "35.jpg",
			tags:     []int{88},
			size:     345,
			time:     now,
			ext:      extensions.Ext{},
			res: map[int]File{
				1: {
					ID:       1,
					Filename: "1.jpg",
					Size:     15,
					Tags:     []int{1, 2, 3},
					AddTime:  now,
					Type:     extensions.Ext{},
					Origin:   storage.config.DataFolder + "/" + "1",
				},
				2: {
					ID:       2,
					Filename: "35.jpg",
					Size:     345,
					Tags:     []int{88},
					AddTime:  now,
					Type:     extensions.Ext{},
					Origin:   storage.config.DataFolder + "/" + "2",
				},
			},
		},
		{
			filename: "545.jpg",
			tags:     nil,
			size:     666,
			time:     now,
			ext: extensions.Ext{
				Ext:         ".jpg",
				FileType:    extensions.FileTypeImage,
				Supported:   true,
				PreviewType: extensions.PreviewTypeImage,
			},
			res: map[int]File{
				1: {
					ID:       1,
					Filename: "1.jpg",
					Size:     15,
					Tags:     []int{1, 2, 3},
					AddTime:  now,
					Type:     extensions.Ext{},
					Origin:   storage.config.DataFolder + "/" + "1",
				},
				2: {
					ID:       2,
					Filename: "35.jpg",
					Size:     345,
					Tags:     []int{88},
					AddTime:  now,
					Type:     extensions.Ext{},
					Origin:   storage.config.DataFolder + "/" + "2",
				},
				3: {
					ID:       3,
					Filename: "545.jpg",
					Size:     666,
					Tags:     []int{},
					AddTime:  now,
					Type: extensions.Ext{
						Ext:         ".jpg",
						FileType:    extensions.FileTypeImage,
						Supported:   true,
						PreviewType: extensions.PreviewTypeImage,
					},
					Origin:  storage.config.DataFolder + "/" + "3",
					Preview: storage.config.ResizedImagesFolder + "/" + "3",
				},
			},
		},
	}

	for i, tt := range tests {
		storage.addFile(tt.filename, tt.ext, tt.tags, tt.size, now)

		assert.Equalf(tt.res, storage.files, "iteration #%d", i+1)
	}
}

func TestGetFile(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

	addDefaultFiles(storage)

	tests := []struct {
		id      int
		res     string // filename
		isError bool
	}{
		{1, "1", false},
		{2, "2", false},
		{50, "", true},
		{-1, "", true},
	}

	for i, tt := range tests {
		f, err := storage.getFile(tt.id)
		if !assert.Equalf(tt.isError, err != nil, "iteration #%d, error: %v", i+1, err) {
			continue
		}

		assert.Equalf(tt.res, f.Filename, "iteration #%d", i+1)
	}
}

func TestGetFilesWithIDs(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

	addDefaultFiles(storage)

	tests := []struct {
		ids []int
		res []int // files ids
	}{
		{
			ids: []int{1, 2, 3},
			res: []int{1, 2, 3},
		},
		{
			ids: []int{1, 2, 10, 20, 33},
			res: []int{1, 2},
		},
		{
			ids: []int{10, 20, 30},
			res: []int{},
		},
	}

	for i, tt := range tests {
		files := storage.getFilesWithIDs(tt.ids...)
		res := []int{}
		for _, f := range files {
			res = append(res, f.ID)
		}

		assert.Equal(tt.res, res, "iteration #%d", i+1)
	}
}

func TestGetFiles(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

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
		if !assert.Equalf(len(r.result), len(files), "iteration #%d", i+1) {
			continue
		}

		filenames := make([]string, len(files))
		for i := range files {
			filenames[i] = files[i].Filename
		}

		assert.ElementsMatchf(r.result, filenames, "iteration #%d", i+1)
	}
}

func TestRenameFile(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()
	addDefaultFiles(storage)

	tests := []struct {
		id          int
		newFilename string
		res         map[int]string // id - filename
		isError     bool
	}{
		{1, "123", map[int]string{1: "123", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6"}, false},
		{2, "123", map[int]string{1: "123", 2: "123", 3: "3", 4: "4", 5: "5", 6: "6"}, false},
		{6, "654", map[int]string{1: "123", 2: "123", 3: "3", 4: "4", 5: "5", 6: "654"}, false},
		{88, "000", map[int]string{1: "123", 2: "123", 3: "3", 4: "4", 5: "5", 6: "654"}, true},
		{5, "", map[int]string{1: "123", 2: "123", 3: "3", 4: "4", 5: "5", 6: "654"}, true},
	}

	for i, tt := range tests {
		_, err := storage.renameFile(tt.id, tt.newFilename)
		if !assert.Equalf(tt.isError, err != nil, "iteration #%d, error: %v", i+1, err) {
			continue
		}

		res := make(map[int]string)
		for id, f := range storage.files {
			res[id] = f.Filename
		}

		assert.Equalf(tt.res, res, "iteration #%d", i+1)
	}
}

func TestUpdateTags(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()
	addDefaultFiles(storage)

	tests := []struct {
		id      int
		newTags []int
		res     map[int][]int // id - tags
		isError bool
	}{
		{1, []int{1, 5}, map[int][]int{
			1: {1, 5},
			2: {1, 2, 7},
			3: {3},
			4: {2, 3},
			5: {4, 5, 6},
			6: {},
		}, false},
		{2, []int{}, map[int][]int{
			1: {1, 5},
			2: {},
			3: {3},
			4: {2, 3},
			5: {4, 5, 6},
			6: {},
		}, false},
		{3, nil, map[int][]int{
			1: {1, 5},
			2: {},
			3: {},
			4: {2, 3},
			5: {4, 5, 6},
			6: {},
		}, false},
		{88, nil, map[int][]int{
			1: {1, 5},
			2: {},
			3: {},
			4: {2, 3},
			5: {4, 5, 6},
			6: {},
		}, true},
	}

	for i, tt := range tests {
		_, err := storage.updateFileTags(tt.id, tt.newTags)
		if !assert.Equalf(tt.isError, err != nil, "iteration #%d, error: %v", i+1, err) {
			continue
		}

		res := make(map[int][]int)
		for id, f := range storage.files {
			res[id] = f.Tags
		}

		assert.Equalf(tt.res, res, "iteration #%d", i+1)
	}
}

func TestUpdateDescription(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()
	addDefaultFiles(storage)

	tests := []struct {
		id      int
		newDesc string
		res     map[int]string // id - description
		isError bool
	}{
		{1, "help", map[int]string{1: "help", 2: "", 3: "", 4: "", 5: "", 6: ""}, false},
		{5, "help - 5", map[int]string{1: "help", 2: "", 3: "", 4: "", 5: "help - 5", 6: ""}, false},
		{5, "АБВГД", map[int]string{1: "help", 2: "", 3: "", 4: "", 5: "АБВГД", 6: ""}, false},
		{88, "123", map[int]string{1: "help", 2: "", 3: "", 4: "", 5: "АБВГД", 6: ""}, true},
	}

	for i, tt := range tests {
		_, err := storage.updateFileDescription(tt.id, tt.newDesc)
		if !assert.Equalf(tt.isError, err != nil, "iteration #%d, error: %v", i+1, err) {
			continue
		}

		res := make(map[int]string)
		for id, f := range storage.files {
			res[id] = f.Description
		}

		assert.Equalf(tt.res, res, "iteration #%d", i+1)
	}
}

func TestDeleteFileForce(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()
	addDefaultFiles(storage)

	tests := []struct {
		id  int
		res []int // ids
	}{
		{1, []int{2, 3, 4, 5, 6}},
		{6, []int{2, 3, 4, 5}},
		{4, []int{2, 3, 5}},
		{1, []int{2, 3, 5}},
	}

	for i, tt := range tests {
		storage.deleteFileForce(tt.id)

		var ids []int
		for id := range storage.files {
			ids = append(ids, id)
		}

		assert.ElementsMatchf(tt.res, ids, "iteration #%d", i+1)
	}
}

func TestDeleteAndRecover(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()
	addDefaultFiles(storage)

	tests := []struct {
		id            int
		shouldRecover bool
		//
		res map[int]bool // id - deleted
	}{
		{
			id:            1,
			shouldRecover: false,
			res: map[int]bool{
				1: true,
				2: false,
				3: false,
				4: false,
				5: false,
				6: false,
			},
		},
		{
			id:            1,
			shouldRecover: false,
			res: map[int]bool{
				1: true,
				2: false,
				3: false,
				4: false,
				5: false,
				6: false,
			},
		},
		{
			id:            5,
			shouldRecover: false,
			res: map[int]bool{
				1: true,
				2: false,
				3: false,
				4: false,
				5: true,
				6: false,
			},
		},
		{
			id:            1,
			shouldRecover: true,
			res: map[int]bool{
				1: false,
				2: false,
				3: false,
				4: false,
				5: true,
				6: false,
			},
		},
		{
			id:            1,
			shouldRecover: true,
			res: map[int]bool{
				1: false,
				2: false,
				3: false,
				4: false,
				5: true,
				6: false,
			},
		},
		// non-existing files
		{
			id:            150,
			shouldRecover: false,
			res: map[int]bool{
				1: false,
				2: false,
				3: false,
				4: false,
				5: true,
				6: false,
			},
		},
		{
			id:            150,
			shouldRecover: true,
			res: map[int]bool{
				1: false,
				2: false,
				3: false,
				4: false,
				5: true,
				6: false,
			},
		},
	}

	for i, tt := range tests {
		if tt.shouldRecover {
			storage.recover(tt.id)
		} else {
			storage.deleteFile(tt.id)
		}

		res := make(map[int]bool)
		for id, f := range storage.files {
			res[id] = f.Deleted
		}

		assert.Equalf(tt.res, res, "iteration #%d", i+1)
	}
}

func TestAddTagsToFiles(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

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
			ok := storage.checkFile(id + 1)
			if !assert.Equalf(true, ok, "iteration #%d: file with id %d doesn't exist", i+1, id+1) {
				break
			}

			f := storage.files[id+1]
			assert.ElementsMatchf(res, f.Tags, "iteration #%d: file with id %d has wrong tags", i+1, id+1)
		}
	}
}

func TestRemoveTagsFromFiles(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

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
			ok := storage.checkFile(id + 1)
			if !assert.Equalf(true, ok, "iteration #%d: file with id %d doesn't exist", i+1, id+1) {
				break
			}

			f := storage.files[id+1]
			assert.ElementsMatchf(res, f.Tags, "iteration #%d: file with id %d has wrong tags", i+1, id+1)
		}
	}
}

func TestDeleteTagFromFiles(t *testing.T) {
	assert := assert.New(t)

	storage := newStorage()
	defer func() {
		storage.shutdown()
		os.Remove(storage.config.FilesJSONFile)
	}()

	addDefaultFiles(storage)

	tests := []struct {
		idToDelete int
		res        map[int][]int // id - tags
	}{
		{1, map[int][]int{
			1: {2, 3},
			2: {2, 7},
			3: {3},
			4: {2, 3},
			5: {4, 5, 6},
			6: {},
		}},
		{5, map[int][]int{
			1: {2, 3},
			2: {2, 7},
			3: {3},
			4: {2, 3},
			5: {4, 6},
			6: {},
		}},
		{20, map[int][]int{
			1: {2, 3},
			2: {2, 7},
			3: {3},
			4: {2, 3},
			5: {4, 6},
			6: {},
		}},
	}

	for i, tt := range tests {
		storage.deleteTagFromFiles(tt.idToDelete)

		res := make(map[int][]int)
		for _, f := range storage.files {
			res[f.ID] = f.Tags
		}

		assert.Equalf(tt.res, res, "iteration #%d", i+1)
	}
}

// newStorage creates new jsonFileStorage and call init() function
func newStorage() *jsonFileStorage {
	cnf := Config{
		Debug:               false,
		DataFolder:          "./data",
		ResizedImagesFolder: "./data/resizing",
		StorageType:         "json",
		FilesJSONFile:       "files.json",
		Encrypt:             true,
		PassPhrase:          sha256.Sum256([]byte("sha256")),
	}
	st := newJsonFileStorage(cnf, clog.NewProdLogger())
	st.init()

	return st
}

// addDefaultFiles adds next files into storage:
//  | id | name | tags      |
//  | -- | ---- | --------- |
//  | 1  | 1    | [1, 2, 3] |
//  | 2  | 2    | [1, 2, 7] |
//  | 3  | 3    | [3]       |
//  | 4  | 4    | [2, 3]    |
//  | 5  | 5    | [4, 5, 6] |
//  | 6  | 6    | []        |
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
