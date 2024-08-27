package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"time"

	"github.com/ZaurDjan/TestTask/errcodes"
	"github.com/ZaurDjan/TestTask/pkg/models"
)

// Authorization выполняет авторизацию пользователя.
// Проверяет логин и пароль, удаляет предыдущие сессии и создает новую сессию.
func (s *Service) Authorization(ctx context.Context, login, password, ipAddress string) (models.Session, error) {

	// Получаем пользователя из базы данных по логину
	dbUser, err := s.repo.GetUser(ctx, login)
	if err != nil {
		return models.Session{}, errcodes.Errorf(errcodes.ReasonUnauthenticated, "getting user: %w", err)
	}

	// Проверка пароля
	if !checkPasswordHash(password, dbUser.PasswordHash) {
		return models.Session{}, errcodes.New(errcodes.ReasonUnauthenticated, "invalid login or password")
	}

	// Удаление всех предыдущих сессий пользователя
	err = s.repo.DeleteSession(ctx, dbUser.ID)
	if err != nil {
		return models.Session{}, fmt.Errorf("deleting previous session: %w", err)
	}

	// Генерация нового идентификатора сессии
	sessionID := generateToken()

	// Создание новой сессии
	err = s.repo.CreateSession(ctx, sessionID, dbUser.ID, ipAddress)
	if err != nil {
		return models.Session{}, fmt.Errorf("creating session: %w", err)
	}

	// Возвращаем созданную сессию
	return models.Session{ID: sessionID, UserID: dbUser.ID, CreatedAt: time.Now(), IPAddress: ipAddress}, nil
}

// checkPasswordHash проверяет соответствие пароля и его хеша.
func checkPasswordHash(password, hash string) bool {
	hashedPassword := hashPassword(password)
	return hashedPassword == hash
}

// hashPassword хеширует пароль с использованием SHA-256.
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// generateToken генерирует случайный токен для сессии.
func generateToken() string {
	return hex.EncodeToString(randomBytes(16))
}

// randomBytes генерирует случайный массив байтов заданной длины.
func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
