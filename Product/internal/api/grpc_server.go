package api

// creo que se puede mejorar la forma en la que se manejan los imports

import (
	gen "Taller2/Product/gen/go/api" // Importación ABSOLUTA desde el módulo
	"Taller2/Product/internal/command"
	commands "Taller2/Product/internal/command/application/commands"
	handlersCommand "Taller2/Product/internal/command/application/handlers"
	"Taller2/Product/internal/query"
	handlersQuery "Taller2/Product/internal/query/application/handlers"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "log"
)

type ProductServer struct {
	gen.UnimplementedProductServiceServer
	commandApp *command.Application
	queryApp   *query.Application
}

func NewProductServer(commandApp *command.Application, queryApp *query.Application) *ProductServer {
	return &ProductServer{
		commandApp: commandApp,
		queryApp:   queryApp,
	}
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.CreateProductResponse, error) {
	// Asegúrate que command.CreateProductCommand esté correctamente definido
	cmd := commands.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
	}

	// El método Handle debería devolver solo (string, error) o ajustar la asignación
	productID, err := s.commandApp.CreateProduct.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &gen.CreateProductResponse{ProductId: productID}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.GetProductResponse, error) {
	// Asegúrate que query.GetProductHandler.Handle devuelva un tipo compatible
	productData, err := s.queryApp.GetProduct.Handle(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Convierte productData al tipo esperado por gen.GetProductResponse
	return &gen.GetProductResponse{
		Product: &gen.Product{
			Id:          productData.ID,
			Name:        productData.Name,
			Description: productData.Description,
			Price:       productData.Price,
			Stock:       productData.Stock,
		},
	}, nil
}

// ListProducts maneja la solicitud de listar productos
func (s *ProductServer) ListProducts(ctx context.Context, req *gen.ListProductsRequest) (*gen.ListProductsResponse, error) {
	// Valores por defecto
	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}

	limit := int(req.GetLimit())
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Llama al caso de uso de Query
	products, total, err := s.queryApp.ListProducts.Handle(ctx, handlersQuery.ListProductsQuery{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	// Convierte a formato gRPC
	var pbProducts []*gen.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &gen.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
		})
	}

	return &gen.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(total),
	}, nil
}

// DeleteProduct maneja la solicitud de eliminar un producto
func (s *ProductServer) DeleteProduct(ctx context.Context, req *gen.DeleteProductRequest) (*gen.DeleteProductResponse, error) {
	err := s.commandApp.DeleteProduct.Handle(ctx, handlersCommand.DeleteProductCommand{
		ID: req.GetId(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &gen.DeleteProductResponse{Success: true}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *gen.UpdateProductRequest) (*gen.UpdateProductResponse, error) {
	err := s.commandApp.UpdateProduct.Handle(ctx, handlersCommand.UpdateProductCommand{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &gen.UpdateProductResponse{Success: true}, nil
}
