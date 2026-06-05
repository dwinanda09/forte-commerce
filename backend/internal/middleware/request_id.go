package middleware

import (
	"context"

	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Request().Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Response().Header().Set("X-Request-ID", requestID)

		// Inject into standard context so util.Logger can read it
		ctx := context.WithValue(c.Request().Context(), util.RequestIDKey, requestID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
