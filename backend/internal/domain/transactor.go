package domain

import (
	"context"
	"database/sql"
)

type txKey struct{}

// Transactor defines the contract for managing database transactions.
type Transactor interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// WithTx adds a transaction to the context.
func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromContext retrieves a transaction from the context if present.
func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}
