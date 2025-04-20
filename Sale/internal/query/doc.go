package query

import (
	handlers "Taller2/Sale/internal/query/application/handlers"
	repositories "Taller2/Sale/internal/query/domain/repositories"
)

// Application contiene los casos de uso de Query para Sale
type Application struct {
	ListSales *handlers.ListSalesHandler
	GetSale   *handlers.GetSaleHandler
	// ...otros handlers
}

func NewApplication(repo repositories.SaleReadRepository) *Application {
	return &Application{
		ListSales: handlers.NewListSalesHandler(repo),
		GetSale:   handlers.NewGetSaleHandler(repo),
	}
}
