package rest

import (
	"net/http"
	"strings"
)

// DeleteAsset обрабатывает запросы на удаление ассетов.
// Он проверяет наличие имени ассета в URL, получает сессию пользователя,
// удаляет ассет с использованием сервиса и отправляет ответ.
func (h Handler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
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

	// Удаляем ассет с использованием сервиса
	err = h.service.DeleteAsset(ctx, session.UserID, assetName)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusNoContent)
}
