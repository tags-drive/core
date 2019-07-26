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

type ShareService struct {
	storage internalStorage
}

func NewShareStorage(cnf Config, fs FileStorage, lg *clog.Logger) (*ShareService, error) {
	storage := &ShareService{}

	// Init an internal storage
	st := newJsonShareStorage(cnf, fs, lg)

	err := st.init()
	if err != nil {
		return nil, err
	}

	storage.storage = st

	return storage, nil
}

func (st ShareService) GetAllTokens() map[string][]int {
	return st.storage.getAllTokens()
}

func (st ShareService) CreateToken(ids []int) (token string) {
	return st.storage.createToken(ids)
}

func (st ShareService) DeleteToken(token string) {
	st.storage.deleteToken(token)
}

func (st ShareService) GetFilesIDs(token string) ([]int, error) {
	return st.storage.getFilesIDs(token)
}

func (st ShareService) CheckToken(token string) bool {
	return st.storage.checkToken(token)
}

func (st ShareService) CheckFile(token string, id int) bool {
	return st.storage.checkFile(token, id)
}

func (st ShareService) DeleteFile(id int) {
	st.storage.deleteFile(id)
}

func (st ShareService) FilterFiles(token string, files []filesPck.File) ([]filesPck.File, error) {
	return st.storage.filterFiles(token, files)
}

func (st ShareService) FilterTags(token string, tags tagsPck.Tags) (tagsPck.Tags, error) {
	return st.storage.filterTags(token, tags)
}

func (st ShareService) Shutdown() error {
	return st.storage.shutdown()
}
