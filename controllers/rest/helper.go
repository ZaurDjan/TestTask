package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/ZaurDjan/TestTask/errcodes"
	"github.com/ZaurDjan/TestTask/pkg/logger"
	"github.com/ZaurDjan/TestTask/pkg/models"
)

// AuthMiddleware добавляет проверку аутентификации для обработчика HTTP-запросов.
// Он проверяет наличие заголовка Authorization, извлекает токен и проверяет его валидность.
func (h Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		// Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Извлекаем токен из заголовка
		token := ""
		fmt.Sscanf(authHeader, "Bearer %s", &token)
		if token == "" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Проверяем валидность токена
		session, err := h.service.GetSession(ctx, token)
		if err != nil {
			WriteError(w, err)
			return
		}

		// Добавляем сессию в контекст
		ctx = WithSession(ctx, session)
		r = r.WithContext(ctx)

		// Передаем управление следующему обработчику
		next(w, r)
	}
}

// sessionKey представляет ключ для хранения сессии в контексте.
type sessionKey struct {
}

// GetSession извлекает сессию из контекста.
func GetSession(ctx context.Context) (models.Session, error) {
	s, ok := ctx.Value(sessionKey{}).(models.Session)
	if !ok {
		return s, errors.New("Unauthorized")
	}
	return s, nil
}

// WithSession добавляет сессию в контекст.
func WithSession(ctx context.Context, session models.Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, session)
}

// WriteError записывает ошибку в HTTP-ответ в формате JSON.
func WriteError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError

	var e errcodes.Error
	if errors.As(err, &e) {
		switch e.Reason() {
		case errcodes.ReasonBadRequest:
			code = http.StatusBadRequest
		case errcodes.ReasonNotFound:
			code = http.StatusNotFound
		case errcodes.ReasonUnauthenticated:
			code = http.StatusUnauthorized
		case errcodes.ReasonAccessDenied:
			code = http.StatusForbidden
		}
	}
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func WriteErrorString(w http.ResponseWriter, message string, statusCode int) {

	w.WriteHeader(statusCode)

	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// LogMiddleware добавляет логирование для обработчика HTTP-запросов.
// Он логирует запросы и передает управление следующему обработчику.
func (h Handler) LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем запрос
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			WriteError(w, err)
			return
		}
		r = r.WithContext(logger.With(r.Context(), "request", string(b)))
		logger.FromCtx(r.Context()).Error("test")

		// Передаем управление следующему обработчику
		next(w, r)
	}
}
