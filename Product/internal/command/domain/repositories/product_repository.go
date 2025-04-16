package repositories

import (
	"Taller2/Product/internal/command/domain/entities"
	"context"
)

type ProductRepository interface {
	Save(ctx context.Context, product *entities.Product) error
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id string) error
}
