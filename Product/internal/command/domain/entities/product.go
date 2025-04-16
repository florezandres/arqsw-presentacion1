package entities

import (
	"errors"
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Stock       int
}

func NewProduct(name, description string, price float64, stock int) (*Product, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("price must be positive")
	}

	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}, nil
}

// Métodos de negocio (ej: reducir stock)
func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return errors.New("insufficient stock")
	}
	p.Stock -= quantity
	return nil
}
