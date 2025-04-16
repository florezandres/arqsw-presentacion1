// internal/command/command.go
package command

import (
	handlers "Taller2/Product/internal/command/application/handlers"
	"Taller2/Product/internal/command/domain/repositories"
	"Taller2/Product/internal/command/infrastructure/kafka"
)

type Application struct {
	CreateProduct *handlers.CreateProductHandler
	DeleteProduct *handlers.DeleteProductHandler
	UpdateProduct *handlers.UpdateProductHandler
}

// Usa kafka.EventProducer como tipo
func NewApplication(repo repositories.ProductRepository, producer kafka.EventProducer) *Application {
	return &Application{
		CreateProduct: handlers.NewCreateProductHandler(repo, producer),
		DeleteProduct: handlers.NewDeleteProductHandler(repo, producer),
		UpdateProduct: handlers.NewUpdateProductHandler(repo, producer),
	}
}
