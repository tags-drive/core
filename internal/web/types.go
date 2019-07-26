package web

import (
	"context"
	"time"

	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
)

type Config struct {
	Debug   bool
	Version string

	Port  string
	IsTLS bool

	Login    string
	Password string

	SkipLogin      bool
	AuthCookieName string
	MaxTokenLife   time.Duration
}

type AuthServiceInterface interface {
	CheckToken(token string) bool

	GenerateToken() string

	AddToken(token string)

	DeleteToken(token string)

	Shutdown() error
}

type ShareServiceInterface interface {
	// Tokens

	GetAllTokens() map[string][]int

	CheckToken(token string) bool

	CreateToken(filesIDs []int) (token string)

	GetFilesIDs(token string) ([]int, error)

	DeleteToken(token string)

	// Files and tags

	CheckFile(token string, id int) bool

	FilterFiles(token string, files []files.File) ([]files.File, error)

	FilterTags(token string, tags tags.Tags) (tags.Tags, error)

	DeleteFile(id int)

	//

	Shutdown() error
}

// requestState stores state of current request. It is passed by request's context
type requestState struct {
	// authorized it always true. It can be false only when shareAccess is true.
	// So, handlers must process shareAccess first
	authorized bool

	shareAccess bool
	// shareToken can't be empty when shareAccess is true
	shareToken string
}

// requestStateKey is a key for an instance of requestState within context
const requestStateKey = "requestState"

func storeRequestState(ctx context.Context, state *requestState) context.Context {
	return context.WithValue(ctx, requestStateKey, state)
}

func getRequestState(ctx context.Context) (*requestState, bool) {
	state, ok := ctx.Value(requestStateKey).(*requestState)
	if !ok {
		return nil, false
	}

	if state == nil {
		return nil, false
	}

	return state, true
}
