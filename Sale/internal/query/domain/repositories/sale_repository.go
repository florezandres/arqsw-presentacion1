package repositories

import (
	"Taller2/Sale/internal/query/domain/models"
	"context"
)

type SaleReadRepository interface {
	FindByID(ctx context.Context, id string) (*models.SaleRead, error)
	ListAll(ctx context.Context) ([]*models.SaleRead, error)
	Save(ctx context.Context, sale *models.SaleRead) error
	ListWithPagination(ctx context.Context, limit, offset int) ([]*models.SaleRead, int, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, sale *models.SaleRead) error
}
