package shortener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/noctusha/url-shortener/internal/storage"
)

const AliasLength = 6

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

func (s *Service) SaveURL(ctx context.Context, url, alias string) (int32, string, error) {
	const op = "service.shortener.SaveURL"
	if alias == "" {
		alias = Random(AliasLength)
	}

	id, err := s.repo.Save(ctx, url, alias)
	if err != nil {
		if errors.Is(err, storage.ErrAliasExists) {
			s.log.Warn("alias already exists", "alias", alias)
			return 0, "", ErrAliasAlreadyExists
		}
		s.log.Error("failed to save url", "error", err)
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("url saved",
		"id", id,
		"alias", alias,
		"url", url,
	)

	return id, alias, nil
}
