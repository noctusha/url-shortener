package shortener

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrURLNotFound        = fmt.Errorf("url not found")
)

type URLRepository interface {
	SaveURL(ctx context.Context, url string, alias string) (int32, error)
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
}
