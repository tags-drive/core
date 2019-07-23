package share

import (
	clog "github.com/ShoshinNikita/log/v2"
	"github.com/pkg/errors"

	filesPck "github.com/tags-drive/core/internal/storage/files"
	tagsPck "github.com/tags-drive/core/internal/storage/tags"
)

var (
	ErrInvalidToken = errors.New("invalid share token")
)

type ShareStorage struct {
	storage internalStorage
}

func NewShareStorage(cnf Config, fs FileStorage, lg *clog.Logger) (*ShareStorage, error) {
	storage := &ShareStorage{}

	// Init an internal storage
	st := newJsonShareStorage(cnf, fs, lg)

	err := st.init()
	if err != nil {
		return nil, err
	}

	storage.storage = st

	return storage, nil
}

func (st ShareStorage) GetAllTokens() map[string][]int {
	return st.storage.getAllTokens()
}

func (st ShareStorage) CreateToken(ids []int) (token string) {
	return st.storage.createToken(ids)
}

func (st ShareStorage) DeleteToken(token string) {
	st.storage.deleteToken(token)
}

func (st ShareStorage) GetFilesIDs(token string) ([]int, error) {
	return st.storage.getFilesIDs(token)
}

func (st ShareStorage) CheckToken(token string) bool {
	return st.storage.checkToken(token)
}

func (st ShareStorage) CheckFile(token string, id int) bool {
	return st.storage.checkFile(token, id)
}

func (st ShareStorage) DeleteFile(id int) {
	st.storage.deleteFile(id)
}

func (st ShareStorage) FilterFiles(token string, files []filesPck.File) ([]filesPck.File, error) {
	return st.storage.filterFiles(token, files)
}

func (st ShareStorage) FilterTags(token string, tags tagsPck.Tags) (tagsPck.Tags, error) {
	return st.storage.filterTags(token, tags)
}

func (st ShareStorage) Shutdown() error {
	return st.storage.shutdown()
}
