package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/internal/mocks"
	"github.com/dwinanda09/forte-commerce/internal/promotion"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func txPassthrough(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestCheckoutUsecaseGetSession(t *testing.T) {
	sessionID := uuid.New()

	tests := []struct {
		name            string
		checkoutID      uuid.UUID
		setupMocks      func(*mocks.MockCheckoutRepository)
		expectSession   bool
		expectErrorCode string
	}{
		{
			name:       "session found",
			checkoutID: sessionID,
			setupMocks: func(m *mocks.MockCheckoutRepository) {
				m.EXPECT().FindByID(gomock.Any(), sessionID).Return(&domain.CheckoutSession{
					ID:        sessionID,
					Status:    domain.CheckoutPending,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}, nil)
			},
			expectSession: true,
		},
		{
			name:       "session not found",
			checkoutID: uuid.New(),
			setupMocks: func(m *mocks.MockCheckoutRepository) {
				m.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("not found"))
			},
			expectErrorCode: "ERR-UC-015-404",
		},
		{
			name:       "repository error",
			checkoutID: uuid.New(),
			setupMocks: func(m *mocks.MockCheckoutRepository) {
				m.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("db error"))
			},
			expectErrorCode: "ERR-UC-015-404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockCheckout)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			session, err := uc.GetSession(context.Background(), tt.checkoutID)

			if tt.expectSession {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tt.checkoutID, session.ID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, session)
				appErr, ok := util.IsAppError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectErrorCode, appErr.Code)
			}
		})
	}
}

func TestCheckoutUsecaseGetOrder(t *testing.T) {
	orderID := uuid.New()
	checkoutID := uuid.New()

	tests := []struct {
		name        string
		orderID     uuid.UUID
		setupMocks  func(*mocks.MockOrderRepository)
		expectOrder bool
		expectCode  string
	}{
		{
			name:    "order found",
			orderID: orderID,
			setupMocks: func(m *mocks.MockOrderRepository) {
				m.EXPECT().FindByID(gomock.Any(), orderID).Return(&domain.Order{
					ID:                orderID,
					CheckoutSessionID: checkoutID,
					Status:            domain.OrderPaid,
					Items:             []byte("[]"),
					PromotionsApplied: []byte("[]"),
					Subtotal:          100.0,
					Total:             100.0,
				}, nil)
			},
			expectOrder: true,
		},
		{
			name:    "order not found",
			orderID: uuid.New(),
			setupMocks: func(m *mocks.MockOrderRepository) {
				m.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("not found"))
			},
			expectCode: "ERR-UC-036-404",
		},
		{
			name:    "repository error",
			orderID: uuid.New(),
			setupMocks: func(m *mocks.MockOrderRepository) {
				m.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("db error"))
			},
			expectCode: "ERR-UC-036-404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockOrder)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			order, err := uc.GetOrder(context.Background(), tt.orderID)

			if tt.expectOrder {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.orderID, order.ID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, order)
				appErr, ok := util.IsAppError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, appErr.Code)
			}
		})
	}
}

func TestCheckoutUsecaseListOrders(t *testing.T) {
	order1 := domain.Order{ID: uuid.New(), Status: domain.OrderPaid, Items: []byte("[]"), PromotionsApplied: []byte("[]")}
	order2 := domain.Order{ID: uuid.New(), Status: domain.OrderPending, Items: []byte("[]"), PromotionsApplied: []byte("[]")}

	tests := []struct {
		name        string
		setupMocks  func(*mocks.MockOrderRepository)
		expectCount int
		expectCode  string
	}{
		{
			name: "multiple orders",
			setupMocks: func(m *mocks.MockOrderRepository) {
				m.EXPECT().FindAll(gomock.Any()).Return([]domain.Order{order1, order2}, nil)
			},
			expectCount: 2,
		},
		{
			name: "empty order list",
			setupMocks: func(m *mocks.MockOrderRepository) {
				m.EXPECT().FindAll(gomock.Any()).Return([]domain.Order{}, nil)
			},
			expectCount: 0,
		},
		{
			name: "repository error",
			setupMocks: func(m *mocks.MockOrderRepository) {
				m.EXPECT().FindAll(gomock.Any()).Return(nil, fmt.Errorf("db error"))
			},
			expectCode: "ERR-UC-037",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockOrder)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			orders, err := uc.ListOrders(context.Background())

			if tt.expectCode != "" {
				assert.Error(t, err)
				appErr, ok := util.IsAppError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, appErr.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectCount, len(orders))
			}
		})
	}
}

func TestCheckoutUsecasePayOrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProduct := mocks.NewMockProductRepository(ctrl)
	mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
	mockOrder := mocks.NewMockOrderRepository(ctrl)
	mockUser := mocks.NewMockUserRepository(ctrl)
	mockTx := mocks.NewMockTransactor(ctrl)

	mockOrder.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("not found"))

	uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
		promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

	order, err := uc.PayOrder(context.Background(), uuid.New())

	assert.Error(t, err)
	assert.Nil(t, order)
	appErr, ok := util.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, "ERR-UC-026-404", appErr.Code)
}

func TestCheckoutUsecaseCancelOrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProduct := mocks.NewMockProductRepository(ctrl)
	mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
	mockOrder := mocks.NewMockOrderRepository(ctrl)
	mockUser := mocks.NewMockUserRepository(ctrl)
	mockTx := mocks.NewMockTransactor(ctrl)

	mockOrder.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("not found"))

	uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
		promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

	order, err := uc.CancelOrder(context.Background(), uuid.New())

	assert.Error(t, err)
	assert.Nil(t, order)
	appErr, ok := util.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, "ERR-UC-030-404", appErr.Code)
}

func TestCheckoutUsecasePayOrder(t *testing.T) {
	orderID := uuid.New()
	checkoutID := uuid.New()

	tests := []struct {
		name        string
		orderID     uuid.UUID
		setupMocks  func(*mocks.MockOrderRepository, *mocks.MockTransactor)
		expectCode  string
	}{
		{
			name:    "order not found returns 404",
			orderID: uuid.New(),
			setupMocks: func(o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				o.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("not found"))
			},
			expectCode: "ERR-UC-026-404",
		},
		{
			name:    "repository error returns 404",
			orderID: orderID,
			setupMocks: func(o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				o.EXPECT().FindByID(gomock.Any(), orderID).Return(nil, fmt.Errorf("db error"))
			},
			expectCode: "ERR-UC-026-404",
		},
		{
			name:    "success updates status to paid",
			orderID: orderID,
			setupMocks: func(o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				o.EXPECT().FindByID(gomock.Any(), orderID).Return(&domain.Order{
					ID:                orderID,
					CheckoutSessionID: checkoutID,
					Status:            domain.OrderPending,
					Items:             []byte("[]"),
					PromotionsApplied: []byte("[]"),
					Subtotal:          100.0,
					Total:             100.0,
				}, nil)
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				o.EXPECT().UpdateStatus(gomock.Any(), orderID, domain.OrderPaid).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockOrder, mockTx)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			order, err := uc.PayOrder(context.Background(), tt.orderID)

			if tt.expectCode != "" {
				assert.Error(t, err)
				assert.Nil(t, order)
				appErr, ok := util.IsAppError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, appErr.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, domain.OrderPaid, order.Status)
			}
		})
	}
}

func TestCheckoutUsecaseCancelOrder(t *testing.T) {
	orderID := uuid.New()
	checkoutID := uuid.New()
	itemsJSON := []byte(`[{"sku":"ITEM1","name":"Item 1","qty":2,"price":50.0,"total":100.0}]`)
	promosJSON := []byte("[]")

	tests := []struct {
		name       string
		orderID    uuid.UUID
		setupMocks func(*mocks.MockProductRepository, *mocks.MockOrderRepository, *mocks.MockTransactor)
		expectCode string
	}{
		{
			name:    "order not found returns 404",
			orderID: uuid.New(),
			setupMocks: func(p *mocks.MockProductRepository, o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				o.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("not found"))
			},
			expectCode: "ERR-UC-030-404",
		},
		{
			name:    "repository error on find returns 404",
			orderID: orderID,
			setupMocks: func(p *mocks.MockProductRepository, o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				o.EXPECT().FindByID(gomock.Any(), orderID).Return(nil, fmt.Errorf("db error"))
			},
			expectCode: "ERR-UC-030-404",
		},
		{
			name:    "success restores inventory and cancels order",
			orderID: orderID,
			setupMocks: func(p *mocks.MockProductRepository, o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				o.EXPECT().FindByID(gomock.Any(), orderID).Return(&domain.Order{
					ID:                orderID,
					CheckoutSessionID: checkoutID,
					Status:            domain.OrderPending,
					Items:             itemsJSON,
					PromotionsApplied: promosJSON,
					Subtotal:          100.0,
					Total:             100.0,
				}, nil)
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				p.EXPECT().RestoreInventory(gomock.Any(), "ITEM1", 2).Return(nil)
				o.EXPECT().UpdateStatus(gomock.Any(), orderID, domain.OrderCancelled).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockProduct, mockOrder, mockTx)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			order, err := uc.CancelOrder(context.Background(), tt.orderID)

			if tt.expectCode != "" {
				assert.Error(t, err)
				assert.Nil(t, order)
				appErr, ok := util.IsAppError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, appErr.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, domain.OrderCancelled, order.Status)
			}
		})
	}
}

func TestCheckoutUsecaseSubmit(t *testing.T) {
	tests := []struct {
		name        string
		skus        []string
		setupMocks  func(*mocks.MockCheckoutRepository, *mocks.MockTransactor, *mocks.MockQueuePublisher)
		expectError bool
		expectCode  string
	}{
		{
			name: "success",
			skus: []string{"SKU1", "SKU2"},
			setupMocks: func(c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor, q *mocks.MockQueuePublisher) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				c.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, s *domain.CheckoutSession) error {
						s.ID = uuid.New()
						return nil
					})
				q.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "checkout create error",
			skus: []string{"SKU1"},
			setupMocks: func(c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor, q *mocks.MockQueuePublisher) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				c.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fmt.Errorf("db error"))
			},
			expectError: true,
			expectCode:  "ERR-UC-003",
		},
		{
			name: "queue publish error",
			skus: []string{"SKU1"},
			setupMocks: func(c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor, q *mocks.MockQueuePublisher) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				c.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, s *domain.CheckoutSession) error {
						s.ID = uuid.New()
						return nil
					})
				q.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(fmt.Errorf("mq error"))
			},
			expectError: true,
			expectCode:  "ERR-UC-006",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)
			mockQueue := mocks.NewMockQueuePublisher(ctrl)

			tt.setupMocks(mockCheckout, mockTx, mockQueue)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), mockQueue, mockTx, util.NewLogger())

			id, err := uc.Submit(context.Background(), tt.skus)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uuid.Nil, id)
				if tt.expectCode != "" {
					appErr, ok := util.IsAppError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.expectCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, id)
			}
		})
	}
}

func TestCheckoutUsecaseProcessJob(t *testing.T) {
	checkoutID := uuid.New()
	product := domain.Product{
		SKU:          "SKU1",
		Name:         "Product 1",
		Price:        100.0,
		InventoryQty: 10,
		ReservedQty:  2,
	}

	tests := []struct {
		name        string
		job         CheckoutJob
		setupMocks  func(*mocks.MockProductRepository, *mocks.MockCheckoutRepository, *mocks.MockTransactor)
		expectError bool
	}{
		{
			name: "success",
			job:  CheckoutJob{CheckoutID: checkoutID, Items: map[string]int{"SKU1": 2}},
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor) {
				gomock.InOrder(
					p.EXPECT().FindBySKUs(gomock.Any(), gomock.Any()).Return([]domain.Product{product}, nil),
					tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough),
				)
				p.EXPECT().FindBySKUs(gomock.Any(), gomock.Any()).Return([]domain.Product{product}, nil)
				p.EXPECT().IncrementReserved(gomock.Any(), "SKU1", 2).Return(nil)
				c.EXPECT().UpdateStatus(gomock.Any(), checkoutID, domain.CheckoutCompleted, gomock.Any(), nil).Return(nil)
			},
		},
		{
			name: "product not found sets checkout failed",
			job:  CheckoutJob{CheckoutID: checkoutID, Items: map[string]int{"MISSING": 1}},
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor) {
				p.EXPECT().FindBySKUs(gomock.Any(), gomock.Any()).Return([]domain.Product{}, nil)
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				p.EXPECT().FindBySKUs(gomock.Any(), gomock.Any()).Return([]domain.Product{}, nil)
				c.EXPECT().UpdateStatus(gomock.Any(), checkoutID, domain.CheckoutFailed, nil, gomock.Any()).Return(nil)
			},
		},
		{
			name: "insufficient inventory sets checkout failed",
			job:  CheckoutJob{CheckoutID: checkoutID, Items: map[string]int{"SKU1": 10}},
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor) {
				p.EXPECT().FindBySKUs(gomock.Any(), gomock.Any()).Return([]domain.Product{product}, nil)
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				p.EXPECT().FindBySKUs(gomock.Any(), gomock.Any()).Return([]domain.Product{product}, nil)
				c.EXPECT().UpdateStatus(gomock.Any(), checkoutID, domain.CheckoutFailed, nil, gomock.Any()).Return(nil)
			},
		},
		{
			name: "first FindBySKUs error",
			job:  CheckoutJob{CheckoutID: checkoutID, Items: map[string]int{"SKU1": 1}},
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor) {
				p.EXPECT().FindBySKUs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("db error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockProduct, mockCheckout, mockTx)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			err := uc.ProcessJob(context.Background(), tt.job)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckoutUsecaseConfirm(t *testing.T) {
	checkoutID := uuid.New()

	items := []domain.CheckoutItem{
		{SKU: "SKU1", Name: "Product 1", Qty: 2, Price: 50.0, Total: 100.0},
	}
	result := domain.CheckoutResult{
		Items:             items,
		PromotionsApplied: []domain.AppliedPromotion{},
		Subtotal:          100.0,
		TotalDiscount:     0.0,
		Total:             100.0,
	}
	resultJSON, _ := json.Marshal(result)

	tests := []struct {
		name        string
		setupMocks  func(*mocks.MockProductRepository, *mocks.MockCheckoutRepository, *mocks.MockOrderRepository, *mocks.MockTransactor)
		expectError bool
		expectCode  string
	}{
		{
			name: "session not found",
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				c.EXPECT().FindByID(gomock.Any(), checkoutID).Return(nil, fmt.Errorf("not found"))
			},
			expectError: true,
			expectCode:  "ERR-UC-016-404",
		},
		{
			name: "session not completed",
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				c.EXPECT().FindByID(gomock.Any(), checkoutID).Return(&domain.CheckoutSession{
					ID:        checkoutID,
					Status:    domain.CheckoutPending,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}, nil)
			},
			expectError: true,
			expectCode:  "ERR-UC-017",
		},
		{
			name: "session expired",
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				c.EXPECT().FindByID(gomock.Any(), checkoutID).Return(&domain.CheckoutSession{
					ID:        checkoutID,
					Status:    domain.CheckoutCompleted,
					ExpiresAt: time.Now().Add(-10 * time.Minute),
					Result:    resultJSON,
				}, nil)
			},
			expectError: true,
			expectCode:  "ERR-UC-018",
		},
		{
			name: "success creates order",
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, o *mocks.MockOrderRepository, tx *mocks.MockTransactor) {
				c.EXPECT().FindByID(gomock.Any(), checkoutID).Return(&domain.CheckoutSession{
					ID:        checkoutID,
					Status:    domain.CheckoutCompleted,
					ExpiresAt: time.Now().Add(10 * time.Minute),
					Result:    resultJSON,
				}, nil)
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				p.EXPECT().DecrementInventory(gomock.Any(), "SKU1", 2).Return(nil)
				p.EXPECT().DecrementReserved(gomock.Any(), "SKU1", 2).Return(nil)
				o.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, ord *domain.Order) error {
						ord.ID = uuid.New()
						return nil
					})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockProduct, mockCheckout, mockOrder, mockTx)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			order, err := uc.Confirm(context.Background(), checkoutID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, order)
				if tt.expectCode != "" {
					appErr, ok := util.IsAppError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.expectCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
			}
		})
	}
}

func TestCheckoutUsecaseReleaseExpired(t *testing.T) {
	expiredSession := domain.CheckoutSession{
		ID:        uuid.New(),
		Status:    domain.CheckoutPending,
		ExpiresAt: time.Now().Add(-5 * time.Minute),
		Items:     []byte(`{"SKU1":2}`),
	}

	tests := []struct {
		name        string
		setupMocks  func(*mocks.MockProductRepository, *mocks.MockCheckoutRepository, *mocks.MockTransactor)
		expectError bool
	}{
		{
			name: "no expired sessions",
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor) {
				c.EXPECT().FindExpiredPending(gomock.Any()).Return([]domain.CheckoutSession{}, nil)
			},
		},
		{
			name: "find expired error",
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor) {
				c.EXPECT().FindExpiredPending(gomock.Any()).Return(nil, fmt.Errorf("db error"))
			},
			expectError: true,
		},
		{
			name: "releases expired sessions",
			setupMocks: func(p *mocks.MockProductRepository, c *mocks.MockCheckoutRepository, tx *mocks.MockTransactor) {
				c.EXPECT().FindExpiredPending(gomock.Any()).Return([]domain.CheckoutSession{expiredSession}, nil)
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(txPassthrough)
				p.EXPECT().DecrementReserved(gomock.Any(), "SKU1", 2).Return(nil)
				c.EXPECT().UpdateStatus(gomock.Any(), expiredSession.ID, domain.CheckoutExpired, nil, nil).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockProduct := mocks.NewMockProductRepository(ctrl)
			mockCheckout := mocks.NewMockCheckoutRepository(ctrl)
			mockOrder := mocks.NewMockOrderRepository(ctrl)
			mockUser := mocks.NewMockUserRepository(ctrl)
			mockTx := mocks.NewMockTransactor(ctrl)

			tt.setupMocks(mockProduct, mockCheckout, mockTx)

			uc := NewCheckoutUsecase(mockProduct, mockCheckout, mockOrder, mockUser,
				promotion.NewStaticAdapter(), nil, mockTx, util.NewLogger())

			err := uc.ReleaseExpired(context.Background())

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
