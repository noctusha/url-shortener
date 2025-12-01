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
	pgxCfg.MaxConnIdleTime = 5 * time.Minute
	pgxCfg.MaxConnLifetime = 30 * time.Minute
	pgxCfg.HealthCheckPeriod = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: pgxpool connect error: %w", op, err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
	defer pingCancel()
	if err := dbpool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	log.Info("connected to postgres",
		"addr", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		"database", cfg.Name)

	return &Storage{db: dbpool, log: log}, nil
}

func (s *Storage) Conn() *pgxpool.Pool {
	return s.db
}

func (s *Storage) Close() {
	s.db.Close()
}
