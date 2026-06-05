package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		latency := time.Since(start).Milliseconds()

		requestID, _ := c.Get("request_id").(string)
		slog.Info("HTTP Request",
			slog.String("method", c.Request().Method),
			slog.String("path", c.Request().RequestURI),
			slog.Int("status", c.Response().Status),
			slog.Int64("latency_ms", latency),
			slog.String("request_id", requestID),
		)

		return err
	}
}
