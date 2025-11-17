package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noctusha/url-shortener/internal/config"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	pgxCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("%s: pgxpool parse config error", op)
	}

	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MinConns = cfg.MinConns

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: pgxpool connect error", op)
	}

	for i := 0; i < 10; i++ {
		err = dbpool.Ping(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s: pgxpool connect error", op)
		}
		time.Sleep(250 * time.Millisecond)
	}

	slog.Info("connected to postgres",
		"addr", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		"db", cfg.Name)

	_, err = dbpool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS url(
	  id INTEGER PRIMARY KEY,
	  alias TEXT NOT NULL UNIQUE,
	  url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: dbpool}, nil
}

func (s *Storage) Conn() *pgxpool.Pool {
	return s.db
}

func (s *Storage) Close() {
	s.db.Close()
}
