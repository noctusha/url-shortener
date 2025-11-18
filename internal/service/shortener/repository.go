package shortener

import (
	"context"
)

//var (
//	ErrAliasExists = fmt.Errorf("alias URL already exists")
//	ErrURLNotFound = fmt.Errorf("url not found")
//)

type URLRepository interface {
	SaveURL(ctx context.Context, url string, alias string) (int32, error)
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
	//GetURLByAlias(ctx context.Context, alias string) (string, error)
}

//type Repository struct {
//	q  *sqlc.Queries
//	db *pgxpool.Pool
//}
//
//func NewURLRepository(db *pgxpool.Pool) *Repository {
//	return &Repository{q: sqlc.New(db), db: db}
//}
//
//func (r *Repository) SaveURL(ctx context.Context, url, alias string) (int32, error) {
//	id, err := r.q.SaveURL(ctx, sqlc.SaveURLParams{
//		Url:   url,
//		Alias: alias,
//	})
//	if err != nil {
//		// 23505 = unique_violation
//		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
//			return 0, ErrAliasExists
//		}
//		return 0, err
//	}
//	return id, nil
//}
