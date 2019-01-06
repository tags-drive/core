package tags

import (
	"bytes"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/ShoshinNikita/log"
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

func newStorage() *jsonTagStorage {
	return &jsonTagStorage{
		tags:   make(Tags),
		mutex:  new(sync.RWMutex),
		logger: log.NewLogger(),
	}
}

func TestMain(m *testing.M) {
	// All tests are called sequentially. Every test creates new instance of jsonTagStorage.
	// So we hadn't to remove tags.json file because of it is trunced in jsonTagStorage.write() func.

	// Create folder storage/tags/configs
	err := os.Mkdir("configs", 0666)
	if err != nil && !os.IsExist(err) {
		log.Fatalln(err)
		return
	}

	// We will create tags.json in TestInit function
	code := m.Run()

	// Remove test file
	err = os.RemoveAll("configs")
	if err != nil {
		log.Fatalln(err)
		return
	}

	os.Exit(code)
}

func TestInit(t *testing.T) {
	testStorage := newStorage()

	err := testStorage.init()
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("configs/tags.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, []byte("{}")) {
		t.Errorf("Wrong file content: %s", string(data))
	}

	// Write new tag. It must be saved in file
	testStorage.addTag(Tag{Name: "first", Color: "#ffffff"})
}

func TestInit2(t *testing.T) {
	testStorage := newStorage()

	err := testStorage.init()
	if err != nil {
		t.Fatal(err)
	}
	tags := testStorage.getAll()

	tagInFile := Tag{ID: 1, Name: "first", Color: "#ffffff"}

	if len(tags) != 1 || !testStorage.check(1) || tags[1] != tagInFile {
		t.Errorf("wrong file content: len(tags): %d, allTags: %v", len(tags), tags)
	}
}

func TestAdd(t *testing.T) {
	testStorage := newStorage()

	tags := []Tag{
		Tag{Name: "test1", Color: "#fffff0"},
		Tag{Name: "test2", Color: "#ffff0f"},
		Tag{Name: "test3", Color: "#fff0ff"},
		Tag{Name: "test4", Color: "#ff0fff"},
		Tag{Name: "test5", Color: "#f0ffff"},
		Tag{Name: "test6", Color: "#0fffff"},
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
}

func TestDelete(t *testing.T) {
	testStorage := newStorage()

	// Default tags
	startTags := []Tag{
		Tag{Name: "test1", Color: "#fffff0"},
		Tag{Name: "test2", Color: "#ffff0f"},
		Tag{Name: "test3", Color: "#fff0ff"},
		Tag{Name: "test4", Color: "#ff0fff"},
		Tag{Name: "test5", Color: "#f0ffff"},
		Tag{Name: "test6", Color: "#0fffff"},
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
		Tag{Name: "123", Color: "#ff0000"},
		Tag{Name: "456", Color: "#00ff00"},
		Tag{Name: "789", Color: "#0000ff"},
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
}

func TestUpdate(t *testing.T) {
	startTags := []Tag{
		Tag{Name: "test1", Color: "#fffff0"},
		Tag{Name: "test2", Color: "#ffff0f"},
		Tag{Name: "test3", Color: "#fff0ff"},
		Tag{Name: "test4", Color: "#ff0fff"},
		Tag{Name: "test5", Color: "#f0ffff"},
		Tag{Name: "test6", Color: "#0fffff"},
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
}
