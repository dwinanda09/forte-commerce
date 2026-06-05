package util

import (
	"context"
	"log/slog"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

type Logger struct {
	log *slog.Logger
}

func NewLogger() *Logger {
	return &Logger{
		log: slog.Default(),
	}
}

func requestIDAttr(ctx context.Context) slog.Attr {
	if id, ok := ctx.Value(RequestIDKey).(string); ok && id != "" {
		return slog.String("request_id", id)
	}
	return slog.Attr{}
}

func (l *Logger) Start(ctx context.Context, method string) time.Time {
	l.log.InfoContext(ctx, "START", slog.String("method", method), requestIDAttr(ctx))
	return time.Now()
}

func (l *Logger) Finish(ctx context.Context, method string, start time.Time, err error) {
	duration := time.Since(start).Milliseconds()
	if err != nil {
		l.log.ErrorContext(ctx, "FINISH",
			slog.String("method", method),
			slog.Int64("processed_ms", duration),
			slog.String("error", err.Error()),
			requestIDAttr(ctx),
		)
	} else {
		l.log.InfoContext(ctx, "FINISH",
			slog.String("method", method),
			slog.Int64("processed_ms", duration),
			requestIDAttr(ctx),
		)
	}
}
