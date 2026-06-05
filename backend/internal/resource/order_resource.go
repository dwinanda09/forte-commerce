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

type OrderResource struct {
	db     *sqlx.DB
	logger *util.Logger
}

func NewOrderResource(db *sqlx.DB, logger *util.Logger) *OrderResource {
	return &OrderResource{
		db:     db,
		logger: logger,
	}
}

// dbConn is a helper interface that both *sqlx.DB and *sql.Tx implement
type dbConnOrder interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// conn returns either a transaction from context or the database connection
func (r *OrderResource) conn(ctx context.Context) dbConnOrder {
	if tx, ok := domain.TxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *OrderResource) Create(ctx context.Context, order *domain.Order) error {
	start := r.logger.Start(ctx, "OrderResource.Create")
	defer func() { r.logger.Finish(ctx, "OrderResource.Create", start, nil) }()

	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		INSERT INTO orders (id, checkout_session_id, status, items, promotions_applied, subtotal, total_discount, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`)
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.Create", start, err)
		return util.Wrap("ERR-RS-023", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx,
		order.ID, order.CheckoutSessionID, order.Status, string(order.Items),
		string(order.PromotionsApplied), order.Subtotal, order.TotalDiscount, order.Total,
		order.CreatedAt, order.UpdatedAt,
	).Scan(&order.ID)

	if err != nil {
		r.logger.Finish(ctx, "OrderResource.Create", start, err)
		return util.Wrap("ERR-RS-023", "Failed to create order", err)
	}

	return nil
}

func (r *OrderResource) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	start := r.logger.Start(ctx, "OrderResource.FindByID")
	defer func() { r.logger.Finish(ctx, "OrderResource.FindByID", start, nil) }()

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, checkout_session_id, status, items, promotions_applied, subtotal, total_discount, total, created_at, updated_at
		FROM orders
		WHERE id = $1
	`)
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-024", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	var order domain.Order
	err = stmt.QueryRowContext(ctx, id).Scan(
		&order.ID, &order.CheckoutSessionID, &order.Status, &order.Items,
		&order.PromotionsApplied, &order.Subtotal, &order.TotalDiscount, &order.Total,
		&order.CreatedAt, &order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		r.logger.Finish(ctx, "OrderResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-024-404", "Order not found", err)
	}
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-024", "Failed to fetch order", err)
	}

	return &order, nil
}

func (r *OrderResource) FindAll(ctx context.Context) ([]domain.Order, error) {
	start := r.logger.Start(ctx, "OrderResource.FindAll")
	defer func() { r.logger.Finish(ctx, "OrderResource.FindAll", start, nil) }()

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, checkout_session_id, status, items, promotions_applied, subtotal, total_discount, total, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`)
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.FindAll", start, err)
		return nil, util.Wrap("ERR-RS-025", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.FindAll", start, err)
		return nil, util.Wrap("ERR-RS-025", "Failed to fetch orders", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.CheckoutSessionID, &o.Status, &o.Items,
			&o.PromotionsApplied, &o.Subtotal, &o.TotalDiscount, &o.Total,
			&o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			r.logger.Finish(ctx, "OrderResource.FindAll", start, err)
			return nil, util.Wrap("ERR-RS-025", "Failed to scan orders", err)
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (r *OrderResource) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	start := r.logger.Start(ctx, "OrderResource.UpdateStatus")
	defer func() { r.logger.Finish(ctx, "OrderResource.UpdateStatus", start, nil) }()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		UPDATE orders
		SET status = $1, updated_at = $2
		WHERE id = $3
	`)
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.UpdateStatus", start, err)
		return util.Wrap("ERR-RS-026", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	queryResult, err := stmt.ExecContext(ctx, status, time.Now(), id)
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.UpdateStatus", start, err)
		return util.Wrap("ERR-RS-026", "Failed to update order", err)
	}

	rows, err := queryResult.RowsAffected()
	if err != nil {
		r.logger.Finish(ctx, "OrderResource.UpdateStatus", start, err)
		return util.Wrap("ERR-RS-027", "Failed to check affected rows", err)
	}

	if rows == 0 {
		notFoundErr := fmt.Errorf("order not found")
		r.logger.Finish(ctx, "OrderResource.UpdateStatus", start, notFoundErr)
		return util.Wrap("ERR-RS-028", "Order not found", notFoundErr)
	}

	return nil
}
