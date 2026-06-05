package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/dwinanda09/forte-commerce/internal/domain"
)

func TestCampaignParsedConditions(t *testing.T) {
	tests := []struct {
		name           string
		conditions     json.RawMessage
		expectedCount  int
		expectedError  bool
		expectNotNil   bool
		validateFunc   func(t *testing.T, conds []domain.Condition)
	}{
		{
			name:          "valid single condition - cart_total_gte",
			conditions:    json.RawMessage(`[{"type":"cart_total_gte","amount":100.0}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, conds []domain.Condition) {
				assert.Equal(t, domain.CondCartTotalGTE, conds[0].Type)
				assert.Equal(t, 100.0, conds[0].Amount)
			},
		},
		{
			name:          "valid single condition - cart_has_sku",
			conditions:    json.RawMessage(`[{"type":"cart_has_sku","sku":"PRODUCT-001"}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, conds []domain.Condition) {
				assert.Equal(t, domain.CondCartHasSKU, conds[0].Type)
				assert.Equal(t, "PRODUCT-001", conds[0].SKU)
			},
		},
		{
			name:          "valid single condition - item_qty_gte",
			conditions:    json.RawMessage(`[{"type":"item_qty_gte","sku":"PRODUCT-001","min_qty":5}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, conds []domain.Condition) {
				assert.Equal(t, domain.CondItemQtyGTE, conds[0].Type)
				assert.Equal(t, "PRODUCT-001", conds[0].SKU)
				assert.Equal(t, 5, conds[0].MinQty)
			},
		},
		{
			name:          "valid single condition - cart_item_count_gte",
			conditions:    json.RawMessage(`[{"type":"cart_item_count_gte","count":3}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, conds []domain.Condition) {
				assert.Equal(t, domain.CondCartItemCountGTE, conds[0].Type)
				assert.Equal(t, 3, conds[0].Count)
			},
		},
		{
			name:          "valid multiple conditions",
			conditions:    json.RawMessage(`[{"type":"cart_total_gte","amount":100.0},{"type":"cart_has_sku","sku":"PRODUCT-001"}]`),
			expectedCount: 2,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, conds []domain.Condition) {
				assert.Equal(t, domain.CondCartTotalGTE, conds[0].Type)
				assert.Equal(t, domain.CondCartHasSKU, conds[1].Type)
			},
		},
		{
			name:          "empty array",
			conditions:    json.RawMessage(`[]`),
			expectedCount: 0,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, conds []domain.Condition) {
				// Empty slice should be present
				assert.NotNil(t, conds)
				assert.Equal(t, 0, len(conds))
			},
		},
		{
			name:          "null value",
			conditions:    json.RawMessage(`null`),
			expectedCount: 0,
			expectedError: false,
			expectNotNil:  false,
			validateFunc: nil,
		},
		{
			name:          "invalid JSON",
			conditions:    json.RawMessage(`{invalid json}`),
			expectedCount: 0,
			expectedError: true,
			expectNotNil:  false,
			validateFunc:  nil,
		},
		{
			name:          "malformed condition object",
			conditions:    json.RawMessage(`[{"type":"unknown_type"}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, conds []domain.Condition) {
				// Should still parse even with unknown type
				assert.Equal(t, domain.ConditionType("unknown_type"), conds[0].Type)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := &domain.Campaign{
				ID:         uuid.New(),
				Name:       "Test Campaign",
				Conditions: tt.conditions,
			}

			conds, err := campaign.ParsedConditions()

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, conds)
			} else {
				assert.NoError(t, err)
				if tt.expectNotNil {
					assert.NotNil(t, conds)
					assert.Equal(t, tt.expectedCount, len(conds))
					if tt.validateFunc != nil {
						tt.validateFunc(t, conds)
					}
				} else {
					assert.Nil(t, conds)
				}
			}
		})
	}
}

func TestCampaignParsedActions(t *testing.T) {
	tests := []struct {
		name          string
		actions       json.RawMessage
		expectedCount int
		expectedError bool
		expectNotNil  bool
		validateFunc  func(t *testing.T, acts []domain.Action)
	}{
		{
			name:          "valid single action - free_item",
			actions:       json.RawMessage(`[{"type":"free_item","trigger_sku":"PRODUCT-001","sku":"PRODUCT-002"}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				assert.Equal(t, domain.ActionFreeItem, acts[0].Type)
				assert.Equal(t, "PRODUCT-001", acts[0].TriggerSKU)
				assert.Equal(t, "PRODUCT-002", acts[0].SKU)
			},
		},
		{
			name:          "valid single action - buy_n_get_m",
			actions:       json.RawMessage(`[{"type":"buy_n_get_m","trigger_sku":"PRODUCT-001","sku":"PRODUCT-001","buy_n":2,"pay_m":1}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				assert.Equal(t, domain.ActionBuyNGetM, acts[0].Type)
				assert.Equal(t, 2, acts[0].BuyN)
				assert.Equal(t, 1, acts[0].PayM)
			},
		},
		{
			name:          "valid single action - pct_discount_on_sku",
			actions:       json.RawMessage(`[{"type":"pct_discount_on_sku","sku":"PRODUCT-001","pct":10.0}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				assert.Equal(t, domain.ActionPctDiscountOnSKU, acts[0].Type)
				assert.Equal(t, "PRODUCT-001", acts[0].SKU)
				assert.Equal(t, 10.0, acts[0].Pct)
			},
		},
		{
			name:          "valid single action - pct_discount_on_cart",
			actions:       json.RawMessage(`[{"type":"pct_discount_on_cart","pct":5.0}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				assert.Equal(t, domain.ActionPctDiscountOnCart, acts[0].Type)
				assert.Equal(t, 5.0, acts[0].Pct)
			},
		},
		{
			name:          "valid single action - fixed_discount",
			actions:       json.RawMessage(`[{"type":"fixed_discount","amount":25.0}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				assert.Equal(t, domain.ActionFixedDiscount, acts[0].Type)
				assert.Equal(t, 25.0, acts[0].Amount)
			},
		},
		{
			name:          "valid multiple actions",
			actions:       json.RawMessage(`[{"type":"pct_discount_on_cart","pct":5.0},{"type":"fixed_discount","amount":10.0}]`),
			expectedCount: 2,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				assert.Equal(t, domain.ActionPctDiscountOnCart, acts[0].Type)
				assert.Equal(t, domain.ActionFixedDiscount, acts[1].Type)
			},
		},
		{
			name:          "empty array",
			actions:       json.RawMessage(`[]`),
			expectedCount: 0,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				assert.NotNil(t, acts)
				assert.Equal(t, 0, len(acts))
			},
		},
		{
			name:          "null value",
			actions:       json.RawMessage(`null`),
			expectedCount: 0,
			expectedError: false,
			expectNotNil:  false,
			validateFunc:  nil,
		},
		{
			name:          "invalid JSON",
			actions:       json.RawMessage(`{invalid json}`),
			expectedCount: 0,
			expectedError: true,
			expectNotNil:  false,
			validateFunc:  nil,
		},
		{
			name:          "malformed action object",
			actions:       json.RawMessage(`[{"type":"unknown_action_type"}]`),
			expectedCount: 1,
			expectedError: false,
			expectNotNil:  true,
			validateFunc: func(t *testing.T, acts []domain.Action) {
				// Should still parse even with unknown type
				assert.Equal(t, domain.ActionType("unknown_action_type"), acts[0].Type)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := &domain.Campaign{
				ID:      uuid.New(),
				Name:    "Test Campaign",
				Actions: tt.actions,
			}

			acts, err := campaign.ParsedActions()

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, acts)
			} else {
				assert.NoError(t, err)
				if tt.expectNotNil {
					assert.NotNil(t, acts)
					assert.Equal(t, tt.expectedCount, len(acts))
					if tt.validateFunc != nil {
						tt.validateFunc(t, acts)
					}
				} else {
					assert.Nil(t, acts)
				}
			}
		})
	}
}

func TestCampaignWithBothConditionsAndActions(t *testing.T) {
	conditionsJSON := json.RawMessage(`[{"type":"cart_total_gte","amount":100.0}]`)
	actionsJSON := json.RawMessage(`[{"type":"pct_discount_on_cart","pct":10.0}]`)

	campaign := &domain.Campaign{
		ID:          uuid.New(),
		Name:        "10% off orders over $100",
		Description: "Get 10% discount on cart total when order exceeds $100",
		IsActive:    true,
		Priority:    1,
		Conditions:  conditionsJSON,
		Actions:     actionsJSON,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test conditions
	conds, err := campaign.ParsedConditions()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(conds))
	assert.Equal(t, domain.CondCartTotalGTE, conds[0].Type)
	assert.Equal(t, 100.0, conds[0].Amount)

	// Test actions
	acts, err := campaign.ParsedActions()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(acts))
	assert.Equal(t, domain.ActionPctDiscountOnCart, acts[0].Type)
	assert.Equal(t, 10.0, acts[0].Pct)
}

func TestCampaignComplexConditionsAndActions(t *testing.T) {
	// Complex multi-condition, multi-action campaign
	conditionsJSON := json.RawMessage(`[
		{"type":"cart_total_gte","amount":50.0},
		{"type":"cart_has_sku","sku":"PREMIUM-001"},
		{"type":"cart_item_count_gte","count":2}
	]`)

	actionsJSON := json.RawMessage(`[
		{"type":"pct_discount_on_sku","sku":"PREMIUM-001","pct":15.0},
		{"type":"fixed_discount","amount":5.0}
	]`)

	campaign := &domain.Campaign{
		ID:         uuid.New(),
		Conditions: conditionsJSON,
		Actions:    actionsJSON,
	}

	// Test parsing multiple conditions
	conds, err := campaign.ParsedConditions()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(conds))
	assert.Equal(t, domain.CondCartTotalGTE, conds[0].Type)
	assert.Equal(t, domain.CondCartHasSKU, conds[1].Type)
	assert.Equal(t, domain.CondCartItemCountGTE, conds[2].Type)

	// Test parsing multiple actions
	acts, err := campaign.ParsedActions()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(acts))
	assert.Equal(t, domain.ActionPctDiscountOnSKU, acts[0].Type)
	assert.Equal(t, domain.ActionFixedDiscount, acts[1].Type)
}
