package handlers

import (
	"Taller2/Product/internal/command/application/commands"
	"Taller2/Product/internal/command/domain/entities"
	"Taller2/Product/internal/command/domain/events"
	"Taller2/Product/internal/command/domain/repositories"
	"Taller2/Product/internal/command/infrastructure/kafka"
	"context"
	"github.com/google/uuid"
)

type CreateProductHandler struct {
	repo     repositories.ProductRepository
	producer kafka.EventProducer
}

func NewCreateProductHandler(repo repositories.ProductRepository, producer kafka.EventProducer) *CreateProductHandler {
	return &CreateProductHandler{repo: repo, producer: producer}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd commands.CreateProductCommand) (string, error) {
	product, err := entities.NewProduct(cmd.Name, cmd.Description, cmd.Price, cmd.Stock)
	if err != nil {
		return "", err
	}

	if err := h.repo.Save(ctx, product); err != nil {
		return "", err
	}

	event := events.NewProductCreatedEvent(product)
	if err := h.producer.PublishEvent(ctx, event); err != nil {
		return "", err
	}

	return product.ID.String(), nil
}

//delete

type DeleteProductCommand struct {
	ID string
}

type DeleteProductHandler struct {
	repo     repositories.ProductRepository
	producer kafka.EventProducer
}

func NewDeleteProductHandler(repo repositories.ProductRepository, producer kafka.EventProducer) *DeleteProductHandler {
	return &DeleteProductHandler{
		repo:     repo,
		producer: producer,
	}
}

// internal/command/application/handlers/delete_product_handler.go
func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	event := events.NewProductDeletedEvent(cmd.ID)
	return h.producer.PublishEvent(ctx, event)
}

// Update
type UpdateProductCommand struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
}

type UpdateProductHandler struct {
	repo     repositories.ProductRepository
	producer kafka.EventProducer
}

func NewUpdateProductHandler(repo repositories.ProductRepository, producer kafka.EventProducer) *UpdateProductHandler {
	return &UpdateProductHandler{
		repo:     repo,
		producer: producer,
	}
}

func (h *UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProductCommand) error {
	// 1. Actualizar en la base de datos
	product := &entities.Product{
		ID:          uuid.MustParse(cmd.ID),
		Name:        cmd.Name,
		Description: cmd.Description,
		Price:       cmd.Price,
		Stock:       cmd.Stock,
	}

	if err := h.repo.Update(ctx, product); err != nil {
		return err
	}

	// 2. Publicar evento
	event := events.NewProductUpdatedEvent(product)
	return h.producer.PublishEvent(ctx, event)
}
