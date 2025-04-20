package handlers

import (
	"Taller2/Sale/internal/query/domain/models"
	"Taller2/Sale/internal/query/domain/repositories"
	"context"
)

type GetSaleHandler struct {
	repo repositories.SaleReadRepository
}

func NewGetSaleHandler(repo repositories.SaleReadRepository) *GetSaleHandler {
	return &GetSaleHandler{repo: repo}
}

func (h *GetSaleHandler) Handle(ctx context.Context, id string) (*models.SaleRead, error) {
	// Llama al repositorio para encontrar la venta por ID
	sale, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Retorna el modelo de lectura
	return sale, nil
}

// Para listar ventas

type ListSalesQuery struct {
	Page  int
	Limit int
}

type ListSalesHandler struct {
	repo repositories.SaleReadRepository
}

func NewListSalesHandler(repo repositories.SaleReadRepository) *ListSalesHandler {
	return &ListSalesHandler{repo: repo}
}

func (h *ListSalesHandler) Handle(ctx context.Context, query ListSalesQuery) ([]*models.SaleRead, int, error) {
	offset := (query.Page - 1) * query.Limit
	return h.repo.ListWithPagination(ctx, query.Limit, offset)
}
