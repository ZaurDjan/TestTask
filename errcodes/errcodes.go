package errcodes

import (
	"errors"
	"fmt"
)

// Reason представляет тип причины ошибки.
type Reason int

const (
	// ReasonInternal обозначает внутреннюю ошибку.
	ReasonInternal Reason = iota
	// ReasonBadRequest обозначает ошибку неверного запроса.
	ReasonBadRequest
	// ReasonNotFound обозначает ошибку, когда ресурс не найден.
	ReasonNotFound
	// ReasonUnauthenticated обозначает ошибку неаутентифицированного доступа.
	ReasonUnauthenticated
	// ReasonAccessDenied обозначает ошибку отказа в доступе.
	ReasonAccessDenied
)

// New создает новую ошибку с указанной причиной и сообщением.
func New(r Reason, message string) Error {
	return err{
		reason:  r,
		message: errors.New(message),
	}
}

// Errorf создает новую ошибку с указанной причиной и форматированным сообщением.
func Errorf(r Reason, format string, args ...any) Error {
	return err{
		reason:  r,
		message: fmt.Errorf(format, args...),
	}
}

// Error представляет интерфейс для ошибок с причиной.
type Error interface {
	error
	Reason() Reason
}

// err представляет реализацию интерфейса Error.
type err struct {
	reason  Reason
	message error
}

// err представляет реализацию интерфейса Error.
func (e err) Error() string {
	return e.message.Error()
}

// Reason возвращает причину ошибки.
func (e err) Reason() Reason {
	return e.reason
}

// Unwrap возвращает вложенную ошибку.
func (e err) Unwrap() error {
	return e.message
}
