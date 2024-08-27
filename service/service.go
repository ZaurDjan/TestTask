package service

import (
	"context"

	"github.com/ZaurDjan/TestTask/pkg/models"
	"github.com/ZaurDjan/TestTask/repo"
)

// New создает новый экземпляр Service.
func New(rep repo.Repo) *Service {
	return &Service{
		repo: rep,
	}
}

// Service представляет сервис для работы с данными.
type Service struct {
	repo repo.Repo
}

// GetSession получает сессию пользователя по идентификатору сессии.
func (s *Service) GetSession(ctx context.Context, sessionID string) (models.Session, error) {
	return s.repo.GetSession(ctx, sessionID)
}
