package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/noctusha/url-shortener/internal/storage"
	sqlc "github.com/noctusha/url-shortener/internal/storage/sqlc"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q  *sqlc.Queries
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		q:  sqlc.New(db),
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, url, alias string) (int32, error) {
	const op = "storage.postgres.Save"
	id, err := r.q.SaveURL(ctx, sqlc.SaveURLParams{
		Url:   url,
		Alias: alias,
	})
	if err != nil {
		// 23505 = unique_violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *Repository) Get(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.Get"
	url, err := r.q.GetURL(ctx, alias)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (r *Repository) Delete(ctx context.Context, alias string) error {
	const op = "storage.postgres.Delete"
	affected, err := r.q.DeleteURL(ctx, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if affected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	return nil
}
