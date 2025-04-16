package repositories

import (
	"Taller2/Product/internal/query/domain/models"
	"context"
)

type ProductReadRepository interface {
	FindByID(ctx context.Context, id string) (*models.ProductRead, error)
	ListAll(ctx context.Context) ([]*models.ProductRead, error)
	Save(ctx context.Context, product *models.ProductRead) error
	ListWithPagination(ctx context.Context, limit, offset int) ([]*models.ProductRead, int, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, product *models.ProductRead) error
}
