package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newEchoContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestOK(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		requestID      string
		expectSuccess  bool
		expectData     interface{}
	}{
		{
			name:          "ok with object data",
			data:          map[string]interface{}{"id": 1, "name": "test"},
			requestID:     "req-123",
			expectSuccess: true,
		},
		{
			name:          "ok with string data",
			data:          "success message",
			requestID:     "req-456",
			expectSuccess: true,
		},
		{
			name:          "ok with nil data",
			data:          nil,
			requestID:     "req-789",
			expectSuccess: true,
		},
		{
			name:          "ok without request id",
			data:          map[string]string{"status": "ok"},
			requestID:     "",
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newEchoContext()
			if tt.requestID != "" {
				c.Set("request_id", tt.requestID)
			}

			err := OK(c, tt.data)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)

			var respBody Response
			parseErr := json.Unmarshal(rec.Body.Bytes(), &respBody)
			assert.NoError(t, parseErr)

			assert.Equal(t, tt.expectSuccess, respBody.Success)
			assert.Equal(t, tt.requestID, respBody.Meta.RequestID)
			assert.NotEmpty(t, respBody.Meta.Timestamp)
		})
	}
}

func TestOKDataInResponse(t *testing.T) {
	c, rec := newEchoContext()
	testData := map[string]interface{}{"id": 42, "name": "John"}

	err := OK(c, testData)

	assert.NoError(t, err)

	var respBody Response
	json.Unmarshal(rec.Body.Bytes(), &respBody)

	assert.Equal(t, true, respBody.Success)
	assert.NotNil(t, respBody.Data)

	// Assert data is present by checking raw JSON
	assert.Contains(t, rec.Body.String(), "\"id\":42")
	assert.Contains(t, rec.Body.String(), "\"name\":\"John\"")
}

func TestFail(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		code           string
		message        string
		requestID      string
		expectSuccess  bool
		expectMessage  string
	}{
		{
			name:          "fail with 404",
			status:        http.StatusNotFound,
			code:          "ERR-USER-404",
			message:       "user not found",
			requestID:     "req-123",
			expectSuccess: false,
			expectMessage: "user not found",
		},
		{
			name:          "fail with 500",
			status:        http.StatusInternalServerError,
			code:          "ERR-SERVER-500",
			message:       "internal server error",
			requestID:     "req-456",
			expectSuccess: false,
			expectMessage: "internal server error",
		},
		{
			name:          "fail with 400",
			status:        http.StatusBadRequest,
			code:          "ERR-VALIDATION-400",
			message:       "invalid request body",
			requestID:     "req-789",
			expectSuccess: false,
			expectMessage: "invalid request body",
		},
		{
			name:          "fail without request id",
			status:        http.StatusUnauthorized,
			code:          "ERR-UNAUTHORIZED-401",
			message:       "unauthorized access",
			requestID:     "",
			expectSuccess: false,
			expectMessage: "unauthorized access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newEchoContext()
			if tt.requestID != "" {
				c.Set("request_id", tt.requestID)
			}

			err := Fail(c, tt.status, tt.code, tt.message)

			assert.NoError(t, err)
			assert.Equal(t, tt.status, rec.Code)

			var respBody Response
			parseErr := json.Unmarshal(rec.Body.Bytes(), &respBody)
			assert.NoError(t, parseErr)

			assert.Equal(t, tt.expectSuccess, respBody.Success)
			assert.Equal(t, tt.expectMessage, respBody.Message)
			assert.Equal(t, tt.requestID, respBody.Meta.RequestID)
			assert.NotEmpty(t, respBody.Meta.Timestamp)
		})
	}
}

func TestFailNoData(t *testing.T) {
	c, rec := newEchoContext()

	err := Fail(c, http.StatusNotFound, "ERR-NOT-FOUND", "resource not found")

	assert.NoError(t, err)

	var respBody Response
	json.Unmarshal(rec.Body.Bytes(), &respBody)

	assert.Equal(t, false, respBody.Success)
	assert.Nil(t, respBody.Data)
	assert.Equal(t, "resource not found", respBody.Message)
}
