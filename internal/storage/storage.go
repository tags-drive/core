package storage

import (
	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/files/aggregation"
)

// Errors
var (
	ErrFileIsNotExist    = files.ErrFileIsNotExist
	ErrAlreadyExist      = files.ErrAlreadyExist
	ErrFileDeletedAgain  = files.ErrFileDeletedAgain
	ErrOffsetOutOfBounds = files.ErrOffsetOutOfBounds
	//
	ErrBadExpessionSyntax = aggregation.ErrBadSyntax
)
