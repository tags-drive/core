package web

import (
	"context"
	"time"
)

type Config struct {
	Debug bool

	DataFolder string

	Port  string
	IsTLS bool

	Login          string
	Password       string
	SkipLogin      bool
	AuthCookieName string
	MaxTokenLife   time.Duration

	AuthTokensJSONFile  string
	ShareTokensJSONFile string

	Encrypt    bool
	PassPhrase [32]byte

	Version string
}

// ServerInterface provides methods for interactions web server
type ServerInterface interface {
	Start() error

	// Shutdown gracefully shutdowns server
	Shutdown() error
}

// requestState stores state of current request. It is passed by request's context
type requestState struct {
	authorized bool

	shareAccess bool
	shareToken  string // shareToken can't be empty when shareAccess is true
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
