package persistence

import (
	"Taller2/Product/internal/query/domain/models"
	"Taller2/Product/internal/query/domain/repositories"
	"context"
	"database/sql"
	"fmt"
)

type PostgresProductReadRepository struct {
	db *sql.DB
}

func (r *PostgresProductReadRepository) Save(ctx context.Context, product *models.ProductRead) error {
	query := `INSERT INTO products_query (id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Description, product.Price, product.Stock)
	return err
}

func (r *PostgresProductReadRepository) ListAll(ctx context.Context) ([]*models.ProductRead, error) {
	query := `SELECT id, name, description, price, stock FROM products_query`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.ProductRead
	for rows.Next() {
		var product models.ProductRead
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func NewPostgresProductReadRepository(dsn string) (repositories.ProductReadRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresProductReadRepository{db: db}, nil
}

func (r *PostgresProductReadRepository) FindByID(ctx context.Context, id string) (*models.ProductRead, error) {
	query := `SELECT id, name, price, stock FROM products_query WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var product models.ProductRead
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
		return nil, err
	}
	return &product, nil
}

// Implementar ListAll...

func (r *PostgresProductReadRepository) ListWithPagination(ctx context.Context, limit, offset int) ([]*models.ProductRead, int, error) {
	// 1. Obtener conteo total
	var total int
	countQuery := `SELECT COUNT(*) FROM products_query`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("error counting products: %w", err)
	}

	// 2. Obtener productos paginados
	query := `
        SELECT id, name, description, price, stock 
        FROM products_query 
        ORDER BY name ASC
        LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []*models.ProductRead
	for rows.Next() {
		var p models.ProductRead
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, 0, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, &p)
	}

	return products, total, nil
}

func (r *PostgresProductReadRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM products_query WHERE id = $1", id)
	return err
}

// Update
func (r *PostgresProductReadRepository) Update(ctx context.Context, product *models.ProductRead) error {
	query := `
        UPDATE products_query 
        SET name = $1, description = $2, price = $3, stock = $4 
        WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ID)

	return err
}
