package resource

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/dwinanda09/forte-commerce/internal/config"
	"github.com/dwinanda09/forte-commerce/internal/domain"
)

func NewDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)

	return db, nil
}

// DBTransactor implements domain.Transactor interface
type DBTransactor struct {
	db *sqlx.DB
}

// NewDBTransactor creates a new transaction manager
func NewDBTransactor(db *sqlx.DB) *DBTransactor {
	return &DBTransactor{db: db}
}

// RunInTx executes the provided function within a database transaction
func (t *DBTransactor) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := domain.WithTx(ctx, tx)
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
