package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZaurDjan/TestTask/controllers/rest"
	"github.com/ZaurDjan/TestTask/pkg/logger"
	"github.com/ZaurDjan/TestTask/repo/postgress"
	"github.com/ZaurDjan/TestTask/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()
	// Устанавливаем логгер по умолчанию
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	// Получаем строку подключения к базе данных из переменной окружения
	dsn := os.Getenv("db_conn")

	var err error
	// Подключаемся к базе данных
	database, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return
	}
	defer database.Close()

	// Выполняем миграции базы данных
	err = runMigrations(ctx, database, "migrations")
	if err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return
	}

	logger.FromCtx(ctx).Info("Migrations applied successfully!")

	// Создаем новый репозиторий
	repo := postgress.NewRepo(database)

	// Создаем новый сервис
	service := service.New(repo)

	// Создаем новый контроллер REST API
	rest := rest.New(service)

	// Настраиваем маршруты HTTP
	http.HandleFunc("POST /api/auth", rest.LogMiddleware(rest.Authorization))
	http.HandleFunc("POST /api/upload-asset/", rest.AuthMiddleware(rest.UploadAsset))
	http.HandleFunc("GET /api/asset/", rest.AuthMiddleware(rest.Download))
	http.HandleFunc("GET /api/assets", rest.AuthMiddleware(rest.ListAssets))
	http.HandleFunc("DELETE /api/asset/", rest.AuthMiddleware(rest.DeleteAsset)) // Добавляем новый маршрут для удаления файлов

	// Загружаем сертификат и ключ для HTTPS
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return
	}

	// Настраиваем и запускаем HTTPS сервер
	s := &http.Server{
		Addr:    ":8080",
		Handler: nil,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{cer},
			InsecureSkipVerify:       true,
		},

		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = s.ListenAndServeTLS("", "")
	if err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return
	}
}

// runMigrations выполняет миграции базы данных.
// Он читает SQL файлы из указанной директории и применяет их,
// если они еще не были применены.
func runMigrations(ctx context.Context, conn *pgxpool.Pool, migrationsDir string) error {
	// Создаем таблицу для отслеживания примененных миграций, если она не существует
	_, err := conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS migrations (id SERIAL PRIMARY KEY,name VARCHAR(255) NOT NULL,applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Читаем файлы миграций из указанной директории
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationName := file.Name()
			migrationPath := filepath.Join(migrationsDir, migrationName)

			// Проверяем, была ли миграция уже применена
			var count int
			err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM migrations WHERE name = $1", migrationName).Scan(&count)
			if err != nil {
				return fmt.Errorf("failed to check migration status: %w", err)
			}

			if count == 0 {
				// Применяем миграцию
				err := applyMigration(ctx, conn, migrationPath)
				if err != nil {
					return fmt.Errorf("failed to apply migration %s: %w", migrationName, err)
				}

				// Записываем миграцию как примененную
				_, err = conn.Exec(ctx, "INSERT INTO migrations (name) VALUES ($1)", migrationName)
				if err != nil {
					return fmt.Errorf("failed to record migration %s: %w", migrationName, err)
				}
			} else {
				logger.FromCtx(ctx).Info("Migration  has already been applied, skipping", slog.String("name", migrationName))
			}
		}
	}

	return nil
}

// applyMigration применяет миграцию из указанного файла.
// Он читает SQL команды из файла и выполняет их в базе данных.
func applyMigration(ctx context.Context, conn *pgxpool.Pool, migrationPath string) error {
	// Читаем файл миграции
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Выполняем SQL команды из файла миграции
	_, err = conn.Exec(ctx, string(migration))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	logger.FromCtx(ctx).Info("Applied migration", slog.String("name", migrationPath))
	return nil
}
