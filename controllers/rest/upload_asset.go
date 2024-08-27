package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// UploadAsset обрабатывает запросы на загрузку ассетов
func (h Handler) UploadAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем сессию из контекста
	session, err := GetSession(ctx)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем имя ассета из URL
	assetName := strings.TrimPrefix(r.URL.Path, "/api/upload-asset/")
	if assetName == "" {
		WriteErrorString(w, "asset name missing", http.StatusBadRequest)
		return
	}

	// Читаем данные из тела запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		WriteErrorString(w, "no data provided", http.StatusBadRequest)
		return
	}

	// Загружаем ассет с использованием сервиса
	err = h.service.UploadAsset(ctx, session.UserID, assetName, data)
	if err != nil {
		WriteErrorString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Логируем успешную загрузку ассета
	log.Printf("user uploaded asset %s with data: %s", assetName, string(data))

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
