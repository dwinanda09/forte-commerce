package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ConditionType string
type ActionType string

const (
	CondCartHasSKU       ConditionType = "cart_has_sku"
	CondItemQtyGTE       ConditionType = "item_qty_gte"
	CondCartTotalGTE     ConditionType = "cart_total_gte"
	CondCartItemCountGTE ConditionType = "cart_item_count_gte"

	ActionFreeItem          ActionType = "free_item"
	ActionBuyNGetM          ActionType = "buy_n_get_m"
	ActionPctDiscountOnSKU  ActionType = "pct_discount_on_sku"
	ActionPctDiscountOnCart ActionType = "pct_discount_on_cart"
	ActionFixedDiscount     ActionType = "fixed_discount"
)

type Condition struct {
	Type   ConditionType `json:"type"`
	SKU    string        `json:"sku,omitempty"`
	MinQty int           `json:"min_qty,omitempty"`
	Qty    int           `json:"qty,omitempty"`
	Amount float64       `json:"amount,omitempty"`
	Count  int           `json:"count,omitempty"`
}

type Action struct {
	Type       ActionType `json:"type"`
	SKU        string     `json:"sku,omitempty"`
	TriggerSKU string     `json:"trigger_sku,omitempty"`
	BuyN       int        `json:"buy_n,omitempty"`
	PayM       int        `json:"pay_m,omitempty"`
	Pct        float64    `json:"pct,omitempty"`
	Amount     float64    `json:"amount,omitempty"`
}

type Campaign struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description"`
	IsActive    bool            `db:"is_active" json:"is_active"`
	Priority    int             `db:"priority" json:"priority"`
	Conditions  json.RawMessage `db:"conditions" json:"conditions"`
	Actions     json.RawMessage `db:"actions" json:"actions"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}

func (c *Campaign) ParsedConditions() ([]Condition, error) {
	var conds []Condition
	if err := json.Unmarshal(c.Conditions, &conds); err != nil {
		return nil, err
	}
	return conds, nil
}

func (c *Campaign) ParsedActions() ([]Action, error) {
	var acts []Action
	if err := json.Unmarshal(c.Actions, &acts); err != nil {
		return nil, err
	}
	return acts, nil
}

type CampaignRepository interface {
	FindAll(ctx context.Context) ([]Campaign, error)
	FindActive(ctx context.Context) ([]Campaign, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Campaign, error)
	Create(ctx context.Context, c *Campaign) error
	Update(ctx context.Context, c *Campaign) error
	Delete(ctx context.Context, id uuid.UUID) error
	Toggle(ctx context.Context, id uuid.UUID, active bool) error
}
