package handlers

import (
	"Taller2/Product/internal/query/domain/models"
	"Taller2/Product/internal/query/domain/repositories"
	"context"
)

type GetProductHandler struct {
	repo repositories.ProductReadRepository
}

func NewGetProductHandler(repo repositories.ProductReadRepository) *GetProductHandler {
	return &GetProductHandler{repo: repo}
}

func (h *GetProductHandler) Handle(ctx context.Context, id string) (*models.ProductRead, error) {
	// Call the repository to find the product by ID
	product, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Return the product read model
	return product, nil
}

// para el list

type ListProductsQuery struct {
	Page  int
	Limit int
}

type ListProductsHandler struct {
	repo repositories.ProductReadRepository
}

func NewListProductsHandler(repo repositories.ProductReadRepository) *ListProductsHandler {
	return &ListProductsHandler{repo: repo}
}

func (h *ListProductsHandler) Handle(ctx context.Context, query ListProductsQuery) ([]*models.ProductRead, int, error) {
	offset := (query.Page - 1) * query.Limit
	return h.repo.ListWithPagination(ctx, query.Limit, offset)
}
