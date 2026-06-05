package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/dwinanda09/forte-commerce/internal/domain"
)

func TestProductAvailable(t *testing.T) {
	tests := []struct {
		name         string
		inventoryQty int
		reservedQty  int
		expected     int
	}{
		{
			name:         "normal case: inventory exceeds reserved",
			inventoryQty: 10,
			reservedQty:  3,
			expected:     7,
		},
		{
			name:         "edge case: inventory equals reserved",
			inventoryQty: 5,
			reservedQty:  5,
			expected:     0,
		},
		{
			name:         "edge case: both zero",
			inventoryQty: 0,
			reservedQty:  0,
			expected:     0,
		},
		{
			name:         "negative case: reserved exceeds inventory",
			inventoryQty: 3,
			reservedQty:  5,
			expected:     -2,
		},
		{
			name:         "large quantities",
			inventoryQty: 10000,
			reservedQty:  2500,
			expected:     7500,
		},
		{
			name:         "small quantities",
			inventoryQty: 1,
			reservedQty:  0,
			expected:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &domain.Product{
				ID:           uuid.New(),
				SKU:          "TEST-SKU-001",
				Name:         "Test Product",
				Price:        99.99,
				InventoryQty: tt.inventoryQty,
				ReservedQty:  tt.reservedQty,
			}

			available := product.Available()

			assert.Equal(t, tt.expected, available)
		})
	}
}

func TestProductAvailableWithZeroInventory(t *testing.T) {
	product := &domain.Product{
		ID:           uuid.New(),
		SKU:          "OUT-OF-STOCK",
		Name:         "Out of Stock Product",
		Price:        49.99,
		InventoryQty: 0,
		ReservedQty:  0,
	}

	assert.Equal(t, 0, product.Available())
}

func TestProductAvailableAfterPartialReservation(t *testing.T) {
	product := &domain.Product{
		ID:           uuid.New(),
		SKU:          "PARTIAL-RESERVE",
		Name:         "Partially Reserved Product",
		Price:        199.99,
		InventoryQty: 100,
		ReservedQty:  25,
	}

	assert.Equal(t, 75, product.Available())
}
