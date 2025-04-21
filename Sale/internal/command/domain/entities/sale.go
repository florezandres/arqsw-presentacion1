package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Sale struct {
	ID         uuid.UUID
	ProductID  string
	Quantity   int
	TotalPrice float64
	SaleDate   time.Time
}

func NewSale(productID string, quantity int, totalPrice float64, saleDate time.Time) (*Sale, error) {
	if productID == "" {
		return nil, errors.New("productID cannot be empty")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if totalPrice <= 0 {
		return nil, errors.New("totalPrice must be positive")
	}
	if saleDate.IsZero() {
		return nil, errors.New("saleDate cannot be zero")
	}

	return &Sale{
		ID:         uuid.New(),
		ProductID:  productID,
		Quantity:   quantity,
		TotalPrice: totalPrice,
		SaleDate:   saleDate,
	}, nil
}

// Métodos de negocio
func (s *Sale) UpdateTotalPrice(price float64) {
	s.TotalPrice = price
}
