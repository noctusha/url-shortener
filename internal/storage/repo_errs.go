package storage

import (
	"errors"
)

var (
	ErrAliasExists = errors.New("alias exists")
	ErrURLNotFound = errors.New("url not found")
)
