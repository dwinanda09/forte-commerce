package promotion

import (
	"testing"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPromotionEngine(t *testing.T) {
	tests := []struct {
		name              string
		cart              map[string]*domain.CartItem
		expectedDiscount  float64
		expectedPromosLen int
		expectedTotal     float64
	}{
		{
			name: "MacBook Pro with Raspberry Pi - Free Item Promotion",
			cart: map[string]*domain.CartItem{
				"43N23P": {
					SKU:   "43N23P",
					Name:  "MacBook Pro",
					Price: 5399.99,
					Qty:   1,
				},
				"234234": {
					SKU:   "234234",
					Name:  "Raspberry Pi B",
					Price: 30.00,
					Qty:   1,
				},
			},
			expectedDiscount:  30.00,
			expectedPromosLen: 1,
			expectedTotal:     5399.99,
		},
		{
			name: "Google Home Bundle - 3 for 2",
			cart: map[string]*domain.CartItem{
				"120P90": {
					SKU:   "120P90",
					Name:  "Google Home",
					Price: 49.99,
					Qty:   3,
				},
			},
			expectedDiscount:  49.99,
			expectedPromosLen: 1,
			expectedTotal:     99.98,
		},
		{
			name: "Alexa Speaker Quantity Discount - 3 units",
			cart: map[string]*domain.CartItem{
				"A304SD": {
					SKU:   "A304SD",
					Name:  "Alexa Speaker",
					Price: 109.50,
					Qty:   3,
				},
			},
			expectedDiscount:  32.85,
			expectedPromosLen: 1,
			expectedTotal:     295.65,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine()
			applied, totalDiscount := engine.Apply(tt.cart)

			assert.Equal(t, tt.expectedPromosLen, len(applied))
			assert.InDelta(t, tt.expectedDiscount, totalDiscount, 0.01)

			// Calculate total
			var subtotal float64
			for _, item := range tt.cart {
				subtotal += item.Price * float64(item.Qty)
			}
			total := subtotal - totalDiscount
			assert.InDelta(t, tt.expectedTotal, total, 0.01)
		})
	}
}
