package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noctusha/url-shortener/internal/config"
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) (*Storage, error) {
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
		return nil, fmt.Errorf("%s: pgxpool parse config error: %w", op, err)
	}

	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MinConns = cfg.MinConns

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: pgxpool connect error: %w", op, err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	log.Info("connected to postgres",
		"addr", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		"sql", cfg.Name)

	return &Storage{db: dbpool, log: log}, nil
}

func (s *Storage) Conn() *pgxpool.Pool {
	return s.db
}

func (s *Storage) Close() {
	s.db.Close()
}
