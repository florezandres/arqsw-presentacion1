package persistence

import (
	"Taller2/Sale/internal/command/domain/entities"
	"Taller2/Sale/internal/command/domain/repositories"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Importar el driver de SQLite
)

type SQLiteSaleRepository struct {
	db *sql.DB
}

func NewSQLiteSaleRepository(dsn string) (repositories.SaleRepository, error) {
	db, err := sql.Open("sqlite3", dsn) // Cambiar a "sqlite3" como el driver
	if err != nil {
		return nil, err
	}
	return &SQLiteSaleRepository{db: db}, nil
}

func (r *SQLiteSaleRepository) Save(ctx context.Context, sale *entities.Sale) error {
	query := `INSERT INTO sales (id, product_id, quantity, total_price, sale_date) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, sale.ID, sale.ProductID, sale.Quantity, sale.TotalPrice, sale.SaleDate)
	return err
}

// Delete
func (r *SQLiteSaleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sales WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sale not found")
	}

	return nil
}

// Update
func (r *SQLiteSaleRepository) Update(ctx context.Context, sale *entities.Sale) error {
	query := `
        UPDATE sales
        SET product_id = ?, quantity = ?, total_price = ?, sale_date = ?
        WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		sale.ProductID,
		sale.Quantity,
		sale.TotalPrice,
		sale.SaleDate,
		sale.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sale not found")
	}

	return nil
}
