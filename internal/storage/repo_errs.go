package storage

import (
	"errors"
	"fmt"
)

var (
	ErrAliasExists = errors.New("alias exists")
	ErrURLNotFound = fmt.Errorf("url not found")
)
