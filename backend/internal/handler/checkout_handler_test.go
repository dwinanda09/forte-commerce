package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/internal/mocks"
	"github.com/dwinanda09/forte-commerce/internal/promotion"
	"github.com/dwinanda09/forte-commerce/internal/usecase"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type checkoutTestMocks struct {
	product  *mocks.MockProductRepository
	checkout *mocks.MockCheckoutRepository
	order    *mocks.MockOrderRepository
	user     *mocks.MockUserRepository
	queue    *mocks.MockQueuePublisher
	tx       *mocks.MockTransactor
}

func newCheckoutTestMocks(ctrl *gomock.Controller) checkoutTestMocks {
	return checkoutTestMocks{
		product:  mocks.NewMockProductRepository(ctrl),
		checkout: mocks.NewMockCheckoutRepository(ctrl),
		order:    mocks.NewMockOrderRepository(ctrl),
		user:     mocks.NewMockUserRepository(ctrl),
		queue:    mocks.NewMockQueuePublisher(ctrl),
		tx:       mocks.NewMockTransactor(ctrl),
	}
}

func (m checkoutTestMocks) newHandler(logger *util.Logger) *CheckoutHandler {
	uc := usecase.NewCheckoutUsecase(
		m.product, m.checkout, m.order, m.user,
		promotion.NewStaticAdapter(), m.queue, m.tx, logger,
	)
	return NewCheckoutHandler(uc)
}

func TestCheckoutHandlerSubmitInvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/checkout/submit", bytes.NewReader([]byte("invalid")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.Set("request_id", uuid.New().String())

	err := handler.Submit(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCheckoutHandlerGetSession(t *testing.T) {
	sessionID := uuid.New()

	tests := []struct {
		name         string
		paramID      string
		setupMocks   func(m checkoutTestMocks)
		expectStatus int
	}{
		{
			name:    "get existing session",
			paramID: sessionID.String(),
			setupMocks: func(m checkoutTestMocks) {
				m.checkout.EXPECT().FindByID(gomock.Any(), sessionID).Return(&domain.CheckoutSession{
					ID:        sessionID,
					Status:    domain.CheckoutPending,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:    "session not found",
			paramID: uuid.New().String(),
			setupMocks: func(m checkoutTestMocks) {
				m.checkout.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "invalid session ID format",
			paramID:      "invalid-uuid",
			setupMocks:   func(m checkoutTestMocks) {},
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := newCheckoutTestMocks(ctrl)
			tt.setupMocks(m)
			handler := m.newHandler(util.NewLogger())

			req := httptest.NewRequest(http.MethodGet, "/checkout/:id", nil)
			e := echo.New()
			c := e.NewContext(req, httptest.NewRecorder())
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := handler.GetSession(c)
			require.NoError(t, err)

			recorder := c.Response().Writer.(*httptest.ResponseRecorder)
			assert.Equal(t, tt.expectStatus, recorder.Code)
		})
	}
}

func TestCheckoutHandlerConfirmSessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	m.checkout.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found"))
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/checkout/:id/confirm", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(uuid.New().String())

	err := handler.Confirm(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestCheckoutHandlerConfirmInvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/checkout/:id/confirm", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	err := handler.Confirm(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCheckoutHandlerPayOrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	m.order.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found"))
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/order/:id/pay", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(uuid.New().String())

	err := handler.PayOrder(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestCheckoutHandlerPayOrderInvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/order/:id/pay", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	err := handler.PayOrder(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCheckoutHandlerCancelOrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	m.order.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found"))
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/order/:id/cancel", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(uuid.New().String())

	err := handler.CancelOrder(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestCheckoutHandlerCancelOrderInvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/order/:id/cancel", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	err := handler.CancelOrder(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCheckoutHandlerGetOrder(t *testing.T) {
	orderID := uuid.New()

	items := []domain.CheckoutItem{
		{SKU: "SKU1", Name: "Product 1", Qty: 1, Price: 100.0, Total: 100.0},
	}
	itemsJSON, _ := json.Marshal(items)

	tests := []struct {
		name         string
		paramID      string
		setupMocks   func(m checkoutTestMocks)
		expectStatus int
	}{
		{
			name:    "get existing order",
			paramID: orderID.String(),
			setupMocks: func(m checkoutTestMocks) {
				m.order.EXPECT().FindByID(gomock.Any(), orderID).Return(&domain.Order{
					ID:                orderID,
					Status:            domain.OrderPaid,
					Items:             itemsJSON,
					PromotionsApplied: []byte("[]"),
					Subtotal:          100.0,
					Total:             100.0,
				}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:    "order not found",
			paramID: uuid.New().String(),
			setupMocks: func(m checkoutTestMocks) {
				m.order.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "invalid order ID format",
			paramID:      "invalid-uuid",
			setupMocks:   func(m checkoutTestMocks) {},
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := newCheckoutTestMocks(ctrl)
			tt.setupMocks(m)
			handler := m.newHandler(util.NewLogger())

			req := httptest.NewRequest(http.MethodGet, "/order/:id", nil)
			e := echo.New()
			c := e.NewContext(req, httptest.NewRecorder())
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := handler.GetOrder(c)
			require.NoError(t, err)

			recorder := c.Response().Writer.(*httptest.ResponseRecorder)
			assert.Equal(t, tt.expectStatus, recorder.Code)
		})
	}
}

func TestCheckoutHandlerListOrders(t *testing.T) {
	items := []domain.CheckoutItem{
		{SKU: "SKU1", Name: "Product 1", Qty: 1, Price: 100.0, Total: 100.0},
	}
	itemsJSON, _ := json.Marshal(items)

	order1 := domain.Order{
		ID:                uuid.New(),
		Status:            domain.OrderPaid,
		Items:             itemsJSON,
		PromotionsApplied: []byte("[]"),
		Subtotal:          100.0,
		Total:             100.0,
	}
	order2 := domain.Order{
		ID:                uuid.New(),
		Status:            domain.OrderPending,
		Items:             itemsJSON,
		PromotionsApplied: []byte("[]"),
		Subtotal:          100.0,
		Total:             100.0,
	}

	tests := []struct {
		name         string
		setupMocks   func(m checkoutTestMocks)
		expectStatus int
		expectCount  int
	}{
		{
			name: "list multiple orders",
			setupMocks: func(m checkoutTestMocks) {
				m.order.EXPECT().FindAll(gomock.Any()).Return([]domain.Order{order1, order2}, nil)
			},
			expectStatus: http.StatusOK,
			expectCount:  2,
		},
		{
			name: "empty order list",
			setupMocks: func(m checkoutTestMocks) {
				m.order.EXPECT().FindAll(gomock.Any()).Return([]domain.Order{}, nil)
			},
			expectStatus: http.StatusOK,
			expectCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := newCheckoutTestMocks(ctrl)
			tt.setupMocks(m)
			handler := m.newHandler(util.NewLogger())

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			e := echo.New()
			c := e.NewContext(req, httptest.NewRecorder())

			err := handler.ListOrders(c)
			require.NoError(t, err)

			recorder := c.Response().Writer.(*httptest.ResponseRecorder)
			assert.Equal(t, tt.expectStatus, recorder.Code)

			var resp map[string]any
			json.Unmarshal(recorder.Body.Bytes(), &resp)
			data, ok := resp["data"].([]any)
			if !ok {
				data = []any{}
			}
			assert.Equal(t, tt.expectCount, len(data))
		})
	}
}

func TestCheckoutHandlerGetSessionWithResult(t *testing.T) {
	sessionID := uuid.New()

	items := []domain.CheckoutItem{
		{SKU: "SKU1", Name: "Product 1", Qty: 1, Price: 100.0, Total: 100.0},
	}
	result := domain.CheckoutResult{
		Items:             items,
		PromotionsApplied: []domain.AppliedPromotion{},
		Subtotal:          100.0,
		TotalDiscount:     0.0,
		Total:             100.0,
	}
	resultJSON, _ := json.Marshal(result)

	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	m.checkout.EXPECT().FindByID(gomock.Any(), sessionID).Return(&domain.CheckoutSession{
		ID:        sessionID,
		Status:    domain.CheckoutCompleted,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Result:    resultJSON,
	}, nil)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodGet, "/checkout/:id", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(sessionID.String())

	err := handler.GetSession(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var resp map[string]any
	json.Unmarshal(recorder.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	assert.NotNil(t, data["result"])
}

func TestCheckoutHandlerGetSessionWithErrorMessage(t *testing.T) {
	sessionID := uuid.New()
	errMsg := "Insufficient inventory for SKU1"

	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	m.checkout.EXPECT().FindByID(gomock.Any(), sessionID).Return(&domain.CheckoutSession{
		ID:           sessionID,
		Status:       domain.CheckoutFailed,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
		ErrorMessage: &errMsg,
	}, nil)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodGet, "/checkout/:id", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(sessionID.String())

	err := handler.GetSession(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var resp map[string]any
	json.Unmarshal(recorder.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, errMsg, data["error_message"])
}

func TestCheckoutHandlerConfirmSessionNotCompleted(t *testing.T) {
	checkoutID := uuid.New()

	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	m.checkout.EXPECT().FindByID(gomock.Any(), checkoutID).Return(&domain.CheckoutSession{
		ID:        checkoutID,
		Status:    domain.CheckoutPending,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}, nil)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/checkout/:id/confirm", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(checkoutID.String())

	err := handler.Confirm(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestCheckoutHandlerConfirmSessionExpired(t *testing.T) {
	checkoutID := uuid.New()

	ctrl := gomock.NewController(t)
	m := newCheckoutTestMocks(ctrl)
	m.checkout.EXPECT().FindByID(gomock.Any(), checkoutID).Return(&domain.CheckoutSession{
		ID:        checkoutID,
		Status:    domain.CheckoutCompleted,
		ExpiresAt: time.Now().Add(-10 * time.Minute),
	}, nil)
	handler := m.newHandler(util.NewLogger())

	req := httptest.NewRequest(http.MethodPost, "/checkout/:id/confirm", nil)
	e := echo.New()
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(checkoutID.String())

	err := handler.Confirm(c)
	require.NoError(t, err)

	recorder := c.Response().Writer.(*httptest.ResponseRecorder)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
