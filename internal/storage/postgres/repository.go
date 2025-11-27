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

func (r *Repository) SaveURL(ctx context.Context, url, alias string) (int32, error) {
	const op = "storage.postgres.SaveURL"
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

func (r *Repository) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	url, err := r.q.GetURL(ctx, alias)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (r *Repository) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.postgres.DeleteURL"
	err := r.q.DeleteURL(ctx, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
