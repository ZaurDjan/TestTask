package rest

import (
	"encoding/json"

	"net/http"

	"github.com/ZaurDjan/TestTask/pkg/models"
)

// Authorization обрабатывает запросы на авторизацию.
// Он декодирует учетные данные из тела запроса, выполняет авторизацию
// с использованием сервиса и возвращает токен сессии в ответе
func (h Handler) Authorization(w http.ResponseWriter, r *http.Request) {

	// Получаем контекст запроса
	ctx := r.Context()

	// Создаем переменную для хранения данных пользователя
	var user models.User

	// Получаем IP-адрес клиента из запроса
	ipAddress := r.RemoteAddr

	// Декодируем JSON-данные из тела запроса в структуру user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Если произошла ошибка при декодировании, отправляем ответ с ошибкой 400 (Bad Request)
		WriteErrorString(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем метод авторизации сервиса, передавая контекст, логин, пароль и IP-адрес пользователя
	session, err := h.service.Authorization(ctx, user.Login, user.Password, ipAddress)
	if err != nil {
		// Если произошла ошибка при авторизации, отправляем ответ с соответствующей ошибкой
		WriteErrorString(w, err.Error(), http.StatusUnauthorized)
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Кодируем токен сессии в JSON и отправляем в ответе
	json.NewEncoder(w).Encode(map[string]string{"session_id": session.ID})
}
