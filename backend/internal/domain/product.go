package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID `db:"id" json:"id"`
	SKU          string    `db:"sku" json:"sku"`
	Name         string    `db:"name" json:"name"`
	Price        float64   `db:"price" json:"price"`
	InventoryQty int       `db:"inventory_qty" json:"inventory_qty"`
	ReservedQty  int       `db:"reserved_qty" json:"reserved_qty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func (p *Product) Available() int {
	return p.InventoryQty - p.ReservedQty
}

type ProductRepository interface {
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	FindBySKU(ctx context.Context, sku string) (*Product, error)
	FindBySKUs(ctx context.Context, skus []string) ([]Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementReserved(ctx context.Context, sku string, qty int) error
	DecrementReserved(ctx context.Context, sku string, qty int) error
	DecrementInventory(ctx context.Context, sku string, qty int) error
	RestoreInventory(ctx context.Context, sku string, qty int) error
}
