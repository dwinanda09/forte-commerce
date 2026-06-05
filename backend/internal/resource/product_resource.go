package resource

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductResource struct {
	db     *sqlx.DB
	logger *util.Logger
}

func NewProductResource(db *sqlx.DB, logger *util.Logger) *ProductResource {
	return &ProductResource{
		db:     db,
		logger: logger,
	}
}

// dbConn is a helper interface that both *sqlx.DB and *sql.Tx implement
type dbConn interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// conn returns either a transaction from context or the database connection
func (r *ProductResource) conn(ctx context.Context) dbConn {
	if tx, ok := domain.TxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *ProductResource) FindAll(ctx context.Context) ([]domain.Product, error) {
	start := r.logger.Start(ctx, "ProductResource.FindAll")
	defer func() { r.logger.Finish(ctx, "ProductResource.FindAll", start, nil) }()

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, sku, name, price, inventory_qty, reserved_qty, created_at
		FROM products
		ORDER BY created_at DESC
	`)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.FindAll", start, err)
		return nil, util.Wrap("ERR-RS-001", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.FindAll", start, err)
		return nil, util.Wrap("ERR-RS-001", "Failed to fetch products", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.InventoryQty, &p.ReservedQty, &p.CreatedAt); err != nil {
			r.logger.Finish(ctx, "ProductResource.FindAll", start, err)
			return nil, util.Wrap("ERR-RS-001", "Failed to scan products", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		r.logger.Finish(ctx, "ProductResource.FindAll", start, err)
		return nil, util.Wrap("ERR-RS-001", "Failed to iterate products", err)
	}

	return products, nil
}

func (r *ProductResource) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	start := r.logger.Start(ctx, "ProductResource.FindBySKU")
	defer func() { r.logger.Finish(ctx, "ProductResource.FindBySKU", start, nil) }()

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, sku, name, price, inventory_qty, reserved_qty, created_at
		FROM products
		WHERE sku = $1
	`)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.FindBySKU", start, err)
		return nil, util.Wrap("ERR-RS-002", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	var product domain.Product
	err = stmt.QueryRowContext(ctx, sku).Scan(
		&product.ID, &product.SKU, &product.Name, &product.Price,
		&product.InventoryQty, &product.ReservedQty, &product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Finish(ctx, "ProductResource.FindBySKU", start, err)
			return nil, util.Wrap("ERR-RS-002-404", "Product not found", err)
		}
		r.logger.Finish(ctx, "ProductResource.FindBySKU", start, err)
		return nil, util.Wrap("ERR-RS-002", "Failed to find product", err)
	}

	return &product, nil
}

func (r *ProductResource) FindBySKUs(ctx context.Context, skus []string) ([]domain.Product, error) {
	start := r.logger.Start(ctx, "ProductResource.FindBySKUs")
	defer func() { r.logger.Finish(ctx, "ProductResource.FindBySKUs", start, nil) }()

	if len(skus) == 0 {
		return []domain.Product{}, nil
	}

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, sku, name, price, inventory_qty, reserved_qty, created_at
		FROM products
		WHERE sku = ANY($1)
		ORDER BY sku
	`)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.FindBySKUs", start, err)
		return nil, util.Wrap("ERR-RS-004", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, pq.Array(skus))
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.FindBySKUs", start, err)
		return nil, util.Wrap("ERR-RS-004", "Failed to fetch products", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.InventoryQty, &p.ReservedQty, &p.CreatedAt); err != nil {
			r.logger.Finish(ctx, "ProductResource.FindBySKUs", start, err)
			return nil, util.Wrap("ERR-RS-004", "Failed to scan products", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductResource) IncrementReserved(ctx context.Context, sku string, qty int) error {
	start := r.logger.Start(ctx, "ProductResource.IncrementReserved")
	defer func() { r.logger.Finish(ctx, "ProductResource.IncrementReserved", start, nil) }()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		UPDATE products
		SET reserved_qty = reserved_qty + $1
		WHERE sku = $2
	`)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.IncrementReserved", start, err)
		return util.Wrap("ERR-RS-005", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, qty, sku)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.IncrementReserved", start, err)
		return util.Wrap("ERR-RS-005", "Failed to increment reserved quantity", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.IncrementReserved", start, err)
		return util.Wrap("ERR-RS-006", "Failed to check affected rows", err)
	}
	if n == 0 {
		err := fmt.Errorf("product not found")
		r.logger.Finish(ctx, "ProductResource.IncrementReserved", start, err)
		return util.Wrap("ERR-RS-007", "Product not found", err)
	}

	return nil
}

func (r *ProductResource) DecrementReserved(ctx context.Context, sku string, qty int) error {
	start := r.logger.Start(ctx, "ProductResource.DecrementReserved")
	defer func() { r.logger.Finish(ctx, "ProductResource.DecrementReserved", start, nil) }()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		UPDATE products
		SET reserved_qty = reserved_qty - $1
		WHERE sku = $2 AND reserved_qty >= $1
	`)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.DecrementReserved", start, err)
		return util.Wrap("ERR-RS-008", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, qty, sku)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.DecrementReserved", start, err)
		return util.Wrap("ERR-RS-008", "Failed to decrement reserved quantity", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.DecrementReserved", start, err)
		return util.Wrap("ERR-RS-009", "Failed to check affected rows", err)
	}
	if n == 0 {
		err := fmt.Errorf("product not found or insufficient reserved quantity")
		r.logger.Finish(ctx, "ProductResource.DecrementReserved", start, err)
		return util.Wrap("ERR-RS-010", "Failed to decrement reserved quantity", err)
	}

	return nil
}

func (r *ProductResource) DecrementInventory(ctx context.Context, sku string, qty int) error {
	start := r.logger.Start(ctx, "ProductResource.DecrementInventory")
	defer func() { r.logger.Finish(ctx, "ProductResource.DecrementInventory", start, nil) }()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		UPDATE products
		SET inventory_qty = inventory_qty - $1
		WHERE sku = $2 AND inventory_qty >= $1
	`)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.DecrementInventory", start, err)
		return util.Wrap("ERR-RS-011", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, qty, sku)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.DecrementInventory", start, err)
		return util.Wrap("ERR-RS-011", "Failed to decrement inventory", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.DecrementInventory", start, err)
		return util.Wrap("ERR-RS-012", "Failed to check affected rows", err)
	}
	if n == 0 {
		err := fmt.Errorf("product not found or insufficient inventory")
		r.logger.Finish(ctx, "ProductResource.DecrementInventory", start, err)
		return util.Wrap("ERR-RS-013", "Failed to decrement inventory", err)
	}

	return nil
}

func (r *ProductResource) RestoreInventory(ctx context.Context, sku string, qty int) error {
	start := r.logger.Start(ctx, "ProductResource.RestoreInventory")
	defer func() { r.logger.Finish(ctx, "ProductResource.RestoreInventory", start, nil) }()

	stmt, err := r.conn(ctx).PrepareContext(ctx, `
		UPDATE products
		SET inventory_qty = inventory_qty + $1
		WHERE sku = $2
	`)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.RestoreInventory", start, err)
		return util.Wrap("ERR-RS-014", "Failed to prepare statement", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, qty, sku)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.RestoreInventory", start, err)
		return util.Wrap("ERR-RS-014", "Failed to restore inventory", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.RestoreInventory", start, err)
		return util.Wrap("ERR-RS-015", "Failed to check affected rows", err)
	}
	if n == 0 {
		err := fmt.Errorf("product not found")
		r.logger.Finish(ctx, "ProductResource.RestoreInventory", start, err)
		return util.Wrap("ERR-RS-016", "Product not found", err)
	}

	return nil
}

func (r *ProductResource) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	start := r.logger.Start(ctx, "ProductResource.FindByID")
	defer func() { r.logger.Finish(ctx, "ProductResource.FindByID", start, nil) }()

	var p domain.Product
	err := r.db.QueryRowContext(ctx, `
		SELECT id, sku, name, price, inventory_qty, reserved_qty, created_at
		FROM products WHERE id = $1
	`, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.InventoryQty, &p.ReservedQty, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Finish(ctx, "ProductResource.FindByID", start, err)
			return nil, util.Wrap("ERR-RS-031-404", "Product not found", err)
		}
		r.logger.Finish(ctx, "ProductResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-031", "Failed to find product", err)
	}
	return &p, nil
}

func (r *ProductResource) Create(ctx context.Context, product *domain.Product) error {
	start := r.logger.Start(ctx, "ProductResource.Create")
	defer func() { r.logger.Finish(ctx, "ProductResource.Create", start, nil) }()

	product.ID = uuid.New()
	product.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO products (id, sku, name, price, inventory_qty, reserved_qty, created_at)
		VALUES ($1, $2, $3, $4, $5, 0, $6)
	`, product.ID, product.SKU, product.Name, product.Price, product.InventoryQty, product.CreatedAt)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.Create", start, err)
		return util.Wrap("ERR-RS-032", "Failed to create product", err)
	}
	return nil
}

func (r *ProductResource) Update(ctx context.Context, product *domain.Product) error {
	start := r.logger.Start(ctx, "ProductResource.Update")
	defer func() { r.logger.Finish(ctx, "ProductResource.Update", start, nil) }()

	_, err := r.db.ExecContext(ctx, `
		UPDATE products SET name = $1, price = $2, inventory_qty = $3, sku = $4 WHERE id = $5
	`, product.Name, product.Price, product.InventoryQty, product.SKU, product.ID)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.Update", start, err)
		return util.Wrap("ERR-RS-033", "Failed to update product", err)
	}
	return nil
}

func (r *ProductResource) Delete(ctx context.Context, id uuid.UUID) error {
	start := r.logger.Start(ctx, "ProductResource.Delete")
	defer func() { r.logger.Finish(ctx, "ProductResource.Delete", start, nil) }()

	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		r.logger.Finish(ctx, "ProductResource.Delete", start, err)
		return util.Wrap("ERR-RS-034", "Failed to delete product", err)
	}
	return nil
}
