package promotion

import (
	"context"
	"math"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
)

type DynamicEngine struct {
	campaignRepo domain.CampaignRepository
	logger       *util.Logger
}

func NewDynamicEngine(campaignRepo domain.CampaignRepository, logger *util.Logger) *DynamicEngine {
	return &DynamicEngine{campaignRepo: campaignRepo, logger: logger}
}

func (e *DynamicEngine) Apply(ctx context.Context, cart map[string]*domain.CartItem) ([]domain.AppliedPromotion, float64, error) {
	start := e.logger.Start(ctx, "DynamicEngine.Apply")
	defer func() { e.logger.Finish(ctx, "DynamicEngine.Apply", start, nil) }()

	campaigns, err := e.campaignRepo.FindActive(ctx)
	if err != nil {
		e.logger.Finish(ctx, "DynamicEngine.Apply", start, err)
		return nil, 0, util.Wrap("ERR-PE-001", "Failed to load active campaigns", err)
	}

	var cartTotal float64
	var cartCount int
	for _, item := range cart {
		cartTotal += item.Price * float64(item.Qty)
		cartCount += item.Qty
	}

	var allApplied []domain.AppliedPromotion
	var totalDiscount float64

	for _, campaign := range campaigns {
		conds, err := campaign.ParsedConditions()
		if err != nil {
			continue
		}
		acts, err := campaign.ParsedActions()
		if err != nil {
			continue
		}

		if !evaluateAll(conds, cart, cartTotal, cartCount) {
			continue
		}

		discount := computeDiscount(acts, cart, cartTotal)
		if discount <= 0 {
			continue
		}

		rounded := math.Round(discount*100) / 100
		allApplied = append(allApplied, domain.AppliedPromotion{
			Name:        campaign.Name,
			Description: campaign.Description,
			Discount:    rounded,
		})
		totalDiscount += rounded
	}

	return allApplied, math.Round(totalDiscount*100) / 100, nil
}

func evaluateAll(conds []domain.Condition, cart map[string]*domain.CartItem, cartTotal float64, cartCount int) bool {
	for _, c := range conds {
		if !evaluateOne(c, cart, cartTotal, cartCount) {
			return false
		}
	}
	return true
}

func evaluateOne(c domain.Condition, cart map[string]*domain.CartItem, cartTotal float64, cartCount int) bool {
	switch c.Type {
	case domain.CondCartHasSKU:
		item, ok := cart[c.SKU]
		return ok && item.Qty > 0
	case domain.CondItemQtyGTE:
		item, ok := cart[c.SKU]
		return ok && item.Qty >= c.MinQty
	case domain.CondCartTotalGTE:
		return cartTotal >= c.Amount
	case domain.CondCartItemCountGTE:
		return cartCount >= c.Count
	}
	return false
}

func computeDiscount(acts []domain.Action, cart map[string]*domain.CartItem, cartTotal float64) float64 {
	var total float64
	for _, a := range acts {
		total += applyAction(a, cart, cartTotal)
	}
	return total
}

func applyAction(a domain.Action, cart map[string]*domain.CartItem, cartTotal float64) float64 {
	switch a.Type {
	case domain.ActionFreeItem:
		freeItem, ok := cart[a.SKU]
		if !ok || freeItem.Qty == 0 {
			return 0
		}
		triggerQty := 0
		if a.TriggerSKU != "" {
			if trigger, ok := cart[a.TriggerSKU]; ok {
				triggerQty = trigger.Qty
			}
		}
		freeQty := min(triggerQty, freeItem.Qty)
		return float64(freeQty) * freeItem.Price

	case domain.ActionBuyNGetM:
		item, ok := cart[a.SKU]
		if !ok || a.BuyN <= 0 || a.PayM < 0 || a.BuyN <= a.PayM {
			return 0
		}
		groups := item.Qty / a.BuyN
		freeQty := groups * (a.BuyN - a.PayM)
		return float64(freeQty) * item.Price

	case domain.ActionPctDiscountOnSKU:
		item, ok := cart[a.SKU]
		if !ok {
			return 0
		}
		return item.Price * float64(item.Qty) * a.Pct / 100

	case domain.ActionPctDiscountOnCart:
		return cartTotal * a.Pct / 100

	case domain.ActionFixedDiscount:
		return a.Amount
	}
	return 0
}
