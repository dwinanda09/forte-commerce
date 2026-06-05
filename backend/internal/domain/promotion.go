package domain

import "context"

type CartItem struct {
	SKU   string
	Name  string
	Price float64
	Qty   int
}

type AppliedPromotion struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Discount    float64 `json:"discount"`
}

type Promotion interface {
	Name() string
	Apply(cart map[string]*CartItem) []AppliedPromotion
}

// PromotionEngine is the interface for the promotion engine used by the checkout usecase.
// Implementations: promotion.StaticEngine (hardcoded), promotion.DynamicEngine (DB-driven).
type PromotionEngine interface {
	Apply(ctx context.Context, cart map[string]*CartItem) ([]AppliedPromotion, float64, error)
}
