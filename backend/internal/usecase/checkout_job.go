package usecase

import "github.com/google/uuid"

type CheckoutJob struct {
	CheckoutID uuid.UUID      `json:"checkout_id"`
	Items      map[string]int `json:"items"`
}
