package shortener

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/noctusha/url-shortener/internal/storage/postgres"
)

type Service struct {
	repo URLRepository
	log  *slog.Logger
}

func NewService(repo URLRepository, log *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) URLSave(ctx context.Context, url string, alias string) (int32, error) {
	const op = "service.shortener.URLSave"
	id, err := s.repo.SaveURL(ctx, url, alias)
	if err != nil {
		if err == postgres.ErrAliasExists {
			s.log.Warn("alias already exists", "alias", alias)
		} else {
			s.log.Error("failed to save url", "error", err)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	s.log.Info("url saved",
		"id", id,
		"alias", alias,
		"url", url,
	)

	return id, nil
}
