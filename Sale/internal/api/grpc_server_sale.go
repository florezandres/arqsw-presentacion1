package api

import (
	gen "Taller2/Sale/gen/go/api"
	"Taller2/Sale/internal/command"
	commands "Taller2/Sale/internal/command/application/commands"
	handlersCommand "Taller2/Sale/internal/command/application/handlers"
	"Taller2/Sale/internal/query"
	handlersQuery "Taller2/Sale/internal/query/application/handlers"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SaleServer struct {
	gen.UnimplementedSalesServiceServer
	commandApp *command.Application
	queryApp   *query.Application
}

func NewSaleServer(commandApp *command.Application, queryApp *query.Application) *SaleServer {
	return &SaleServer{
		commandApp: commandApp,
		queryApp:   queryApp,
	}
}

func (s *SaleServer) CreateSale(ctx context.Context, req *gen.CreateSaleRequest) (*gen.CreateSaleResponse, error) {
	cmd := commands.CreateSaleCommand{
		ProductID: req.ProductId,
		Quantity:  int(req.Quantity),
	}

	saleID, err := s.commandApp.CreateSale.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &gen.CreateSaleResponse{SaleId: saleID}, nil
}

func (s *SaleServer) GetSale(ctx context.Context, req *gen.GetSaleRequest) (*gen.GetSaleResponse, error) {
	sale, err := s.queryApp.GetSale.Handle(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "sale not found: %v", err)
	}

	return &gen.GetSaleResponse{
		Sale: &gen.Sale{
			Id:        sale.ID,
			ProductId: sale.ProductID,
			Quantity:  int32(sale.Quantity),
			SaleDate:  sale.Date,
		},
	}, nil
}

func (s *SaleServer) ListSales(ctx context.Context, req *gen.ListSalesRequest) (*gen.ListSalesResponse, error) {
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
	sales, total, err := s.queryApp.ListSales.Handle(ctx, handlersQuery.ListSalesQuery{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list sales: %v", err)
	}

	// Convierte a formato gRPC
	var pbSales []*gen.Sale
	for _, s := range sales {
		pbSales = append(pbSales, &gen.Sale{
			Id:        s.ID,
			ProductId: s.ProductID,
			Quantity:  int32(s.Quantity),
			SaleDate:  s.Date,
		})
	}

	return &gen.ListSalesResponse{
		Sales: pbSales,
		Total: int32(total),
	}, nil
}

func (s *SaleServer) DeleteSale(ctx context.Context, req *gen.DeleteSaleRequest) (*gen.DeleteSaleResponse, error) {
	err := s.commandApp.DeleteSale.Handle(ctx, handlersCommand.DeleteSaleCommand{
		ID: req.GetId(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete sale: %v", err)
	}

	return &gen.DeleteSaleResponse{Success: true}, nil
}

func (s *SaleServer) UpdateSale(ctx context.Context, req *gen.UpdateSaleRequest) (*gen.UpdateSaleResponse, error) {
	err := s.commandApp.UpdateSale.Handle(ctx, handlersCommand.UpdateSaleCommand{
		ID:        req.GetId(),
		ProductID: req.GetProductId(),
		Quantity:  int(req.GetQuantity()),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update sale: %v", err)
	}

	return &gen.UpdateSaleResponse{Success: true}, nil
}
