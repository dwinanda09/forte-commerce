package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID                uuid.UUID   `db:"id" json:"id"`
	CheckoutSessionID uuid.UUID   `db:"checkout_session_id" json:"checkout_session_id"`
	Status            OrderStatus `db:"status" json:"status"`
	Items             []byte      `db:"items" json:"-"`
	PromotionsApplied []byte      `db:"promotions_applied" json:"-"`
	Subtotal          float64     `db:"subtotal" json:"subtotal"`
	TotalDiscount     float64     `db:"total_discount" json:"total_discount"`
	Total             float64     `db:"total" json:"total"`
	CreatedAt         time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time   `db:"updated_at" json:"updated_at"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	FindAll(ctx context.Context) ([]Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status OrderStatus) error
}
