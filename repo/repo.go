package repo

import (
	"context"

	"time"

	"github.com/ZaurDjan/TestTask/pkg/models"
)

type Repo interface {
	GetUser(ctx context.Context, login string) (models.User, error)
	DeleteSession(ctx context.Context, id string) error
	CreateSession(ctx context.Context, sessionID string, id string, ipAddress string) error
	GetSession(ctx context.Context, sessionID string) (models.Session, error)
	InsertAsset(ctx context.Context, assetName string, userID string, data []byte, createdAt time.Time) error
	GetAsset(ctx context.Context, assetName string, userID string) (*models.Asset, error)
	ListAssets(ctx context.Context, userID string) ([]models.Asset, error)
	DeleteAsset(ctx context.Context, assetName string, userID string) error
}
