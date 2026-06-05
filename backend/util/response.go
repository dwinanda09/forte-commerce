package util

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Meta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    Meta        `json:"meta"`
}

func OK(c echo.Context, data interface{}) error {
	requestID := ""
	if rid := c.Get("request_id"); rid != nil {
		requestID = rid.(string)
	}

	resp := Response{
		Success: true,
		Data:    data,
		Meta: Meta{
			RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.JSON(200, resp)
}

func Fail(c echo.Context, status int, code, message string) error {
	requestID := ""
	if rid := c.Get("request_id"); rid != nil {
		requestID = rid.(string)
	}

	resp := Response{
		Success: false,
		Message: message,
		Meta: Meta{
			RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.JSON(status, resp)
}
