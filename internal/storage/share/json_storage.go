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

// fileIDs is a slice of sorted file ids
type filesIDs []int

func newFileIDs(ids []int) filesIDs {
	newIDs := append(ids[:0:0], ids...)

	sort.Ints(newIDs)

	return filesIDs(ids)
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

// jsonShareStorage implements ShareStorageInterface interface
type jsonShareStorage struct {
	config Config

	tokens map[string]filesIDs
	mu     sync.RWMutex

	logger *clog.Logger
}

func newJsonShareStorage(cnf Config, lg *clog.Logger) *jsonShareStorage {
	return &jsonShareStorage{
		config: cnf,
		tokens: make(map[string]filesIDs),
		logger: lg,
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

func (jss *jsonShareStorage) CreateToken(ids []int) (token string) {
	jss.mu.Lock()
	defer func() {
		jss.mu.Unlock()
		jss.write()
	}()

	token = utils.GenerateRandomString(maxTokenSize)
	jss.tokens[token] = newFileIDs(ids)

	return token
}

func (jss *jsonShareStorage) CheckToken(token string) bool {
	jss.mu.RLock()
	defer jss.mu.RUnlock()

	_, ok := jss.tokens[token]
	return ok
}

func (jss *jsonShareStorage) CheckFile(token string, id int) bool {
	jss.mu.RLock()
	defer jss.mu.RUnlock()

	ids, ok := jss.tokens[token]
	if !ok {
		return false
	}

	return ids.hasID(id)
}

func (jss *jsonShareStorage) DeleteFile(id int) {
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

func (jss *jsonShareStorage) FilterFiles(token string, files []filesPck.File) ([]filesPck.File, error) {
	if !jss.CheckToken(token) {
		return nil, errors.New("invalid token")
	}

	ids := jss.tokens[token]

	res := make([]filesPck.File, 0, len(files))
	for _, f := range files {
		if ids.hasID(f.ID) {
			res = append(res, f)
		}
	}

	return res, nil
}

// TODO
func (jss *jsonShareStorage) FilterTags(token string, tags tagsPck.Tags) (tagsPck.Tags, error) {
	return tags, nil
}

func (jss *jsonShareStorage) Shutdown() error {
	jss.mu.Lock()
	jss.mu.Unlock()

	jss.write()

	return nil
}
