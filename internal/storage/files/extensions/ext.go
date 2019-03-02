package extensions

import (
	"errors"
	"sync"

	"github.com/tags-drive/core/cmd"
)

var (
	// UnsupportedExt is used for unsupported file extension
	UnsupportedExt    = cmd.Ext{Supported: false, FileType: cmd.FileTypeUnsupported}
	errUnsupportedExt = errors.New("unsupported extension")
)

type extensions struct {
	exts map[string]cmd.Ext
	mut  *sync.RWMutex
}

func (e *extensions) add(ext cmd.Ext) {
	e.mut.Lock()
	defer e.mut.Unlock()

	e.exts[ext.Ext] = ext
}

func (e *extensions) get(ext string) (cmd.Ext, error) {
	e.mut.RLock()
	defer e.mut.RUnlock()

	res, ok := e.exts[ext]
	if !ok {
		return cmd.Ext{}, errUnsupportedExt
	}

	return res, nil
}

var allExtensions extensions

func init() {
	allExtensions = extensions{
		exts: make(map[string]cmd.Ext),
		mut:  new(sync.RWMutex),
	}

	for i := range extensionsList {
		allExtensions.add(extensionsList[i])
	}
}

// GetExt returns Ext according to passed file ext.
// If there's no such extension, it returns cmd.UnsupportedExt
func GetExt(ext string) cmd.Ext {
	res, err := allExtensions.get(ext)
	if err != nil {
		res = UnsupportedExt
		res.Ext = ext
	}

	return res
}
