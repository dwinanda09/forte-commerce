package promotion

import (
	"context"

	"github.com/dwinanda09/forte-commerce/internal/domain"
)

type Engine struct {
	promotions []domain.Promotion
}

func NewEngine() *Engine {
	return &Engine{
		promotions: []domain.Promotion{
			&FreeItemPromotion{},
			&BundlePromotion{},
			&QuantityDiscountPromotion{},
		},
	}
}

func (e *Engine) Apply(cart map[string]*domain.CartItem) ([]domain.AppliedPromotion, float64) {
	var allApplied []domain.AppliedPromotion
	var totalDiscount float64

	for _, promo := range e.promotions {
		applied := promo.Apply(cart)
		if len(applied) > 0 {
			allApplied = append(allApplied, applied...)
			for _, a := range applied {
				totalDiscount += a.Discount
			}
		}
	}

	return allApplied, totalDiscount
}

// StaticAdapter wraps Engine to satisfy domain.PromotionEngine (adds ctx + error).
type StaticAdapter struct {
	engine *Engine
}

func NewStaticAdapter() *StaticAdapter {
	return &StaticAdapter{engine: NewEngine()}
}

func (a *StaticAdapter) Apply(_ context.Context, cart map[string]*domain.CartItem) ([]domain.AppliedPromotion, float64, error) {
	applied, discount := a.engine.Apply(cart)
	return applied, discount, nil
}
