package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CheckoutStatus string

const (
	CheckoutPending   CheckoutStatus = "pending"
	CheckoutCompleted CheckoutStatus = "completed"
	CheckoutExpired   CheckoutStatus = "expired"
	CheckoutFailed    CheckoutStatus = "failed"
)

type CheckoutItem struct {
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
	Total float64 `json:"total"`
}

type CheckoutResult struct {
	Items             []CheckoutItem     `json:"items"`
	PromotionsApplied []AppliedPromotion `json:"promotions_applied"`
	Subtotal          float64            `json:"subtotal"`
	TotalDiscount     float64            `json:"total_discount"`
	Total             float64            `json:"total"`
}

type CheckoutSession struct {
	ID           uuid.UUID       `db:"id" json:"id"`
	Status       CheckoutStatus  `db:"status" json:"status"`
	Items        []byte          `db:"items" json:"-"`
	Result       []byte          `db:"result" json:"-"`
	ErrorMessage *string         `db:"error_message" json:"error_message,omitempty"`
	ExpiresAt    time.Time       `db:"expires_at" json:"expires_at"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

type CheckoutRepository interface {
	Create(ctx context.Context, session *CheckoutSession) error
	FindByID(ctx context.Context, id uuid.UUID) (*CheckoutSession, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status CheckoutStatus, result []byte, errMsg *string) error
	FindExpiredPending(ctx context.Context) ([]CheckoutSession, error)
}
