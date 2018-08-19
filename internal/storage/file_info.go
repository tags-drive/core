package storage

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

type FileInfo struct {
	Filename string    `json:"filename"`
	Size     int64     `json:"size"`
	Tags     []string  `json:"tags"`
	AddTime  time.Time `json:"add_time"`
}

// Files is a map (filename: File) with RWMutex
// files.json keeps only Files.info
type Files struct {
	info  map[string]FileInfo
	mutex *sync.RWMutex
}

// write writes fs.info into params.TagsFile
func (fs Files) write() {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// TODO
	f, _ := os.OpenFile(params.TagsFile, os.O_RDWR, 0600)
	json.NewEncoder(f).Encode(fs.info)
	f.Close()
}

// decode decodes fs.info
func (fs *Files) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&fs.info)
}

// add adds an element into fs.info and call fs.write()
func (fs *Files) add(info FileInfo) error {
	fs.mutex.Lock()
	// TODO
	if _, ok := fs.info[info.Filename]; ok {

	}

	fs.info[info.Filename] = info
	fs.mutex.Unlock()

	fs.write()

	return nil
}

// delete deletes an element and call fs.write()
func (fs *Files) delete(filename string) {
	fs.write()
}
