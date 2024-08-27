package models

import "time"

// User представляет модель пользователя.
type User struct {
	ID           string    `json:"id"`            // Уникальный идентификатор пользователя
	Login        string    `json:"login"`         // Логин пользователя
	Password     string    `json:"password"`      // Пароль пользователя (не храним в базе данных)
	PasswordHash string    `json:"password_hash"` // Хеш пароля пользователя
	CreatedAt    time.Time `json:"created_at"`    // Время создания пользователя
}

// Session представляет модель сессии пользователя.
type Session struct {
	ID        string    `json:"id"`         // Уникальный идентификатор сессии
	UserID    string    `json:"user_id"`    // Идентификатор пользователя, которому принадлежит сессия
	CreatedAt time.Time `json:"created_at"` // Время создания сессии
	IPAddress string    `json:"ip_address"` // IP-адрес пользователя, использовавшего сессию
}

// Asset представляет модель ассета (файла).
type Asset struct {
	ID        string    `json:"id"`         // Уникальный идентификатор ассета
	Data      []byte    `json:"data"`       // Данные ассета в виде байтового массива
	Name      string    `json:"name"`       // Имя ассета
	UserID    string    `json:"user_id"`    // Идентификатор пользователя, загрузившего ассет
	CreatedAt time.Time `json:"created_at"` // Время создания ассета
}
