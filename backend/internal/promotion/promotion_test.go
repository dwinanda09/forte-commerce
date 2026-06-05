package promotion

import (
	"testing"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestFreeItemPromotion(t *testing.T) {
	tests := []struct {
		name               string
		cart               map[string]*domain.CartItem
		expectedApplied    int
		expectedDiscount   float64
		expectedPromoName  string
	}{
		{
			name: "both MacBook and Raspberry Pi present",
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
			expectedApplied:   1,
			expectedDiscount:   30.00,
			expectedPromoName: "MacBook Pro Free Raspberry Pi",
		},
		{
			name: "multiple MacBooks with single Raspberry Pi",
			cart: map[string]*domain.CartItem{
				"43N23P": {
					SKU:   "43N23P",
					Name:  "MacBook Pro",
					Price: 5399.99,
					Qty:   2,
				},
				"234234": {
					SKU:   "234234",
					Name:  "Raspberry Pi B",
					Price: 30.00,
					Qty:   1,
				},
			},
			expectedApplied:   1,
			expectedDiscount:   30.00,
			expectedPromoName: "MacBook Pro Free Raspberry Pi",
		},
		{
			name: "multiple Raspberry Pis with single MacBook",
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
					Qty:   3,
				},
			},
			expectedApplied:   1,
			expectedDiscount:   30.00,
			expectedPromoName: "MacBook Pro Free Raspberry Pi",
		},
		{
			name: "multiple MacBooks with multiple Raspberry Pis",
			cart: map[string]*domain.CartItem{
				"43N23P": {
					SKU:   "43N23P",
					Name:  "MacBook Pro",
					Price: 5399.99,
					Qty:   2,
				},
				"234234": {
					SKU:   "234234",
					Name:  "Raspberry Pi B",
					Price: 30.00,
					Qty:   2,
				},
			},
			expectedApplied:   1,
			expectedDiscount:   60.00,
			expectedPromoName: "MacBook Pro Free Raspberry Pi",
		},
		{
			name: "only MacBook present",
			cart: map[string]*domain.CartItem{
				"43N23P": {
					SKU:   "43N23P",
					Name:  "MacBook Pro",
					Price: 5399.99,
					Qty:   1,
				},
			},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
		{
			name: "only Raspberry Pi present",
			cart: map[string]*domain.CartItem{
				"234234": {
					SKU:   "234234",
					Name:  "Raspberry Pi B",
					Price: 30.00,
					Qty:   1,
				},
			},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
		{
			name:              "empty cart",
			cart:              map[string]*domain.CartItem{},
			expectedApplied:   0,
			expectedDiscount:  0.00,
		},
	}

	promo := &FreeItemPromotion{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applied := promo.Apply(tt.cart)

			assert.Equal(t, tt.expectedApplied, len(applied))
			if len(applied) > 0 {
				assert.Equal(t, tt.expectedPromoName, applied[0].Name)
				assert.InDelta(t, tt.expectedDiscount, applied[0].Discount, 0.01)
			}
		})
	}
}

func TestBundlePromotion(t *testing.T) {
	tests := []struct {
		name              string
		cart              map[string]*domain.CartItem
		expectedApplied   int
		expectedDiscount  float64
		expectedPromoName string
	}{
		{
			name: "3 Google Home units - 1 free",
			cart: map[string]*domain.CartItem{
				"120P90": {
					SKU:   "120P90",
					Name:  "Google Home",
					Price: 49.99,
					Qty:   3,
				},
			},
			expectedApplied:   1,
			expectedDiscount:  49.99,
			expectedPromoName: "Google Home Bundle (3 for 2)",
		},
		{
			name: "6 Google Home units - 2 free",
			cart: map[string]*domain.CartItem{
				"120P90": {
					SKU:   "120P90",
					Name:  "Google Home",
					Price: 49.99,
					Qty:   6,
				},
			},
			expectedApplied:   1,
			expectedDiscount:  99.98,
			expectedPromoName: "Google Home Bundle (3 for 2)",
		},
		{
			name: "9 Google Home units - 3 free",
			cart: map[string]*domain.CartItem{
				"120P90": {
					SKU:   "120P90",
					Name:  "Google Home",
					Price: 49.99,
					Qty:   9,
				},
			},
			expectedApplied:   1,
			expectedDiscount:  149.97,
			expectedPromoName: "Google Home Bundle (3 for 2)",
		},
		{
			name: "2 Google Home units - no discount",
			cart: map[string]*domain.CartItem{
				"120P90": {
					SKU:   "120P90",
					Name:  "Google Home",
					Price: 49.99,
					Qty:   2,
				},
			},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
		{
			name: "1 Google Home unit - no discount",
			cart: map[string]*domain.CartItem{
				"120P90": {
					SKU:   "120P90",
					Name:  "Google Home",
					Price: 49.99,
					Qty:   1,
				},
			},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
		{
			name:             "Google Home not in cart",
			cart:             map[string]*domain.CartItem{},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
	}

	promo := &BundlePromotion{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applied := promo.Apply(tt.cart)

			assert.Equal(t, tt.expectedApplied, len(applied))
			if len(applied) > 0 {
				assert.Equal(t, tt.expectedPromoName, applied[0].Name)
				assert.InDelta(t, tt.expectedDiscount, applied[0].Discount, 0.01)
			}
		})
	}
}

func TestQuantityDiscountPromotion(t *testing.T) {
	tests := []struct {
		name              string
		cart              map[string]*domain.CartItem
		expectedApplied   int
		expectedDiscount  float64
		expectedPromoName string
	}{
		{
			name: "3 Alexa speakers - 10% discount",
			cart: map[string]*domain.CartItem{
				"A304SD": {
					SKU:   "A304SD",
					Name:  "Alexa Speaker",
					Price: 109.50,
					Qty:   3,
				},
			},
			expectedApplied:   1,
			expectedDiscount:  32.85,
			expectedPromoName: "Alexa Speaker 10% Discount",
		},
		{
			name: "5 Alexa speakers - 10% discount",
			cart: map[string]*domain.CartItem{
				"A304SD": {
					SKU:   "A304SD",
					Name:  "Alexa Speaker",
					Price: 100.00,
					Qty:   5,
				},
			},
			expectedApplied:   1,
			expectedDiscount:  50.00,
			expectedPromoName: "Alexa Speaker 10% Discount",
		},
		{
			name: "10 Alexa speakers - 10% discount",
			cart: map[string]*domain.CartItem{
				"A304SD": {
					SKU:   "A304SD",
					Name:  "Alexa Speaker",
					Price: 50.00,
					Qty:   10,
				},
			},
			expectedApplied:   1,
			expectedDiscount:  50.00,
			expectedPromoName: "Alexa Speaker 10% Discount",
		},
		{
			name: "2 Alexa speakers - no discount",
			cart: map[string]*domain.CartItem{
				"A304SD": {
					SKU:   "A304SD",
					Name:  "Alexa Speaker",
					Price: 109.50,
					Qty:   2,
				},
			},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
		{
			name: "1 Alexa speaker - no discount",
			cart: map[string]*domain.CartItem{
				"A304SD": {
					SKU:   "A304SD",
					Name:  "Alexa Speaker",
					Price: 109.50,
					Qty:   1,
				},
			},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
		{
			name: "0 Alexa speakers - no discount",
			cart: map[string]*domain.CartItem{
				"A304SD": {
					SKU:   "A304SD",
					Name:  "Alexa Speaker",
					Price: 109.50,
					Qty:   0,
				},
			},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
		{
			name:             "Alexa not in cart",
			cart:             map[string]*domain.CartItem{},
			expectedApplied:  0,
			expectedDiscount: 0.00,
		},
	}

	promo := &QuantityDiscountPromotion{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applied := promo.Apply(tt.cart)

			assert.Equal(t, tt.expectedApplied, len(applied))
			if len(applied) > 0 {
				assert.Equal(t, tt.expectedPromoName, applied[0].Name)
				assert.InDelta(t, tt.expectedDiscount, applied[0].Discount, 0.01)
			}
		})
	}
}
