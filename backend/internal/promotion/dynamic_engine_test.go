package promotion

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/internal/mocks"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func makeConditions(conds []domain.Condition) json.RawMessage {
	b, _ := json.Marshal(conds)
	return b
}

func makeActions(acts []domain.Action) json.RawMessage {
	b, _ := json.Marshal(acts)
	return b
}

func TestDynamicEngine_Apply(t *testing.T) {
	logger := util.NewLogger()

	tests := []struct {
		name             string
		campaigns        []domain.Campaign
		cart             map[string]*domain.CartItem
		wantDiscountName string
		wantDiscount     float64
		wantCount        int
	}{
		{
			name: "free item — macbook + raspberry pi",
			campaigns: []domain.Campaign{
				{
					ID:          uuid.New(),
					Name:        "MacBook Pro Free Raspberry Pi",
					Description: "Free Raspberry Pi with every MacBook Pro",
					IsActive:    true,
					Conditions: makeConditions([]domain.Condition{
						{Type: domain.CondCartHasSKU, SKU: "43N23P"},
						{Type: domain.CondCartHasSKU, SKU: "234234"},
					}),
					Actions: makeActions([]domain.Action{
						{Type: domain.ActionFreeItem, SKU: "234234", TriggerSKU: "43N23P"},
					}),
				},
			},
			cart: map[string]*domain.CartItem{
				"43N23P": {SKU: "43N23P", Name: "MacBook Pro", Price: 5399.99, Qty: 1},
				"234234": {SKU: "234234", Name: "Raspberry Pi B", Price: 30.00, Qty: 1},
			},
			wantDiscountName: "MacBook Pro Free Raspberry Pi",
			wantDiscount:     30.00,
			wantCount:        1,
		},
		{
			name: "buy 3 get 2 — google home",
			campaigns: []domain.Campaign{
				{
					ID:          uuid.New(),
					Name:        "Google Home Bundle (3 for 2)",
					Description: "Buy 3 Google Home, pay for 2",
					IsActive:    true,
					Conditions: makeConditions([]domain.Condition{
						{Type: domain.CondItemQtyGTE, SKU: "120P90", MinQty: 3},
					}),
					Actions: makeActions([]domain.Action{
						{Type: domain.ActionBuyNGetM, SKU: "120P90", BuyN: 3, PayM: 2},
					}),
				},
			},
			cart: map[string]*domain.CartItem{
				"120P90": {SKU: "120P90", Name: "Google Home", Price: 49.99, Qty: 3},
			},
			wantDiscountName: "Google Home Bundle (3 for 2)",
			wantDiscount:     49.99,
			wantCount:        1,
		},
		{
			name: "10% discount — alexa 3+",
			campaigns: []domain.Campaign{
				{
					ID:          uuid.New(),
					Name:        "Alexa Speaker 10% Discount",
					Description: "10% off Alexa Speaker when buying 3 or more",
					IsActive:    true,
					Conditions: makeConditions([]domain.Condition{
						{Type: domain.CondItemQtyGTE, SKU: "A304SD", MinQty: 3},
					}),
					Actions: makeActions([]domain.Action{
						{Type: domain.ActionPctDiscountOnSKU, SKU: "A304SD", Pct: 10},
					}),
				},
			},
			cart: map[string]*domain.CartItem{
				"A304SD": {SKU: "A304SD", Name: "Alexa Speaker", Price: 109.50, Qty: 3},
			},
			wantDiscountName: "Alexa Speaker 10% Discount",
			wantDiscount:     32.85,
			wantCount:        1,
		},
		{
			name: "condition not met — no discount",
			campaigns: []domain.Campaign{
				{
					ID:          uuid.New(),
					Name:        "Google Home Bundle (3 for 2)",
					IsActive:    true,
					Conditions: makeConditions([]domain.Condition{
						{Type: domain.CondItemQtyGTE, SKU: "120P90", MinQty: 3},
					}),
					Actions: makeActions([]domain.Action{
						{Type: domain.ActionBuyNGetM, SKU: "120P90", BuyN: 3, PayM: 2},
					}),
				},
			},
			cart: map[string]*domain.CartItem{
				"120P90": {SKU: "120P90", Name: "Google Home", Price: 49.99, Qty: 2},
			},
			wantCount:    0,
			wantDiscount: 0,
		},
		{
			name:         "empty campaigns — no discount",
			campaigns:    []domain.Campaign{},
			cart:         map[string]*domain.CartItem{"43N23P": {SKU: "43N23P", Price: 5399.99, Qty: 1}},
			wantCount:    0,
			wantDiscount: 0,
		},
		{
			name: "cart total GTE condition",
			campaigns: []domain.Campaign{
				{
					ID:          uuid.New(),
					Name:        "Big Spender 5% Off",
					IsActive:    true,
					Conditions: makeConditions([]domain.Condition{
						{Type: domain.CondCartTotalGTE, Amount: 1000},
					}),
					Actions: makeActions([]domain.Action{
						{Type: domain.ActionPctDiscountOnCart, Pct: 5},
					}),
				},
			},
			cart: map[string]*domain.CartItem{
				"X": {SKU: "X", Price: 2000, Qty: 1},
			},
			wantDiscountName: "Big Spender 5% Off",
			wantDiscount:     100.00,
			wantCount:        1,
		},
		{
			name: "fixed discount",
			campaigns: []domain.Campaign{
				{
					ID:          uuid.New(),
					Name:        "$50 Off",
					IsActive:    true,
					Conditions: makeConditions([]domain.Condition{
						{Type: domain.CondCartItemCountGTE, Count: 2},
					}),
					Actions: makeActions([]domain.Action{
						{Type: domain.ActionFixedDiscount, Amount: 50},
					}),
				},
			},
			cart: map[string]*domain.CartItem{
				"A": {SKU: "A", Price: 100, Qty: 2},
			},
			wantDiscountName: "$50 Off",
			wantDiscount:     50.00,
			wantCount:        1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockCampaignRepository(ctrl)
			mockRepo.EXPECT().FindActive(gomock.Any()).Return(tt.campaigns, nil)

			engine := NewDynamicEngine(mockRepo, logger)
			applied, totalDiscount, err := engine.Apply(context.Background(), tt.cart)

			require.NoError(t, err)
			assert.Len(t, applied, tt.wantCount)
			assert.InDelta(t, tt.wantDiscount, totalDiscount, 0.01)
			if tt.wantCount > 0 && tt.wantDiscountName != "" {
				assert.Equal(t, tt.wantDiscountName, applied[0].Name)
				assert.InDelta(t, tt.wantDiscount, applied[0].Discount, 0.01)
			}
		})
	}
}

func TestDynamicEngine_Apply_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := util.NewLogger()
	mockRepo := mocks.NewMockCampaignRepository(ctrl)
	mockRepo.EXPECT().FindActive(gomock.Any()).Return(nil, assert.AnError)

	engine := NewDynamicEngine(mockRepo, logger)
	_, _, err := engine.Apply(context.Background(), map[string]*domain.CartItem{})
	assert.Error(t, err)
}
