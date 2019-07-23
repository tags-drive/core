package share

import (
	"os"
	"sort"
	"sync"

	"github.com/ShoshinNikita/log/v2"
	"github.com/pkg/errors"

	filesPck "github.com/tags-drive/core/internal/storage/files"
	tagsPck "github.com/tags-drive/core/internal/storage/tags"
	"github.com/tags-drive/core/internal/utils"
)

const maxTokenSize = 20

// fileIDs is a slice of sorted files ids
type filesIDs []int

func newFileIDs(ids []int) filesIDs {
	newIDs := append(ids[:0:0], ids...)

	sort.Ints(newIDs)

	return filesIDs(newIDs)
}

func (ids filesIDs) hasID(id int) bool {
	// ids is sorted, so we can use sort.SearchInts
	i := sort.SearchInts(ids, id)
	return i < len(ids) && ids[i] == id
}

func (ids *filesIDs) deleteID(id int) {
	// ids is sorted, so we can use sort.SearchInts
	i := sort.SearchInts(*ids, id)
	if i < len(*ids) && (*ids)[i] == id {
		*ids = append((*ids)[:i:i], (*ids)[i+1:]...)
	}
}

type jsonShareStorage struct {
	config Config

	tokens map[string]filesIDs
	mu     sync.RWMutex

	fileStorage FileStorage
	logger      *clog.Logger
}

func newJsonShareStorage(cnf Config, fileStorage FileStorage, lg *clog.Logger) *jsonShareStorage {
	return &jsonShareStorage{
		config:      cnf,
		tokens:      make(map[string]filesIDs),
		fileStorage: fileStorage,
		logger:      lg,
	}
}

func (jss *jsonShareStorage) init() error {
	if f, err := os.Open(jss.config.ShareTokenJSONFile); err == nil {
		err = utils.Decode(f, &jss.tokens, jss.config.Encrypt, jss.config.PassPhrase)
		f.Close()
		return err
	}

	// Have to create a new file
	jss.logger.Debugf("file %s doesn't exist. Need to create a new file\n", jss.config.ShareTokenJSONFile)

	// Just create a new file
	f, err := os.OpenFile(jss.config.ShareTokenJSONFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create a new file")
	}
	f.Close()

	// Write empty tag map
	jss.write()

	return nil

}

func (jss *jsonShareStorage) write() {
	jss.mu.RLock()
	defer jss.mu.RUnlock()

	f, err := os.OpenFile(jss.config.ShareTokenJSONFile, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		jss.logger.Errorf("can't open file %s: %s\n", jss.config.ShareTokenJSONFile, err)
		return
	}
	defer f.Close()

	err = utils.Encode(f, jss.tokens, jss.config.Encrypt, jss.config.PassPhrase)
	if err != nil {
		jss.logger.Warnf("can't write '%s': %s", jss.config.ShareTokenJSONFile, err)
	}
}

func (jss *jsonShareStorage) getAllTokens() map[string][]int {
	jss.mu.RLock()
	defer jss.mu.RUnlock()

	// Clone map
	res := make(map[string][]int)
	for token, ids := range jss.tokens {
		res[token] = append(ids[:0:0], ids...)
	}

	return res
}

func (jss *jsonShareStorage) createToken(ids []int) (token string) {
	jss.mu.Lock()
	defer func() {
		jss.mu.Unlock()
		jss.write()
	}()

	token = utils.GenerateRandomString(maxTokenSize)
	jss.tokens[token] = newFileIDs(ids)

	return token
}

func (jss *jsonShareStorage) deleteToken(token string) {
	jss.mu.Lock()
	defer func() {
		jss.mu.Unlock()
		jss.write()
	}()

	delete(jss.tokens, token)
}

func (jss *jsonShareStorage) getFilesIDs(token string) ([]int, error) {
	jss.mu.RLock()
	defer jss.mu.RUnlock()

	if _, ok := jss.tokens[token]; !ok {
		return nil, ErrInvalidToken
	}

	return jss.tokens[token], nil
}

func (jss *jsonShareStorage) checkToken(token string) bool {
	jss.mu.RLock()
	defer jss.mu.RUnlock()

	_, ok := jss.tokens[token]
	return ok
}

func (jss *jsonShareStorage) checkFile(token string, id int) bool {
	jss.mu.RLock()
	defer jss.mu.RUnlock()

	ids, ok := jss.tokens[token]
	if !ok {
		return false
	}

	return ids.hasID(id)
}

func (jss *jsonShareStorage) deleteFile(id int) {
	jss.mu.Lock()
	defer func() {
		jss.mu.Unlock()
		jss.write()
	}()

	for token, ids := range jss.tokens {
		ids.deleteID(id)
		jss.tokens[token] = ids
	}
}

func (jss *jsonShareStorage) filterFiles(token string, files []filesPck.File) ([]filesPck.File, error) {
	jss.mu.RLock()

	if _, ok := jss.tokens[token]; !ok {
		jss.mu.RUnlock()
		return nil, ErrInvalidToken
	}

	ids := jss.tokens[token]

	jss.mu.RUnlock()

	res := make([]filesPck.File, 0, len(files))
	for _, f := range files {
		if ids.hasID(f.ID) {
			res = append(res, f)
		}
	}

	return res, nil
}

func (jss *jsonShareStorage) filterTags(token string, tags tagsPck.Tags) (tagsPck.Tags, error) {
	ids, err := jss.getFilesIDs(token)
	if err != nil {
		return tags, err
	}

	// We don't need to use mutex

	result := make(tagsPck.Tags)

	files := jss.fileStorage.GetFiles(ids...)

	for id := range tags {
	searchLoop:
		for i := range files {
			for j := range files[i].Tags {
				if files[i].Tags[j] == id {
					result[id] = tags[id]
					break searchLoop
				}
			}
		}
	}

	return result, nil
}

func (jss *jsonShareStorage) shutdown() error {
	jss.mu.Lock()
	jss.mu.Unlock()

	jss.write()

	return nil
}
