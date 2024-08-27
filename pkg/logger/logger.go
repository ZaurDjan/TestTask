package logger

import (
	"context"
	"log/slog"
)

type key struct{}

// ToCtx добавляет логгер в контекст
func ToCtx(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, key{}, l)
}

// FromCtx извлекает логгер из контекста
func FromCtx(ctx context.Context) *slog.Logger {
	log, _ := ctx.Value(key{}).(*slog.Logger)
	if log != nil {
		return log
	}
	return slog.Default()
}

// With добавляет аргументы к логгеру в контексте
func With(ctx context.Context, args ...any) context.Context {
	return ToCtx(ctx, FromCtx(ctx).With(args...))
}
