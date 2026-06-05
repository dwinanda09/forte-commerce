package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CheckoutUsecase struct {
	productRepo  domain.ProductRepository
	checkoutRepo domain.CheckoutRepository
	orderRepo    domain.OrderRepository
	userRepo     domain.UserRepository
	promoEngine  domain.PromotionEngine
	queue        domain.QueuePublisher
	txManager    domain.Transactor
	logger       *util.Logger
}

func NewCheckoutUsecase(
	productRepo domain.ProductRepository,
	checkoutRepo domain.CheckoutRepository,
	orderRepo domain.OrderRepository,
	userRepo domain.UserRepository,
	promoEngine domain.PromotionEngine,
	q domain.QueuePublisher,
	txManager domain.Transactor,
	logger *util.Logger,
) *CheckoutUsecase {
	return &CheckoutUsecase{
		productRepo:  productRepo,
		checkoutRepo: checkoutRepo,
		orderRepo:    orderRepo,
		userRepo:     userRepo,
		promoEngine:  promoEngine,
		queue:        q,
		txManager:    txManager,
		logger:       logger,
	}
}

func (uc *CheckoutUsecase) Submit(ctx context.Context, skus []string) (uuid.UUID, error) {
	start := uc.logger.Start(ctx, "CheckoutUsecase.Submit")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.Submit", start, nil) }()

	// Parse SKUs into map with quantities
	skuQtyMap := make(map[string]int)
	for _, sku := range skus {
		skuQtyMap[sku]++
	}

	// Marshal items to JSON
	itemsJSON, err := json.Marshal(skuQtyMap)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.Submit", start, err)
		return uuid.Nil, util.Wrap("ERR-UC-001", "Failed to marshal items", err)
	}

	// Create checkout session
	session := &domain.CheckoutSession{
		Status:    domain.CheckoutPending,
		Items:     itemsJSON,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err = uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return uc.checkoutRepo.Create(ctx, session)
	})
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.Submit", start, err)
		return uuid.Nil, util.Wrap("ERR-UC-003", "Failed to create checkout session", err)
	}

	// Publish job to queue
	job := CheckoutJob{
		CheckoutID: session.ID,
		Items:      skuQtyMap,
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.Submit", start, err)
		return uuid.Nil, util.Wrap("ERR-UC-005", "Failed to marshal job", err)
	}

	err = uc.queue.Publish(ctx, jobJSON)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.Submit", start, err)
		return uuid.Nil, util.Wrap("ERR-UC-006", "Failed to publish job", err)
	}

	return session.ID, nil
}

func (uc *CheckoutUsecase) ProcessJob(ctx context.Context, job CheckoutJob) error {
	start := uc.logger.Start(ctx, "CheckoutUsecase.ProcessJob")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.ProcessJob", start, nil) }()

	// Get unique SKUs
	skus := make([]string, 0, len(job.Items))
	for sku := range job.Items {
		skus = append(skus, sku)
	}

	// Fetch products first (outside tx, for initial validation)
	products, err := uc.productRepo.FindBySKUs(ctx, skus)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.ProcessJob", start, err)
		return util.Wrap("ERR-UC-008", "Failed to fetch products", err)
	}

	return uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		// Lock products via FOR UPDATE — access tx from context
		if tx, ok := domain.TxFromContext(ctx); ok {
			skuArray := pq.Array(skus)
			_, err := tx.ExecContext(ctx, `SELECT 1 FROM products WHERE sku = ANY($1) FOR UPDATE`, skuArray)
			if err != nil {
				return util.Wrap("ERR-UC-009", "Failed to lock products", err)
			}
		}

		// Re-fetch products after locking
		products, err = uc.productRepo.FindBySKUs(ctx, skus)
		if err != nil {
			return util.Wrap("ERR-UC-010", "Failed to fetch products", err)
		}

		// Build product map for quick lookup
		productMap := make(map[string]*domain.Product)
		for i := range products {
			productMap[products[i].SKU] = &products[i]
		}

		// Validate available inventory
		for sku, qty := range job.Items {
			product, ok := productMap[sku]
			if !ok {
				errMsg := fmt.Sprintf("Product %s not found", sku)
				_ = uc.checkoutRepo.UpdateStatus(ctx, job.CheckoutID, domain.CheckoutFailed, nil, &errMsg)
				return nil // commit with failed status
			}

			available := product.Available()
			if available < qty {
				errMsg := fmt.Sprintf("Insufficient inventory for %s", sku)
				_ = uc.checkoutRepo.UpdateStatus(ctx, job.CheckoutID, domain.CheckoutFailed, nil, &errMsg)
				return nil // commit with failed status
			}
		}

		// Build cart for promotions
		cart := make(map[string]*domain.CartItem)
		var subtotal float64

		for sku, qty := range job.Items {
			product := productMap[sku]
			itemTotal := product.Price * float64(qty)
			cart[sku] = &domain.CartItem{
				SKU:   sku,
				Name:  product.Name,
				Price: product.Price,
				Qty:   qty,
			}
			subtotal += itemTotal
		}

		// Apply promotions
		appliedPromos, totalDiscount, err := uc.promoEngine.Apply(ctx, cart)
		if err != nil {
			return util.Wrap("ERR-UC-011A", "Failed to apply promotions", err)
		}
		total := subtotal - totalDiscount

		// Increment reserved quantities
		for sku, qty := range job.Items {
			err := uc.productRepo.IncrementReserved(ctx, sku, qty)
			if err != nil {
				return util.Wrap("ERR-UC-011", "Failed to increment reserved quantity", err)
			}
		}

		// Build checkout items
		var checkoutItems []domain.CheckoutItem
		for sku, qty := range job.Items {
			product := productMap[sku]
			checkoutItems = append(checkoutItems, domain.CheckoutItem{
				SKU:   sku,
				Name:  product.Name,
				Qty:   qty,
				Price: product.Price,
				Total: product.Price * float64(qty),
			})
		}

		// Build checkout result
		result := domain.CheckoutResult{
			Items:             checkoutItems,
			PromotionsApplied: appliedPromos,
			Subtotal:          subtotal,
			TotalDiscount:     totalDiscount,
			Total:             total,
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			return util.Wrap("ERR-UC-012", "Failed to marshal result", err)
		}

		// Update checkout session
		return uc.checkoutRepo.UpdateStatus(ctx, job.CheckoutID, domain.CheckoutCompleted, resultJSON, nil)
	})
}

func (uc *CheckoutUsecase) GetSession(ctx context.Context, id uuid.UUID) (*domain.CheckoutSession, error) {
	start := uc.logger.Start(ctx, "CheckoutUsecase.GetSession")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.GetSession", start, nil) }()

	session, err := uc.checkoutRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.GetSession", start, err)
		return nil, util.Wrap("ERR-UC-015-404", "Checkout session not found", err)
	}

	return session, nil
}

func (uc *CheckoutUsecase) Confirm(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	start := uc.logger.Start(ctx, "CheckoutUsecase.Confirm")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.Confirm", start, nil) }()

	// Get checkout session
	session, err := uc.checkoutRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.Confirm", start, err)
		return nil, util.Wrap("ERR-UC-016-404", "Checkout session not found", err)
	}

	// Validate session is completed and not expired
	if session.Status != domain.CheckoutCompleted {
		err := fmt.Errorf("checkout session not completed")
		uc.logger.Finish(ctx, "CheckoutUsecase.Confirm", start, err)
		return nil, util.Wrap("ERR-UC-017", "Checkout session is not in completed status", err)
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("checkout session expired")
		uc.logger.Finish(ctx, "CheckoutUsecase.Confirm", start, err)
		return nil, util.Wrap("ERR-UC-018", "Checkout session has expired", err)
	}

	// Check for existing order
	// (This is a simplified check - in production would query by session ID)

	// Parse result
	var result domain.CheckoutResult
	err = json.Unmarshal(session.Result, &result)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.Confirm", start, err)
		return nil, util.Wrap("ERR-UC-019", "Failed to parse checkout result", err)
	}

	// Decrement inventory and reserved for each item
	itemsJSON, _ := json.Marshal(result.Items)
	promosJSON, _ := json.Marshal(result.PromotionsApplied)

	order := &domain.Order{
		CheckoutSessionID: session.ID,
		Status:            domain.OrderPending,
		Items:             itemsJSON,
		PromotionsApplied: promosJSON,
		Subtotal:          result.Subtotal,
		TotalDiscount:     result.TotalDiscount,
		Total:             result.Total,
	}

	err = uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		for _, item := range result.Items {
			// Decrement inventory
			err := uc.productRepo.DecrementInventory(ctx, item.SKU, item.Qty)
			if err != nil {
				return util.Wrap("ERR-UC-021", "Failed to decrement inventory", err)
			}

			// Decrement reserved
			err = uc.productRepo.DecrementReserved(ctx, item.SKU, item.Qty)
			if err != nil {
				return util.Wrap("ERR-UC-022", "Failed to decrement reserved quantity", err)
			}
		}
		return uc.orderRepo.Create(ctx, order)
	})
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.Confirm", start, err)
		return nil, err
	}

	return order, nil
}

func (uc *CheckoutUsecase) ReleaseExpired(ctx context.Context) error {
	start := uc.logger.Start(ctx, "CheckoutUsecase.ReleaseExpired")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.ReleaseExpired", start, nil) }()

	// Find expired pending sessions
	sessions, err := uc.checkoutRepo.FindExpiredPending(ctx)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.ReleaseExpired", start, err)
		return util.Wrap("ERR-UC-025", "Failed to find expired sessions", err)
	}

	for _, session := range sessions {
		// Parse items
		var items map[string]int
		err := json.Unmarshal(session.Items, &items)
		if err != nil {
			continue
		}

		_ = uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
			// Decrement reserved for each item
			for sku, qty := range items {
				_ = uc.productRepo.DecrementReserved(ctx, sku, qty)
			}

			// Update session status
			return uc.checkoutRepo.UpdateStatus(ctx, session.ID, domain.CheckoutExpired, nil, nil)
		})
	}

	return nil
}

func (uc *CheckoutUsecase) PayOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	start := uc.logger.Start(ctx, "CheckoutUsecase.PayOrder")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.PayOrder", start, nil) }()

	order, err := uc.orderRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.PayOrder", start, err)
		return nil, util.Wrap("ERR-UC-026-404", "Order not found", err)
	}

	err = uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return uc.orderRepo.UpdateStatus(ctx, id, domain.OrderPaid)
	})
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.PayOrder", start, err)
		return nil, util.Wrap("ERR-UC-028", "Failed to update order status", err)
	}

	order.Status = domain.OrderPaid
	return order, nil
}

func (uc *CheckoutUsecase) CancelOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	start := uc.logger.Start(ctx, "CheckoutUsecase.CancelOrder")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.CancelOrder", start, nil) }()

	order, err := uc.orderRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.CancelOrder", start, err)
		return nil, util.Wrap("ERR-UC-030-404", "Order not found", err)
	}

	// Parse items
	var items []domain.CheckoutItem
	err = json.Unmarshal(order.Items, &items)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.CancelOrder", start, err)
		return nil, util.Wrap("ERR-UC-031", "Failed to parse order items", err)
	}

	err = uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		// Restore inventory for each item
		for _, item := range items {
			err := uc.productRepo.RestoreInventory(ctx, item.SKU, item.Qty)
			if err != nil {
				return util.Wrap("ERR-UC-033", "Failed to restore inventory", err)
			}
		}

		// Update order status
		return uc.orderRepo.UpdateStatus(ctx, id, domain.OrderCancelled)
	})
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.CancelOrder", start, err)
		return nil, err
	}

	order.Status = domain.OrderCancelled
	return order, nil
}

func (uc *CheckoutUsecase) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	start := uc.logger.Start(ctx, "CheckoutUsecase.GetOrder")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.GetOrder", start, nil) }()

	order, err := uc.orderRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.GetOrder", start, err)
		return nil, util.Wrap("ERR-UC-036-404", "Order not found", err)
	}

	return order, nil
}

func (uc *CheckoutUsecase) ListOrders(ctx context.Context) ([]domain.Order, error) {
	start := uc.logger.Start(ctx, "CheckoutUsecase.ListOrders")
	defer func() { uc.logger.Finish(ctx, "CheckoutUsecase.ListOrders", start, nil) }()

	orders, err := uc.orderRepo.FindAll(ctx)
	if err != nil {
		uc.logger.Finish(ctx, "CheckoutUsecase.ListOrders", start, err)
		return nil, util.Wrap("ERR-UC-037", "Failed to list orders", err)
	}

	return orders, nil
}

