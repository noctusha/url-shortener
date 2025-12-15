package shortener

import "errors"

var (
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrURLNotFound        = errors.New("url not found")
	ErrInvalidAlias       = errors.New("alias is empty")
	ErrURLExpired         = errors.New("url is expired")
)
