package service

import (
	"context"

	"time"

	"github.com/ZaurDjan/TestTask/errcodes"
	"github.com/ZaurDjan/TestTask/pkg/models"
)

// UploadAsset загружает ассет в репозиторий.
func (s *Service) UploadAsset(ctx context.Context, userID, assetName string, data []byte) error {
	return s.repo.InsertAsset(ctx, assetName, userID, data, time.Now())
}

// GetAsset получает ассет из репозитория.
func (s *Service) GetAsset(ctx context.Context, userID, assetName string) (models.Asset, error) {
	data, err := s.repo.GetAsset(ctx, assetName, userID)
	if err != nil {
		return models.Asset{}, err
	}
	if data == nil {
		return models.Asset{}, errcodes.Errorf(errcodes.ReasonNotFound, "%q asset does not exist", assetName)
	}
	return *data, nil
}

// ListAssets возвращает список всех ассетов для данного пользователя.
func (s *Service) ListAssets(ctx context.Context, userID string) ([]models.Asset, error) {
	rows, err := s.repo.ListAssets(ctx, userID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// DeleteAsset удаляет ассет из репозитория.
func (s *Service) DeleteAsset(ctx context.Context, userID, assetName string) error {
	err := s.repo.DeleteAsset(ctx, assetName, userID)
	if err != nil {
		return err
	}
	return nil
}
