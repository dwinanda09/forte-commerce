package handler

import (
	"encoding/json"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/internal/usecase"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CheckoutHandler struct {
	uc *usecase.CheckoutUsecase
}

func NewCheckoutHandler(uc *usecase.CheckoutUsecase) *CheckoutHandler {
	return &CheckoutHandler{uc: uc}
}

type SubmitCheckoutRequest struct {
	SKUs []string `json:"items"`
}

type SubmitCheckoutResponse struct {
	CheckoutID string `json:"checkout_id"`
}

type CheckoutSessionResponse struct {
	ID           string      `json:"id"`
	Status       string      `json:"status"`
	ExpiresAt    string      `json:"expires_at"`
	ErrorMessage *string     `json:"error_message,omitempty"`
	Result       interface{} `json:"result,omitempty"`
}

type OrderItemResponse struct {
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
	Total float64 `json:"total"`
}

type OrderResponse struct {
	ID                string                  `json:"id"`
	CheckoutSessionID string                  `json:"checkout_session_id"`
	Status            string                  `json:"status"`
	Items             []OrderItemResponse     `json:"items"`
	PromotionsApplied []domain.AppliedPromotion `json:"promotions_applied"`
	Subtotal          float64                 `json:"subtotal"`
	TotalDiscount     float64                 `json:"total_discount"`
	Total             float64                 `json:"total"`
	CreatedAt         string                  `json:"created_at"`
	UpdatedAt         string                  `json:"updated_at"`
}

func (h *CheckoutHandler) Submit(c echo.Context) error {
	var req SubmitCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, 400, "ERR-HDL-004", "Invalid request body")
	}

	if len(req.SKUs) == 0 {
		return util.Fail(c, 400, "ERR-HDL-005", "Items list cannot be empty")
	}

	checkoutID, err := h.uc.Submit(c.Request().Context(), req.SKUs)
	if err != nil {
		return util.Fail(c, 500, "ERR-HDL-006", "Failed to submit checkout")
	}

	return c.JSON(202, map[string]interface{}{
		"success": true,
		"data": SubmitCheckoutResponse{
			CheckoutID: checkoutID.String(),
		},
		"meta": map[string]string{
			"request_id": c.Get("request_id").(string),
		},
	})
}

func (h *CheckoutHandler) GetSession(c echo.Context) error {
	idStr := c.Param("id")
	checkoutID, err := uuid.Parse(idStr)
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-007", "Invalid checkout ID")
	}

	session, err := h.uc.GetSession(c.Request().Context(), checkoutID)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-008", "Checkout session not found")
		}
		return util.Fail(c, 500, "ERR-HDL-009", "Failed to get checkout session")
	}

	// Parse result if it exists
	var result interface{}
	if session.Result != nil {
		var r domain.CheckoutResult
		json.Unmarshal(session.Result, &r)
		result = r
	}

	resp := CheckoutSessionResponse{
		ID:        session.ID.String(),
		Status:    string(session.Status),
		ExpiresAt: session.ExpiresAt.String(),
		Result:    result,
	}

	if session.ErrorMessage != nil {
		resp.ErrorMessage = session.ErrorMessage
	}

	return util.OK(c, resp)
}

func (h *CheckoutHandler) Confirm(c echo.Context) error {
	idStr := c.Param("id")
	checkoutID, err := uuid.Parse(idStr)
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-010", "Invalid checkout ID")
	}

	order, err := h.uc.Confirm(c.Request().Context(), checkoutID)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-011", "Checkout session not found")
		}
		return util.Fail(c, 500, "ERR-HDL-012", "Failed to confirm checkout")
	}

	// Parse items
	var items []domain.CheckoutItem
	json.Unmarshal(order.Items, &items)

	var itemResponses []OrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, OrderItemResponse{
			SKU:   item.SKU,
			Name:  item.Name,
			Qty:   item.Qty,
			Price: item.Price,
			Total: item.Total,
		})
	}

	// Parse promotions
	var promos []domain.AppliedPromotion
	json.Unmarshal(order.PromotionsApplied, &promos)

	resp := OrderResponse{
		ID:                order.ID.String(),
		CheckoutSessionID: order.CheckoutSessionID.String(),
		Status:            string(order.Status),
		Items:             itemResponses,
		PromotionsApplied: promos,
		Subtotal:          order.Subtotal,
		TotalDiscount:     order.TotalDiscount,
		Total:             order.Total,
		CreatedAt:         order.CreatedAt.String(),
		UpdatedAt:         order.UpdatedAt.String(),
	}

	return util.OK(c, resp)
}

func (h *CheckoutHandler) PayOrder(c echo.Context) error {
	idStr := c.Param("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-013", "Invalid order ID")
	}

	order, err := h.uc.PayOrder(c.Request().Context(), orderID)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-014", "Order not found")
		}
		return util.Fail(c, 500, "ERR-HDL-015", "Failed to pay order")
	}

	return util.OK(c, map[string]string{
		"status": string(order.Status),
	})
}

func (h *CheckoutHandler) CancelOrder(c echo.Context) error {
	idStr := c.Param("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-016", "Invalid order ID")
	}

	order, err := h.uc.CancelOrder(c.Request().Context(), orderID)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-017", "Order not found")
		}
		return util.Fail(c, 500, "ERR-HDL-018", "Failed to cancel order")
	}

	return util.OK(c, map[string]string{
		"status": string(order.Status),
	})
}

func (h *CheckoutHandler) GetOrder(c echo.Context) error {
	idStr := c.Param("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-019", "Invalid order ID")
	}

	order, err := h.uc.GetOrder(c.Request().Context(), orderID)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-020", "Order not found")
		}
		return util.Fail(c, 500, "ERR-HDL-021", "Failed to get order")
	}

	// Parse items
	var items []domain.CheckoutItem
	json.Unmarshal(order.Items, &items)

	var itemResponses []OrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, OrderItemResponse{
			SKU:   item.SKU,
			Name:  item.Name,
			Qty:   item.Qty,
			Price: item.Price,
			Total: item.Total,
		})
	}

	// Parse promotions
	var promos []domain.AppliedPromotion
	json.Unmarshal(order.PromotionsApplied, &promos)

	resp := OrderResponse{
		ID:                order.ID.String(),
		CheckoutSessionID: order.CheckoutSessionID.String(),
		Status:            string(order.Status),
		Items:             itemResponses,
		PromotionsApplied: promos,
		Subtotal:          order.Subtotal,
		TotalDiscount:     order.TotalDiscount,
		Total:             order.Total,
		CreatedAt:         order.CreatedAt.String(),
		UpdatedAt:         order.UpdatedAt.String(),
	}

	return util.OK(c, resp)
}

func (h *CheckoutHandler) ListOrders(c echo.Context) error {
	orders, err := h.uc.ListOrders(c.Request().Context())
	if err != nil {
		return util.Fail(c, 500, "ERR-HDL-022", "Failed to list orders")
	}

	responses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		// Parse items
		var items []domain.CheckoutItem
		json.Unmarshal(order.Items, &items)

		var itemResponses []OrderItemResponse
		for _, item := range items {
			itemResponses = append(itemResponses, OrderItemResponse{
				SKU:   item.SKU,
				Name:  item.Name,
				Qty:   item.Qty,
				Price: item.Price,
				Total: item.Total,
			})
		}

		// Parse promotions
		var promos []domain.AppliedPromotion
		json.Unmarshal(order.PromotionsApplied, &promos)

		responses = append(responses, OrderResponse{
			ID:                order.ID.String(),
			CheckoutSessionID: order.CheckoutSessionID.String(),
			Status:            string(order.Status),
			Items:             itemResponses,
			PromotionsApplied: promos,
			Subtotal:          order.Subtotal,
			TotalDiscount:     order.TotalDiscount,
			Total:             order.Total,
			CreatedAt:         order.CreatedAt.String(),
			UpdatedAt:         order.UpdatedAt.String(),
		})
	}

	return util.OK(c, responses)
}
