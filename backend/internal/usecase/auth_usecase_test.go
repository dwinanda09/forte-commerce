package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users map[string]*domain.User
	err   error
}

func (m *mockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[username]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.Username] = user
	return nil
}

func TestAuthUsecaseLogin(t *testing.T) {
	logger := util.NewLogger()
	jwtSecret := "test-secret"

	// Create a test password hash
	password := "correctpassword"
	salt := "testsalt"

	// Simulate the hash: sha256(password + salt) then bcrypt
	hashedInput := sha256.Sum256([]byte(password + salt))
	hashedInputStr := fmt.Sprintf("%x", hashedInput)

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(hashedInputStr), bcrypt.DefaultCost)
	require.NoError(t, err)

	userID := uuid.New()
	testUser := &domain.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: string(bcryptHash),
		PasswordSalt: salt,
	}

	tests := []struct {
		name          string
		username      string
		password      string
		mockUsers     map[string]*domain.User
		mockErr       error
		expectToken   bool
		expectError   bool
		expectCode    string
	}{
		{
			name:        "valid credentials",
			username:    "testuser",
			password:    password,
			mockUsers:   map[string]*domain.User{"testuser": testUser},
			expectToken: true,
			expectError: false,
		},
		{
			name:        "user not found",
			username:    "nonexistent",
			password:    password,
			mockUsers:   map[string]*domain.User{},
			expectToken: false,
			expectError: true,
			expectCode:  "ERR-UC-038",
		},
		{
			name:        "wrong password",
			username:    "testuser",
			password:    "wrongpassword",
			mockUsers:   map[string]*domain.User{"testuser": testUser},
			expectToken: false,
			expectError: true,
			expectCode:  "ERR-UC-039",
		},
		{
			name:        "repository error",
			username:    "testuser",
			password:    password,
			mockErr:     fmt.Errorf("db error"),
			expectToken: false,
			expectError: true,
			expectCode:  "ERR-UC-038",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{
				users: tt.mockUsers,
				err:   tt.mockErr,
			}

			uc := NewAuthUsecase(mockRepo, jwtSecret, logger)
			token, err := uc.Login(context.Background(), tt.username, tt.password)

			if tt.expectToken {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			} else {
				assert.Error(t, err)
				assert.Empty(t, token)
				if tt.expectCode != "" {
					appErr, ok := util.IsAppError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.expectCode, appErr.Code)
				}
			}
		})
	}
}
