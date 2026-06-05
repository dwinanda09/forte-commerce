package handler

import (
	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productRepo domain.ProductRepository
}

func NewProductHandler(productRepo domain.ProductRepository) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

type ProductResponse struct {
	ID           string  `json:"id"`
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	AvailableQty int     `json:"available_qty"`
	InventoryQty int     `json:"inventory_qty"`
	ReservedQty  int     `json:"reserved_qty"`
}

type CreateProductRequest struct {
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	InventoryQty int     `json:"inventory_qty"`
}

type UpdateProductRequest struct {
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	InventoryQty int     `json:"inventory_qty"`
}

func (h *ProductHandler) List(c echo.Context) error {
	products, err := h.productRepo.FindAll(c.Request().Context())
	if err != nil {
		return util.Fail(c, 500, "ERR-HDL-003", "Failed to fetch products")
	}

	responses := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, ProductResponse{
			ID:           p.ID.String(),
			SKU:          p.SKU,
			Name:         p.Name,
			Price:        p.Price,
			AvailableQty: p.Available(),
			InventoryQty: p.InventoryQty,
			ReservedQty:  p.ReservedQty,
		})
	}

	return util.OK(c, responses)
}

func (h *ProductHandler) Create(c echo.Context) error {
	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, 400, "ERR-HDL-023", "Invalid request body")
	}
	if req.SKU == "" || req.Name == "" || req.Price <= 0 {
		return util.Fail(c, 400, "ERR-HDL-024", "SKU, name, and price are required")
	}
	p := &domain.Product{
		SKU:          req.SKU,
		Name:         req.Name,
		Price:        req.Price,
		InventoryQty: req.InventoryQty,
	}
	if err := h.productRepo.Create(c.Request().Context(), p); err != nil {
		return util.Fail(c, 500, "ERR-HDL-025", "Failed to create product")
	}
	return util.OK(c, ProductResponse{
		ID:           p.ID.String(),
		SKU:          p.SKU,
		Name:         p.Name,
		Price:        p.Price,
		AvailableQty: p.Available(),
		InventoryQty: p.InventoryQty,
		ReservedQty:  p.ReservedQty,
	})
}

func (h *ProductHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-026", "Invalid product ID")
	}
	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, 400, "ERR-HDL-027", "Invalid request body")
	}
	p, err := h.productRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-028", "Product not found")
		}
		return util.Fail(c, 500, "ERR-HDL-029", "Failed to find product")
	}
	if req.SKU != "" {
		p.SKU = req.SKU
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.InventoryQty >= 0 {
		p.InventoryQty = req.InventoryQty
	}
	if err := h.productRepo.Update(c.Request().Context(), p); err != nil {
		return util.Fail(c, 500, "ERR-HDL-030", "Failed to update product")
	}
	return util.OK(c, ProductResponse{
		ID:           p.ID.String(),
		SKU:          p.SKU,
		Name:         p.Name,
		Price:        p.Price,
		AvailableQty: p.Available(),
		InventoryQty: p.InventoryQty,
		ReservedQty:  p.ReservedQty,
	})
}

func (h *ProductHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-031", "Invalid product ID")
	}
	if err := h.productRepo.Delete(c.Request().Context(), id); err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-032", "Product not found")
		}
		return util.Fail(c, 500, "ERR-HDL-033", "Failed to delete product")
	}
	return util.OK(c, map[string]string{"status": "deleted"})
}
