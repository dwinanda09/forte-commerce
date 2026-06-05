package resource

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CheckoutResource struct {
	db     *sqlx.DB
	logger *util.Logger
}

func NewCheckoutResource(db *sqlx.DB, logger *util.Logger) *CheckoutResource {
	return &CheckoutResource{
		db:     db,
		logger: logger,
	}
}

// dbConn is a helper interface that both *sqlx.DB and *sql.Tx implement
type dbConnCheckout interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// conn returns either a transaction from context or the database connection
func (r *CheckoutResource) conn(ctx context.Context) dbConnCheckout {
	if tx, ok := domain.TxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *CheckoutResource) Create(ctx context.Context, session *domain.CheckoutSession) error {
	start := r.logger.Start(ctx, "CheckoutResource.Create")
	defer func() { r.logger.Finish(ctx, "CheckoutResource.Create", start, nil) }()

	session.ID = uuid.New()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		INSERT INTO checkout_sessions (id, status, items, result, error_message, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`)
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.Create", start, err)
		return util.Wrap("ERR-RS-017", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	var resultJSON any
	if session.Result != nil {
		resultJSON = string(session.Result)
	}

	err = stmt.QueryRowContext(ctx,
		session.ID, session.Status, string(session.Items), resultJSON,
		session.ErrorMessage, session.ExpiresAt, session.CreatedAt, session.UpdatedAt,
	).Scan(&session.ID)

	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.Create", start, err)
		return util.Wrap("ERR-RS-017", "Failed to create checkout session", err)
	}

	return nil
}

func (r *CheckoutResource) FindByID(ctx context.Context, id uuid.UUID) (*domain.CheckoutSession, error) {
	start := r.logger.Start(ctx, "CheckoutResource.FindByID")
	defer func() { r.logger.Finish(ctx, "CheckoutResource.FindByID", start, nil) }()

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, status, items, result, error_message, expires_at, created_at, updated_at
		FROM checkout_sessions
		WHERE id = $1
	`)
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-018", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	var session domain.CheckoutSession
	err = stmt.QueryRowContext(ctx, id).Scan(
		&session.ID, &session.Status, &session.Items, &session.Result,
		&session.ErrorMessage, &session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		r.logger.Finish(ctx, "CheckoutResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-018-404", "Checkout session not found", err)
	}
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-018", "Failed to fetch checkout session", err)
	}

	return &session, nil
}

func (r *CheckoutResource) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.CheckoutStatus, result []byte, errMsg *string) error {
	start := r.logger.Start(ctx, "CheckoutResource.UpdateStatus")
	defer func() { r.logger.Finish(ctx, "CheckoutResource.UpdateStatus", start, nil) }()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		UPDATE checkout_sessions
		SET status = $1, result = $2, error_message = $3, updated_at = $4
		WHERE id = $5
	`)
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.UpdateStatus", start, err)
		return util.Wrap("ERR-RS-019", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	var resultJSON any
	if result != nil {
		resultJSON = string(result)
	}

	queryResult, err := stmt.ExecContext(ctx, status, resultJSON, errMsg, time.Now(), id)
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.UpdateStatus", start, err)
		return util.Wrap("ERR-RS-019", "Failed to update checkout session", err)
	}

	rows, err := queryResult.RowsAffected()
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.UpdateStatus", start, err)
		return util.Wrap("ERR-RS-020", "Failed to check affected rows", err)
	}

	if rows == 0 {
		notFoundErr := fmt.Errorf("checkout session not found")
		r.logger.Finish(ctx, "CheckoutResource.UpdateStatus", start, notFoundErr)
		return util.Wrap("ERR-RS-021", "Checkout session not found", notFoundErr)
	}

	return nil
}

func (r *CheckoutResource) FindExpiredPending(ctx context.Context) ([]domain.CheckoutSession, error) {
	start := r.logger.Start(ctx, "CheckoutResource.FindExpiredPending")
	defer func() { r.logger.Finish(ctx, "CheckoutResource.FindExpiredPending", start, nil) }()

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, status, items, result, error_message, expires_at, created_at, updated_at
		FROM checkout_sessions
		WHERE status = $1 AND expires_at < $2
	`)
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.FindExpiredPending", start, err)
		return nil, util.Wrap("ERR-RS-022", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, domain.CheckoutPending, time.Now())
	if err != nil {
		r.logger.Finish(ctx, "CheckoutResource.FindExpiredPending", start, err)
		return nil, util.Wrap("ERR-RS-022", "Failed to fetch expired pending sessions", err)
	}
	defer rows.Close()

	var sessions []domain.CheckoutSession
	for rows.Next() {
		var s domain.CheckoutSession
		if err := rows.Scan(&s.ID, &s.Status, &s.Items, &s.Result, &s.ErrorMessage, &s.ExpiresAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			r.logger.Finish(ctx, "CheckoutResource.FindExpiredPending", start, err)
			return nil, util.Wrap("ERR-RS-022", "Failed to scan sessions", err)
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}
