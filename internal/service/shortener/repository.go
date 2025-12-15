package shortener

import (
	"context"
	"time"
)

type URLRepository interface {
	Save(ctx context.Context, url string, alias string, expireAt *time.Time) (int32, error)
	Get(ctx context.Context, alias string) (string, error)
	Delete(ctx context.Context, alias string) error
}
