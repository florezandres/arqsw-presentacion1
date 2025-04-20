// internal/command/command.go
package command

import (
	handlers "Taller2/Sale/internal/command/application/handlers"
	"Taller2/Sale/internal/command/domain/repositories"
	"Taller2/Sale/internal/command/infrastructure/kafka"
)

type Application struct {
	CreateSale *handlers.CreateSaleHandler
	DeleteSale *handlers.DeleteSaleHandler
	UpdateSale *handlers.UpdateSaleHandler
}

// Usa kafka.EventProducer como tipo
func NewApplication(repo repositories.SaleRepository, producer kafka.EventProducer) *Application {
	return &Application{
		CreateSale: handlers.NewCreateSaleHandler(repo, producer),
		DeleteSale: handlers.NewDeleteSaleHandler(repo, producer),
		UpdateSale: handlers.NewUpdateSaleHandler(repo, producer),
	}
}
