package query

//EL PROBLEMA ES QUE ESTA VAINA SE LLAMA ES APP NO DOC JAJA
import (
	handlers "Taller2/Product/internal/query/application/handlers"
	repositories "Taller2/Product/internal/query/domain/repositories"
)

// Application contiene los casos de uso de Query
type Application struct {
	ListProducts *handlers.ListProductsHandler
	GetProduct   *handlers.GetProductHandler
	// ...otros handlers
}

func NewApplication(repo repositories.ProductReadRepository) *Application {
	return &Application{
		ListProducts: handlers.NewListProductsHandler(repo),
		GetProduct:   handlers.NewGetProductHandler(repo),
	}
}
