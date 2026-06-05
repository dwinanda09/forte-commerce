package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
	logger    *util.Logger
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	jwtSecret string,
	logger *util.Logger,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (string, error) {
	start := uc.logger.Start(ctx, "AuthUsecase.Login")
	defer func() { uc.logger.Finish(ctx, "AuthUsecase.Login", start, nil) }()

	// Find user
	user, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		uc.logger.Finish(ctx, "AuthUsecase.Login", start, err)
		return "", util.Wrap("ERR-UC-038", "Invalid credentials", err)
	}

	// Verify password
	// Hash: sha256(password + salt), then bcrypt compare
	hashedInput := sha256.Sum256([]byte(password + user.PasswordSalt))
	hashedInputStr := fmt.Sprintf("%x", hashedInput)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(hashedInputStr))
	if err != nil {
		uc.logger.Finish(ctx, "AuthUsecase.Login", start, err)
		return "", util.Wrap("ERR-UC-039", "Invalid credentials", err)
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		uc.logger.Finish(ctx, "AuthUsecase.Login", start, err)
		return "", util.Wrap("ERR-UC-040", "Failed to generate token", err)
	}

	return tokenString, nil
}
