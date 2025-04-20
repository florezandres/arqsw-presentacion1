package repositories

import (
	"Taller2/Sale/internal/command/domain/entities"
	"context"
)

type SaleRepository interface {
	Save(ctx context.Context, sale *entities.Sale) error
	Update(ctx context.Context, sale *entities.Sale) error
	Delete(ctx context.Context, id string) error
}
