package events

import (
	"Taller2/Sale/internal/command/domain/entities"
	"encoding/json"
)

// SaleEventWrapper es la estructura genérica para eventos de venta
type SaleEventWrapper struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// SaleEvent es la interfaz que deben implementar todos los eventos de venta
type SaleEvent interface {
	ToJSON() ([]byte, error)
	EventType() string
}

// SaleCreatedEvent representa el evento de creación de una venta
type SaleCreatedEvent struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"totalprice"`
}

func NewSaleCreatedEvent(sale *entities.Sale) SaleEvent {
	return &SaleCreatedEvent{
		ID:         sale.ID.String(),
		ProductID:  sale.ProductID,
		Quantity:   sale.Quantity,
		TotalPrice: sale.TotalPrice,
	}
}

func (e *SaleCreatedEvent) EventType() string {
	return "sale_created"
}

func (e *SaleCreatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// SaleDeletedEvent representa el evento de eliminación de una venta
type SaleDeletedEvent struct {
	SaleID string `json:"sale_id"`
}

func NewSaleDeletedEvent(saleID string) SaleEvent {
	return &SaleDeletedEvent{
		SaleID: saleID,
	}
}

func (e *SaleDeletedEvent) EventType() string {
	return "sale_deleted"
}

func (e *SaleDeletedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// SaleUpdatedEvent representa el evento de actualización de una venta
type SaleUpdatedEvent struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"totalprice"`
}

func NewSaleUpdatedEvent(sale *entities.Sale) SaleEvent {
	return &SaleUpdatedEvent{
		ID:         sale.ID.String(),
		ProductID:  sale.ProductID,
		Quantity:   sale.Quantity,
		TotalPrice: sale.TotalPrice,
	}
}

func (e *SaleUpdatedEvent) EventType() string {
	return "sale_updated"
}

func (e *SaleUpdatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// Función de ayuda para crear el wrapper genérico
func CreateSaleEventWrapper(event SaleEvent) (*SaleEventWrapper, error) {
	data, err := event.ToJSON()
	if err != nil {
		return nil, err
	}

	var payload interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	return &SaleEventWrapper{
		Type:    event.EventType(),
		Payload: payload,
	}, nil
}
