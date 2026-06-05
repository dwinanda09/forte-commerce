package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	PasswordSalt string    `db:"password_salt"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
}
