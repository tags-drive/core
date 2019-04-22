package tags

import (
	"bytes"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	clog "github.com/ShoshinNikita/log/v2"
)

func areTagsEqual(a, b Tags) bool {
	if len(a) != len(b) {
		return false
	}

	for k, t := range a {
		if tt, ok := b[k]; !ok || t != tt {
			return false
		}
	}

	return true
}

func removeConfigFile(path string) {
	os.Remove(path)
}

func newStorage() *jsonTagStorage {
	cnf := Config{
		Debug:        false,
		StorageType:  "json",
		TagsJSONFile: "tags.json",
		Encrypt:      false,
		// PassPhrase:   sha256.Sum256([]byte("sha256")),
	}
	return newJsonTagStorage(cnf, clog.NewProdLogger())
}

func TestInit(t *testing.T) {
	testStorage := newStorage()
	testStorage.init()

	err := testStorage.init()
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(testStorage.config.TagsJSONFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if !(bytes.Equal(data, []byte("{}")) || bytes.Equal(data, []byte("{}\n"))) {
		t.Errorf("Wrong file content: %s", string(data))
	}

	removeConfigFile(testStorage.config.TagsJSONFile)
}

func TestAdd(t *testing.T) {
	testStorage := newStorage()
	testStorage.init()

	tags := []Tag{
		{Name: "test1", Color: "#fffff0"},
		{Name: "test2", Color: "#ffff0f"},
		{Name: "test3", Color: "#fff0ff"},
		{Name: "test4", Color: "#ff0fff"},
		{Name: "test5", Color: "#f0ffff"},
		{Name: "test6", Color: "#0fffff"},
	}

	answer := Tags{
		1: Tag{ID: 1, Name: "test1", Color: "#fffff0"},
		2: Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
		3: Tag{ID: 3, Name: "test3", Color: "#fff0ff"},
		4: Tag{ID: 4, Name: "test4", Color: "#ff0fff"},
		5: Tag{ID: 5, Name: "test5", Color: "#f0ffff"},
		6: Tag{ID: 6, Name: "test6", Color: "#0fffff"},
	}

	for _, tag := range tags {
		testStorage.addTag(tag)
	}

	result := testStorage.getAll()

	if !areTagsEqual(result, answer) {
		t.Errorf("Want: %v\n\nGot: %v", answer, result)
	}

	removeConfigFile(testStorage.config.TagsJSONFile)
}

func TestDelete(t *testing.T) {
	testStorage := newStorage()
	testStorage.init()

	// Default tags
	startTags := []Tag{
		{Name: "test1", Color: "#fffff0"},
		{Name: "test2", Color: "#ffff0f"},
		{Name: "test3", Color: "#fff0ff"},
		{Name: "test4", Color: "#ff0fff"},
		{Name: "test5", Color: "#f0ffff"},
		{Name: "test6", Color: "#0fffff"},
	}
	for _, tag := range startTags {
		testStorage.addTag(tag)
	}

	idsToDelete := []int{1, 3, 5, 10}
	for _, id := range idsToDelete {
		testStorage.deleteTag(id)
	}

	// Check
	answer := Tags{
		2: Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
		4: Tag{ID: 4, Name: "test4", Color: "#ff0fff"},
		6: Tag{ID: 6, Name: "test6", Color: "#0fffff"},
	}
	result := testStorage.getAll()
	if !areTagsEqual(result, answer) {
		t.Errorf("Want: %v\n\nGot: %v", answer, result)
	}

	// Add new Tags
	newTags := []Tag{
		{Name: "123", Color: "#ff0000"},
		{Name: "456", Color: "#00ff00"},
		{Name: "789", Color: "#0000ff"},
	}
	for _, tag := range newTags {
		testStorage.addTag(tag)
	}

	answer = Tags{
		2: Tag{ID: 2, Name: "test2", Color: "#ffff0f"},
		4: Tag{ID: 4, Name: "test4", Color: "#ff0fff"},
		6: Tag{ID: 6, Name: "test6", Color: "#0fffff"},
		7: Tag{ID: 7, Name: "123", Color: "#ff0000"},
		8: Tag{ID: 8, Name: "456", Color: "#00ff00"},
		9: Tag{ID: 9, Name: "789", Color: "#0000ff"},
	}
	result = testStorage.getAll()
	if !areTagsEqual(result, answer) {
		t.Errorf("Want: %v\n\nGot: %v", answer, result)
	}

	removeConfigFile(testStorage.config.TagsJSONFile)
}

func TestUpdate(t *testing.T) {
	startTags := []Tag{
		{Name: "test1", Color: "#fffff0"},
		{Name: "test2", Color: "#ffff0f"},
		{Name: "test3", Color: "#fff0ff"},
		{Name: "test4", Color: "#ff0fff"},
		{Name: "test5", Color: "#f0ffff"},
		{Name: "test6", Color: "#0fffff"},
	}

	type toUpdate struct {
		id    int
		name  string
		color string
	}

	tests := []struct {
		update toUpdate
		answer Tag
		ok     bool
	}{
		// No changes
		{
			toUpdate{id: 1, name: "", color: ""},
			Tag{ID: 1, Name: "test1", Color: "#fffff0"},
			true,
		},
		// Change name
		{
			toUpdate{id: 5, name: "hello", color: ""},
			Tag{ID: 5, Name: "hello", Color: "#f0ffff"},
			true,
		},
		// Change color (without #)
		{
			toUpdate{id: 4, name: "", color: "ff0000"},
			Tag{ID: 4, Name: "test4", Color: "#ff0000"},
			true,
		},
		// Change name and color
		{
			toUpdate{id: 2, name: "123", color: "#efefef"},
			Tag{ID: 2, Name: "123", Color: "#efefef"},
			true,
		},
		// Wrong id
		{
			toUpdate{id: 89, name: "123", color: "#efefef"},
			Tag{},
			false,
		},
	}

	testStorage := newStorage()
	testStorage.init()

	wg := new(sync.WaitGroup)
	for _, tag := range startTags {
		testStorage.addTag(tag)
	}

	for i, tt := range tests {
		wg.Add(1)

		go func(testID int, up toUpdate, ansTag Tag, ansOk bool) {
			defer wg.Done()

			testStorage.updateTag(up.id, up.name, up.color)
			result := testStorage.getAll()
			tag, ok := result[ansTag.ID]
			if ok != ansOk {
				t.Errorf("Test #%d wrong ok. Want: %v Got: %v", testID, ansOk, ok)
				return
			}

			if tag != ansTag {
				t.Errorf("Test #%d. Want: %v Got: %v", testID, ansTag, tag)
			}
		}(i, tt.update, tt.answer, tt.ok)
	}

	wg.Wait()

	removeConfigFile(testStorage.config.TagsJSONFile)
}
