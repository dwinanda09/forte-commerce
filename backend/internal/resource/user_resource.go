package resource

import (
	"context"
	"database/sql"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserResource struct {
	db     *sqlx.DB
	logger *util.Logger
}

func NewUserResource(db *sqlx.DB, logger *util.Logger) *UserResource {
	return &UserResource{
		db:     db,
		logger: logger,
	}
}

func (r *UserResource) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	start := r.logger.Start(ctx, "UserResource.FindByUsername")
	defer func() { r.logger.Finish(ctx, "UserResource.FindByUsername", start, nil) }()

	var user domain.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, password_salt, role, created_at
		FROM users
		WHERE username = $1
	`, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.PasswordSalt, &user.Role, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Finish(ctx, "UserResource.FindByUsername", start, err)
			return nil, util.Wrap("ERR-RS-029-404", "User not found", err)
		}
		r.logger.Finish(ctx, "UserResource.FindByUsername", start, err)
		return nil, util.Wrap("ERR-RS-029", "Failed to find user", err)
	}

	return &user, nil
}

func (r *UserResource) Create(ctx context.Context, user *domain.User) error {
	start := r.logger.Start(ctx, "UserResource.Create")
	defer func() { r.logger.Finish(ctx, "UserResource.Create", start, nil) }()

	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, username, password_hash, password_salt, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Username, user.PasswordHash, user.PasswordSalt, user.CreatedAt)

	if err != nil {
		r.logger.Finish(ctx, "UserResource.Create", start, err)
		return util.Wrap("ERR-RS-030", "Failed to create user", err)
	}

	return nil
}
