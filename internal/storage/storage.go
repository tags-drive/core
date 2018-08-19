package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/ShoshinNikita/log"
	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

type File struct {
	Filename string
	Size     int64
	Tags     []string
}

var (
	filesMutex = new(sync.RWMutex)
	files      []File
)

// Init reads params.TagsFiles and decode its data
func Init() error {
	f, err := os.OpenFile(params.TagsFile, os.O_RDWR, 0600)
	if err != nil {
		// Have to create a new file
		if os.IsNotExist(err) {
			log.Infof("File %s doesn't exist. Need to create a new file\n", params.TagsFile)
			f, err = os.OpenFile(params.TagsFile, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return errors.Wrap(err, "can't create a new file")
			}
			// Write empty list
			f.Write([]byte("[]"))
			f.Close()
			// Can exit because we don't need to decode files from the file
			return nil
		}

		return errors.Wrapf(err, "can't open file %s", params.TagsFile)
	}

	defer f.Close()
	err = json.NewDecoder(f).Decode(&files)
	if err != nil {
		return errors.Wrapf(err, "can't decode data")
	}

	return nil
}
