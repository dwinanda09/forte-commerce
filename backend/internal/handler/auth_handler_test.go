package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/internal/mocks"
	"github.com/dwinanda09/forte-commerce/internal/usecase"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandlerLogin(t *testing.T) {
	logger := util.NewLogger()
	jwtSecret := "test-jwt-secret"

	password := "correctpassword"
	salt := "testsalt"

	hashedInput := sha256.Sum256([]byte(password + salt))
	hashedInputStr := fmt.Sprintf("%x", hashedInput)

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(hashedInputStr), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: string(bcryptHash),
		PasswordSalt: salt,
	}

	tests := []struct {
		name           string
		requestBody    LoginRequest
		setupMocks     func(r *mocks.MockUserRepository)
		expectStatus   int
		expectHasToken bool
	}{
		{
			name:        "valid credentials",
			requestBody: LoginRequest{Username: "testuser", Password: password},
			setupMocks: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByUsername(gomock.Any(), "testuser").Return(testUser, nil)
			},
			expectStatus:   http.StatusOK,
			expectHasToken: true,
		},
		{
			name:        "wrong password",
			requestBody: LoginRequest{Username: "testuser", Password: "wrongpassword"},
			setupMocks: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByUsername(gomock.Any(), "testuser").Return(testUser, nil)
			},
			expectStatus:   http.StatusUnauthorized,
			expectHasToken: false,
		},
		{
			name:        "user not found",
			requestBody: LoginRequest{Username: "nonexistent", Password: password},
			setupMocks: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByUsername(gomock.Any(), "nonexistent").Return(nil, errors.New("user not found"))
			},
			expectStatus:   http.StatusUnauthorized,
			expectHasToken: false,
		},
		{
			name:        "empty username",
			requestBody: LoginRequest{Username: "", Password: password},
			setupMocks: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByUsername(gomock.Any(), "").Return(nil, errors.New("user not found"))
			},
			expectStatus:   http.StatusUnauthorized,
			expectHasToken: false,
		},
		{
			name:        "empty password",
			requestBody: LoginRequest{Username: "testuser", Password: ""},
			setupMocks: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByUsername(gomock.Any(), "testuser").Return(testUser, nil)
			},
			expectStatus:   http.StatusUnauthorized,
			expectHasToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockUserRepository(ctrl)
			tt.setupMocks(mockRepo)

			authUC := usecase.NewAuthUsecase(mockRepo, jwtSecret, logger)
			handler := NewAuthHandler(authUC)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e := echo.New()
			c := e.NewContext(req, httptest.NewRecorder())

			err := handler.Login(c)
			require.NoError(t, err)

			recorder := c.Response().Writer.(*httptest.ResponseRecorder)
			assert.Equal(t, tt.expectStatus, recorder.Code)

			if tt.expectHasToken {
				var resp map[string]any
				json.Unmarshal(recorder.Body.Bytes(), &resp)
				data := resp["data"].(map[string]any)
				token := data["token"].(string)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestAuthHandlerLoginInvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	authUC := usecase.NewAuthUsecase(mockRepo, "secret", util.NewLogger())
	handler := NewAuthHandler(authUC)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())

	err := handler.Login(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
