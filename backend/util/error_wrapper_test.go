package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		message     string
		err         error
		expectCode  string
		expectMsg   string
		expectErr   error
	}{
		{
			name:       "wrap with all fields",
			code:       "ERR-USER-404",
			message:    "user not found",
			err:        errors.New("database error"),
			expectCode: "ERR-USER-404",
			expectMsg:  "user not found",
			expectErr:  errors.New("database error"),
		},
		{
			name:       "wrap with nil error",
			code:       "ERR-VALIDATION-400",
			message:    "invalid input",
			err:        nil,
			expectCode: "ERR-VALIDATION-400",
			expectMsg:  "invalid input",
			expectErr:  nil,
		},
		{
			name:       "wrap with empty message",
			code:       "ERR-UNKNOWN-500",
			message:    "",
			err:        errors.New("something failed"),
			expectCode: "ERR-UNKNOWN-500",
			expectMsg:  "",
			expectErr:  errors.New("something failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Wrap(tt.code, tt.message, tt.err)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectCode, result.Code)
			assert.Equal(t, tt.expectMsg, result.Message)
			if tt.expectErr != nil {
				assert.NotNil(t, result.Err)
				assert.Equal(t, tt.expectErr.Error(), result.Err.Error())
			} else {
				assert.Nil(t, result.Err)
			}
		})
	}
}

func TestAppErrorError(t *testing.T) {
	tests := []struct {
		name       string
		appErr     *AppError
		expectStr  string
	}{
		{
			name: "error string with wrapped error",
			appErr: &AppError{
				Code:    "ERR-USER-404",
				Message: "user not found",
				Err:     errors.New("db connection lost"),
			},
			expectStr: "ERR-USER-404: user not found (db connection lost)",
		},
		{
			name: "error string without wrapped error",
			appErr: &AppError{
				Code:    "ERR-VALIDATION-400",
				Message: "invalid email format",
				Err:     nil,
			},
			expectStr: "ERR-VALIDATION-400: invalid email format",
		},
		{
			name: "error string with empty message",
			appErr: &AppError{
				Code:    "ERR-UNKNOWN-500",
				Message: "",
				Err:     errors.New("panic recovered"),
			},
			expectStr: "ERR-UNKNOWN-500:  (panic recovered)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appErr.Error()
			assert.Equal(t, tt.expectStr, result)
		})
	}
}

func TestIsAppError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		expectErr  *AppError
		expectOk   bool
	}{
		{
			name: "is app error - true",
			err: &AppError{
				Code:    "ERR-USER-404",
				Message: "not found",
				Err:     nil,
			},
			expectOk: true,
		},
		{
			name:      "is not app error - plain error",
			err:       errors.New("plain error"),
			expectErr: nil,
			expectOk:  false,
		},
		{
			name:      "is not app error - nil",
			err:       nil,
			expectErr: nil,
			expectOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr, ok := IsAppError(tt.err)
			assert.Equal(t, tt.expectOk, ok)
			if tt.expectOk {
				assert.NotNil(t, appErr)
				assert.Equal(t, tt.err.(*AppError), appErr)
			} else {
				assert.Nil(t, appErr)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		expectTrue bool
	}{
		{
			name: "is 404 - true",
			err: &AppError{
				Code:    "ERR-USER-404",
				Message: "user not found",
				Err:     nil,
			},
			expectTrue: true,
		},
		{
			name: "is 404 - with different prefix",
			err: &AppError{
				Code:    "ERR-THING-404",
				Message: "not found",
				Err:     nil,
			},
			expectTrue: true,
		},
		{
			name: "is not 404 - different code",
			err: &AppError{
				Code:    "ERR-USER-500",
				Message: "server error",
				Err:     nil,
			},
			expectTrue: false,
		},
		{
			name: "is not 404 - plain error",
			err: errors.New("not found"),
			expectTrue: false,
		},
		{
			name:       "is not 404 - nil error",
			err:        nil,
			expectTrue: false,
		},
		{
			name: "is not 404 - code contains 404 but not suffix",
			err: &AppError{
				Code:    "ERR-404-USER",
				Message: "not found",
				Err:     nil,
			},
			expectTrue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			assert.Equal(t, tt.expectTrue, result)
		})
	}
}
