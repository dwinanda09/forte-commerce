package promotion

import (
	"github.com/dwinanda09/forte-commerce/internal/domain"
)

type BundlePromotion struct{}

func (p *BundlePromotion) Name() string {
	return "Google Home Bundle (3 for 2)"
}

func (p *BundlePromotion) Apply(cart map[string]*domain.CartItem) []domain.AppliedPromotion {
	const googleHomeSKU = "120P90"

	googleHome, hasGoogleHome := cart[googleHomeSKU]
	if !hasGoogleHome {
		return nil
	}

	// For every group of 3, one is free
	freeQty := googleHome.Qty / 3
	if freeQty == 0 {
		return nil
	}

	discount := googleHome.Price * float64(freeQty)

	return []domain.AppliedPromotion{
		{
			Name:        p.Name(),
			Description: "Buy 3 Google Home, pay for 2",
			Discount:    discount,
		},
	}
}
