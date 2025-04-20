package handlers

import (
	"Taller2/Sale/internal/command/application/commands"
	"Taller2/Sale/internal/command/domain/entities"
	"Taller2/Sale/internal/command/domain/events"
	"Taller2/Sale/internal/command/domain/repositories"
	"Taller2/Sale/internal/command/infrastructure/kafka"
	"context"

	"github.com/google/uuid"
)

// CreateSaleHandler
type CreateSaleHandler struct {
	repo     repositories.SaleRepository
	producer kafka.EventProducer
}

func NewCreateSaleHandler(repo repositories.SaleRepository, producer kafka.EventProducer) *CreateSaleHandler {
	return &CreateSaleHandler{repo: repo, producer: producer}
}

func (h *CreateSaleHandler) Handle(ctx context.Context, cmd commands.CreateSaleCommand) (string, error) {
	sale, err := entities.NewSale(cmd.ProductID, cmd.Quantity, cmd.TotalPrice, cmd.SaleDate)
	if err != nil {
		return "", err
	}

	if err := h.repo.Save(ctx, sale); err != nil {
		return "", err
	}

	event := events.NewSaleCreatedEvent(sale)
	if err := h.producer.PublishEvent(ctx, event); err != nil {
		return "", err
	}

	return sale.ID.String(), nil
}

// DeleteSaleHandler
type DeleteSaleCommand struct {
	ID string
}

type DeleteSaleHandler struct {
	repo     repositories.SaleRepository
	producer kafka.EventProducer
}

func NewDeleteSaleHandler(repo repositories.SaleRepository, producer kafka.EventProducer) *DeleteSaleHandler {
	return &DeleteSaleHandler{repo: repo, producer: producer}
}

func (h *DeleteSaleHandler) Handle(ctx context.Context, cmd DeleteSaleCommand) error {
	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	event := events.NewSaleDeletedEvent(cmd.ID)
	return h.producer.PublishEvent(ctx, event)
}

// UpdateSaleHandler
type UpdateSaleCommand struct {
	ID         string
	ProductID  string
	Quantity   int
	TotalPrice float64
	SaleDate   string
}

type UpdateSaleHandler struct {
	repo     repositories.SaleRepository
	producer kafka.EventProducer
}

func NewUpdateSaleHandler(repo repositories.SaleRepository, producer kafka.EventProducer) *UpdateSaleHandler {
	return &UpdateSaleHandler{repo: repo, producer: producer}
}

func (h *UpdateSaleHandler) Handle(ctx context.Context, cmd UpdateSaleCommand) error {
	sale := &entities.Sale{
		ID:         uuid.MustParse(cmd.ID),
		ProductID:  cmd.ProductID,
		Quantity:   cmd.Quantity,
		TotalPrice: cmd.TotalPrice,
		SaleDate:   cmd.SaleDate,
	}

	if err := h.repo.Update(ctx, sale); err != nil {
		return err
	}

	event := events.NewSaleUpdatedEvent(sale)
	return h.producer.PublishEvent(ctx, event)
}
