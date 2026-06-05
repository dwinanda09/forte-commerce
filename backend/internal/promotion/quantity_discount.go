package promotion

import (
	"github.com/dwinanda09/forte-commerce/internal/domain"
)

type QuantityDiscountPromotion struct{}

func (p *QuantityDiscountPromotion) Name() string {
	return "Alexa Speaker 10% Discount"
}

func (p *QuantityDiscountPromotion) Apply(cart map[string]*domain.CartItem) []domain.AppliedPromotion {
	const alexaSKU = "A304SD"

	alexa, hasAlexa := cart[alexaSKU]
	if !hasAlexa || alexa.Qty < 3 {
		return nil
	}

	// 10% off all Alexa units when >= 3
	discount := alexa.Price * float64(alexa.Qty) * 0.1

	return []domain.AppliedPromotion{
		{
			Name:        p.Name(),
			Description: "10% discount on all Alexa Speaker when buying 3 or more",
			Discount:    discount,
		},
	}
}
