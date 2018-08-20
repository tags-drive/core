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

// filesData is a map (filename: FileInfo) with RWMutex
// files.json keeps only filesData.info
type filesData struct {
	info  map[string]FileInfo
	mutex *sync.RWMutex
}

// write writes fs.info into params.TagsFile
func (fs filesData) write() {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// TODO
	f, _ := os.OpenFile(params.TagsFile, os.O_RDWR, 0600)
	// Write pretty json if Debug mode
	enc := json.NewEncoder(f)
	if params.Debug {
		enc.SetIndent("", "  ")
	}
	enc.Encode(fs.info)

	f.Close()
}

// decode decodes fs.info
func (fs *filesData) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&fs.info)
}

// add adds an element into fs.info and call fs.write()
func (fs *filesData) add(info FileInfo) error {
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
func (fs *filesData) delete(filename string) {
	fs.write()
}
