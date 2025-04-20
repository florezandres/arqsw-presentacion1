package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Sale struct {
	ID         uuid.UUID
	ProductID  string
	Quantity   int
	TotalPrice float64
	SaleDate   string
}

func NewSale(productID string, quantity int, totalPrice float64, saleDate string) (*Sale, error) {
	if productID == "" {
		return nil, errors.New("productID cannot be empty")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if totalPrice <= 0 {
		return nil, errors.New("totalPrice must be positive")
	}
	if saleDate == "" {
		return nil, errors.New("saleDate cannot be empty")
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
