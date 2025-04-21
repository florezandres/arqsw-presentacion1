package handlers

import (
	"Taller2/Sale/internal/command/application/commands"
	"Taller2/Sale/internal/command/domain/entities"
	"Taller2/Sale/internal/command/domain/events"
	"Taller2/Sale/internal/command/domain/repositories"
	"Taller2/Sale/internal/command/infrastructure/kafka"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
	// Validación de campos requeridos
	if cmd.ProductID == "" {
		return "", errors.New("el product_id es requerido")
	}
	if cmd.Quantity <= 0 {
		return "", errors.New("la cantidad debe ser mayor a cero")
	}

	// Parseo de fecha con la nueva función
	saleDate, err := parseSaleDate(cmd.SaleDate)
	if err != nil {
		return "", fmt.Errorf("error en fecha: %v", err)
	}

	// Creación de la entidad Sale
	sale, err := entities.NewSale(
		cmd.ProductID,
		cmd.Quantity,
		cmd.TotalPrice,
		saleDate, // time.Time ya validado
	)
	if err != nil {
		return "", fmt.Errorf("error creando venta: %v", err)
	}

	// Guardar en repositorio
	if err := h.repo.Save(ctx, sale); err != nil {
		return "", fmt.Errorf("error guardando venta: %v", err)
	}

	// Publicar evento
	event := events.NewSaleCreatedEvent(sale)
	if err := h.producer.PublishEvent(ctx, event); err != nil {
		return "", fmt.Errorf("error publicando evento: %v", err)
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
	// Parsear la fecha de string a time.Time
	saleDate, err := time.Parse("2006-01-02", cmd.SaleDate)
	if err != nil {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}

	sale := &entities.Sale{
		ID:         uuid.MustParse(cmd.ID),
		ProductID:  cmd.ProductID,
		Quantity:   cmd.Quantity,
		TotalPrice: cmd.TotalPrice,
		SaleDate:   saleDate, // Usar el time.Time parseado
	}

	if err := h.repo.Update(ctx, sale); err != nil {
		return err
	}

	event := events.NewSaleUpdatedEvent(sale)
	return h.producer.PublishEvent(ctx, event)
}

// parseSaleDate intenta múltiples formatos de fecha
func parseSaleDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, errors.New("la fecha no puede estar vacía")
	}

	formats := []string{
		"2006-01-02",       // Formato YYYY-MM-DD
		"2006-01-02 15:04", // Con hora
		time.RFC3339,       // Formato ISO8601
		"02/01/2006",       // Formato DD/MM/YYYY
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("formato de fecha no reconocido. Use YYYY-MM-DD")
}
