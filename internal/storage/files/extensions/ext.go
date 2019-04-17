package extensions

import (
	"errors"
	"strings"
	"sync"
)

var (
	// UnsupportedExt is used for unsupported file extension
	UnsupportedExt    = Ext{Supported: false, FileType: FileTypeUnsupported}
	errUnsupportedExt = errors.New("unsupported extension")
)

type extensions struct {
	exts map[string]Ext
	mut  *sync.RWMutex
}

func (e *extensions) add(ext Ext) {
	e.mut.Lock()
	defer e.mut.Unlock()

	e.exts[ext.Ext] = ext
}

func (e *extensions) get(ext string) (Ext, error) {
	e.mut.RLock()
	defer e.mut.RUnlock()

	res, ok := e.exts[ext]
	if !ok {
		return Ext{}, errUnsupportedExt
	}

	return res, nil
}

var allExtensions extensions

func init() {
	allExtensions = extensions{
		exts: make(map[string]Ext),
		mut:  new(sync.RWMutex),
	}

	for i := range extensionsList {
		allExtensions.add(extensionsList[i])
	}
}

// GetExt returns Ext according to passed file ext.
// If there's no such extension, it returns UnsupportedExt
func GetExt(ext string) Ext {
	if len(ext) == 0 {
		return UnsupportedExt
	}

	ext = strings.ToLower(ext)
	if ext[0] != '.' {
		ext = "." + ext
	}

	res, err := allExtensions.get(ext)
	if err != nil {
		res = UnsupportedExt
		res.Ext = ext
	}

	return res
}
