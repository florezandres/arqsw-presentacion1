package persistence

import (
	"Taller2/Product/internal/command/domain/entities"
	"Taller2/Product/internal/command/domain/repositories"
	"context"
	"database/sql"
	"fmt"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(dsn string) (repositories.ProductRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresProductRepository{db: db}, nil
}

func (r *PostgresProductRepository) Save(ctx context.Context, product *entities.Product) error {
	query := `INSERT INTO products_command (id, name, price, stock) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Price, product.Stock)
	return err
}

// Delete
func (r *PostgresProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products_command WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Update
func (r *PostgresProductRepository) Update(ctx context.Context, product *entities.Product) error {
	query := `
        UPDATE products_command 
        SET name = $1, description = $2, price = $3, stock = $4 
        WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
