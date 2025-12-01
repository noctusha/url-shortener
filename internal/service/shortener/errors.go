package shortener

import "errors"

var (
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrURLNotFound        = errors.New("url not found")
)
