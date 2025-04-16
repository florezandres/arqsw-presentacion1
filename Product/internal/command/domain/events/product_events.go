package events

import (
	"Taller2/Product/internal/command/domain/entities"
	"encoding/json"
)

// ProductEventWrapper es la estructura genérica para eventos
type ProductEventWrapper struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// ProductEvent es la interfaz que deben implementar todos los eventos
type ProductEvent interface {
	ToJSON() ([]byte, error)
	EventType() string
}

// ProductCreatedEvent implementa ProductEvent
type ProductCreatedEvent struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func NewProductCreatedEvent(product *entities.Product) ProductEvent {
	return &ProductCreatedEvent{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}
}

func (e *ProductCreatedEvent) EventType() string {
	return "product_created"
}

func (e *ProductCreatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ProductDeletedEvent implementa ProductEvent
type ProductDeletedEvent struct {
	ProductID string `json:"product_id"`
}

func NewProductDeletedEvent(productID string) ProductEvent {
	return &ProductDeletedEvent{
		ProductID: productID,
	}
}

func (e *ProductDeletedEvent) EventType() string {
	return "product_deleted"
}

func (e *ProductDeletedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// Función de ayuda para crear el wrapper genérico (opcional)
func CreateEventWrapper(event ProductEvent) (*ProductEventWrapper, error) {
	data, err := event.ToJSON()
	if err != nil {
		return nil, err
	}

	var payload interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	return &ProductEventWrapper{
		Type:    event.EventType(),
		Payload: payload,
	}, nil
}

//Update

type ProductUpdatedEvent struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
}

func NewProductUpdatedEvent(product *entities.Product) ProductEvent {
	return &ProductUpdatedEvent{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}
}

func (e *ProductUpdatedEvent) EventType() string {
	return "product_updated"
}

func (e *ProductUpdatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
