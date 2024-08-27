package rest

import "github.com/ZaurDjan/TestTask/service"

// New создает новый экземпляр Handler.
func New(s *service.Service) Handler {
	return Handler{service: s}
}

// Handler представляет обработчик HTTP-запросов
type Handler struct {
	service *service.Service
}
