package tags

import (
	"os"
	"testing"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/stretchr/testify/assert"
)

const testFile = "./tags.json"

func newReadyStorage() (*jsonTagStorage, error) {
	cnf := Config{
		Debug:        false,
		StorageType:  "json",
		TagsJSONFile: testFile,
		Encrypt:      false,
		// PassPhrase:   sha256.Sum256([]byte("sha256")),
	}

	st := newJsonTagStorage(cnf, clog.NewProdLogger())
	err := st.init()

	return st, err
}

func TestInit(t *testing.T) {
	assert := assert.New(t)
	cnf := Config{
		Debug:        false,
		StorageType:  "json",
		TagsJSONFile: testFile,
		Encrypt:      false,
		// PassPhrase:   sha256.Sum256([]byte("sha256")),
	}

	testStorage := newJsonTagStorage(cnf, clog.NewProdLogger())

	err := testStorage.init()
	if !assert.Nil(err, "can't init test storage") {
		t.FailNow()
	}

	err = testStorage.shutdown()
	assert.Nil(err, "can't init test storage")

	os.Remove(testFile)
}

func TestAddAndDelete(t *testing.T) {
	type testType int
	const (
		add testType = iota
		delete
	)

	assert := assert.New(t)

	storage, err := newReadyStorage()
	if !assert.Nil(err, "can't init storage") {
		t.FailNow()
	}

	tests := []struct {
		testType     testType
		tagsToAdd    []Tag // if testType == add
		tagsToDelete []int // if testType == delete
		result       Tags
	}{
		{
			testType: add,
			tagsToAdd: []Tag{
				{Name: "test1", Color: "#fffff0", Group: "test"},
			},
			result: Tags{
				1: Tag{ID: 1, Name: "test1", Color: "#fffff0", Group: "test"},
			},
		},
		{
			testType: add,
			tagsToAdd: []Tag{
				{Name: "test2", Color: "#ffff0f"},
				{Name: "test3", Color: "#fff0ff", Group: "123"},
			},
			result: Tags{
				1: Tag{ID: 1, Name: "test1", Color: "#fffff0", Group: "test"},
				2: Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
				3: Tag{ID: 3, Name: "test3", Color: "#fff0ff", Group: "123"},
			},
		},
		{
			testType: add,
			tagsToAdd: []Tag{
				{Name: "test4", Color: "#ff0fff"},
				{Name: "test5", Color: "#f0ffff"},
				{Name: "test6", Color: "#0fffff"},
			},
			result: Tags{
				1: Tag{ID: 1, Name: "test1", Color: "#fffff0", Group: "test"},
				2: Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
				3: Tag{ID: 3, Name: "test3", Color: "#fff0ff", Group: "123"},
				4: Tag{ID: 4, Name: "test4", Color: "#ff0fff"},
				5: Tag{ID: 5, Name: "test5", Color: "#f0ffff"},
				6: Tag{ID: 6, Name: "test6", Color: "#0fffff"},
			},
		},
		{
			testType: add,
			tagsToAdd: []Tag{
				{Name: "test6", Color: "#111111"},
			},
			result: Tags{
				1: Tag{ID: 1, Name: "test1", Color: "#fffff0", Group: "test"},
				2: Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
				3: Tag{ID: 3, Name: "test3", Color: "#fff0ff", Group: "123"},
				4: Tag{ID: 4, Name: "test4", Color: "#ff0fff"},
				5: Tag{ID: 5, Name: "test5", Color: "#f0ffff"},
				6: Tag{ID: 6, Name: "test6", Color: "#0fffff"},
				7: Tag{ID: 7, Name: "test6", Color: "#111111"},
			},
		},
		{
			testType: delete,
			tagsToDelete: []int{
				1,
				3,
				5,
			},
			result: Tags{
				2: Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
				4: Tag{ID: 4, Name: "test4", Color: "#ff0fff"},
				6: Tag{ID: 6, Name: "test6", Color: "#0fffff"},
				7: Tag{ID: 7, Name: "test6", Color: "#111111"},
			},
		},
		{
			testType: add,
			tagsToAdd: []Tag{
				{Name: "new tag", Color: "#111111"},
				{Name: "adsbcv", Color: "#222222"},
				{Name: "scvxcv", Color: "#111111"},
			},
			result: Tags{
				2:  Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
				4:  Tag{ID: 4, Name: "test4", Color: "#ff0fff"},
				6:  Tag{ID: 6, Name: "test6", Color: "#0fffff"},
				7:  Tag{ID: 7, Name: "test6", Color: "#111111"},
				8:  Tag{ID: 8, Name: "new tag", Color: "#111111"},
				9:  Tag{ID: 9, Name: "adsbcv", Color: "#222222"},
				10: Tag{ID: 10, Name: "scvxcv", Color: "#111111"},
			},
		},
		{
			testType: delete,
			tagsToDelete: []int{
				2,
				4,
				6,
				7,
				9,
			},
			result: Tags{

				8:  Tag{ID: 8, Name: "new tag", Color: "#111111"},
				10: Tag{ID: 10, Name: "scvxcv", Color: "#111111"},
			},
		},
	}

	for i, tt := range tests {
		if tt.testType == add {
			for _, tag := range tt.tagsToAdd {
				storage.addTag(tag)
			}
		} else if tt.testType == delete {
			for _, id := range tt.tagsToDelete {
				storage.deleteTag(id)
			}
		}

		res := storage.getAll()
		assert.Equalf(tt.result, res, "iteration #%d", i)
	}

	storage.shutdown()
	os.Remove(testFile)
}

func TestUpdate(t *testing.T) {
	type testType int
	const (
		updateTag testType = iota
		updateGroup
	)

	assert := assert.New(t)

	storage, err := newReadyStorage()
	if !assert.Nil(err, "can't init storage") {
		t.FailNow()
	}

	startTags := []Tag{
		{Name: "test1", Color: "#fffff0"},
		{Name: "test2", Color: "#ffff0f"},
		{Name: "test3", Color: "#fff0ff"},
		{Name: "test4", Color: "#ff0fff"},
		{Name: "test5", Color: "#f0ffff"},
		{Name: "test6", Color: "#0fffff"},
	}

	tests := []struct {
		testType testType
		id       int
		newName  string
		newColor string
		newGroup string
		result   Tag
	}{
		// No changes
		{
			testType: updateTag,
			id:       1,
			newName:  "",
			newColor: "",
			result:   Tag{ID: 1, Name: "test1", Color: "#fffff0"},
		},
		// Change name
		{
			testType: updateTag,
			id:       5,
			newName:  "hello",
			newColor: "",
			result:   Tag{ID: 5, Name: "hello", Color: "#f0ffff"},
		},
		// Change color (without #)
		{
			testType: updateTag,
			id:       4,
			newName:  "",
			newColor: "ff0000",
			result:   Tag{ID: 4, Name: "test4", Color: "#ff0000"},
		},
		// Change name and color
		{
			testType: updateTag,
			id:       2,
			newName:  "123",
			newColor: "#efefef",
			result:   Tag{ID: 2, Name: "123", Color: "#efefef"},
		},
		// Change group
		{
			testType: updateGroup,
			id:       2,
			newGroup: "test",
			result:   Tag{ID: 2, Name: "123", Color: "#efefef", Group: "test"},
		},
		// Change group (to "")
		{
			testType: updateGroup,
			id:       2,
			newGroup: "",
			result:   Tag{ID: 2, Name: "123", Color: "#efefef", Group: ""},
		},
	}

	// Add start tags
	for _, t := range startTags {
		storage.addTag(t)
	}

	var res Tag
	for i, tt := range tests {
		if tt.testType == updateTag {
			res, err = storage.updateTag(tt.id, tt.newName, tt.newColor)
		} else if tt.testType == updateGroup {
			res, err = storage.updateGroup(tt.id, tt.newGroup)
		}

		if !assert.Nilf(err, "iteration %d: got an error: %s", i, err) {
			continue
		}
		assert.Equalf(tt.result, res, "iteration %d", i)
	}

	storage.shutdown()
	os.Remove(testFile)
}
