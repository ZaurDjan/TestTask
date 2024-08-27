package rest

import (
	"encoding/json"
	"net/http"
)

// ListAssets обрабатывает запросы на получение списка ассетов.
// Он проверяет наличие сессии, получает список ассетов для пользователя
// и возвращает их в формате JSON.
func (h Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем сессию из контекста
	session, err := GetSession(ctx)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем список ассетов для пользователя
	assets, err := h.service.ListAssets(ctx, session.UserID)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Кодируем список ассетов в JSON и отправляем в ответе
	json.NewEncoder(w).Encode(assets)
}
