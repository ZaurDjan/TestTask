package postgress

import (
	"context"

	"time"

	"github.com/ZaurDjan/TestTask/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NewRepo создает новый экземпляр Repo.
func NewRepo(pool *pgxpool.Pool) Repo {
	return Repo{db: pool}
}

// Repo представляет репозиторий для работы с PostgreSQL.
type Repo struct {
	db *pgxpool.Pool
}

// Получение пользователя
func (r Repo) GetUser(ctx context.Context, login string) (models.User, error) {
	var dbUser models.User
	err := r.db.QueryRow(ctx, "SELECT id, login, password_hash FROM users WHERE login = $1", login).Scan(&dbUser.ID, &dbUser.Login, &dbUser.PasswordHash)
	return dbUser, err
}

// Удаление сессии.
func (r Repo) DeleteSession(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sessions WHERE uid = $1", id)
	return err
}

// Создание сессии.
func (r Repo) CreateSession(ctx context.Context, sessionID string, id string, ipAddress string) error {
	_, err := r.db.Exec(ctx, "INSERT INTO sessions (id, uid, ip_address) VALUES ($1, $2, $3)", sessionID, id, ipAddress)
	return err
}

// InsertAsset вставляет новый ассет в базу данных.
func (r Repo) InsertAsset(ctx context.Context, assetName string, userID string, data []byte, createdAt time.Time) error {
	_, err := r.db.Exec(ctx, "INSERT INTO assets (name, uid, data, created_at) VALUES ($1, $2, $3, $4)", assetName, userID, data, time.Now())
	return err
}

// GetAsset получает ассет из базы данных.
func (r Repo) GetAsset(ctx context.Context, assetName string, userID string) (*models.Asset, error) {

	var data models.Asset

	err := r.db.QueryRow(ctx, "SELECT data FROM assets WHERE name = $1 AND uid = $2", assetName, userID).Scan(&data.Data)
	return &data, err
}

// ListAssets возвращает список всех ассетов для данного пользователя.
func (r Repo) ListAssets(ctx context.Context, userID string) ([]models.Asset, error) {
	rows, err := r.db.Query(ctx, "SELECT name, uid, created_at FROM assets WHERE uid = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset

	for rows.Next() {
		var asset models.Asset
		if err := rows.Scan(&asset.Name, &asset.UserID, &asset.CreatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

// DeleteAsset удаляет ассет из базы данных.
func (r Repo) DeleteAsset(ctx context.Context, assetName string, userID string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM assets WHERE name = $1 AND uid = $2", assetName, userID)
	return err
}

// GetSession получает сессию пользователя по токену.
func (r Repo) GetSession(ctx context.Context, sessionID string) (models.Session, error) {
	var session models.Session
	err := r.db.QueryRow(ctx, "SELECT id, uid, created_at,ip_address  FROM sessions WHERE id = $1", sessionID).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.IPAddress)
	return session, err
}
