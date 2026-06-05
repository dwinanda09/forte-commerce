package promotion

import (
	"github.com/dwinanda09/forte-commerce/internal/domain"
)

type FreeItemPromotion struct{}

func (p *FreeItemPromotion) Name() string {
	return "MacBook Pro Free Raspberry Pi"
}

func (p *FreeItemPromotion) Apply(cart map[string]*domain.CartItem) []domain.AppliedPromotion {
	const macbookSKU = "43N23P"
	const raspberryPiSKU = "234234"

	macbook, hasMacbook := cart[macbookSKU]
	raspPi, hasRaspPi := cart[raspberryPiSKU]

	if !hasMacbook || !hasRaspPi {
		return nil
	}

	// For every MacBook, get 1 Raspberry Pi free (or min of both quantities)
	freeQty := macbook.Qty
	if raspPi.Qty < freeQty {
		freeQty = raspPi.Qty
	}

	discount := raspPi.Price * float64(freeQty)

	return []domain.AppliedPromotion{
		{
			Name:        p.Name(),
			Description: "Free Raspberry Pi with MacBook Pro",
			Discount:    discount,
		},
	}
}
