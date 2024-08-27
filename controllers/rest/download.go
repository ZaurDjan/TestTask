package rest

import (
	"log"
	"net/http"
	"strings"
)

// Download обрабатывает запросы на скачивание ассетов.
// Он проверяет наличие имени ассета в URL, получает сессию пользователя,
// извлекает данные ассета из сервиса и отправляет их в ответе.
func (h Handler) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Извлекаем имя ассета из URL
	assetName := strings.TrimPrefix(r.URL.Path, "/api/asset/")
	if assetName == "" {
		WriteErrorString(w, "asset name missing", http.StatusBadRequest)
		return
	}

	// Получаем сессию из контекста
	session, err := GetSession(ctx)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем данные ассета из сервиса
	data, err := h.service.GetAsset(ctx, session.UserID, assetName)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Логируем успешное скачивание ассета
	log.Printf("User downloaded asset %s with data: %s", assetName, data.Data)

	// Устанавливаем заголовок Content-Type в application/octet-stream
	w.Header().Set("Content-Type", "application/octet-stream")

	// Отправляем данные ассета в ответе
	w.Write(data.Data)
}
