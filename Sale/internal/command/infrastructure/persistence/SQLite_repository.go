package persistence

import (
	"Taller2/Sale/internal/command/domain/entities"
	"context"
	"database/sql"
	"errors"
	"time"
)

type SQLiteSaleRepository struct {
	db *sql.DB
}

func (r *SQLiteSaleRepository) Save(ctx context.Context, sale *entities.Sale) error {
	//TODO implement me
	panic("implement me")
}

func NewSQLiteSaleRepository(db *sql.DB) *SQLiteSaleRepository {
	return &SQLiteSaleRepository{db: db}
}

func (r *SQLiteSaleRepository) Create(ctx context.Context, sale *entities.Sale) error {
	// Convertir time.Time a string en formato RFC3339
	saleDateStr := sale.SaleDate.Format(time.RFC3339)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sales (id, product_id, quantity, total_price, sale_date) 
		VALUES (?, ?, ?, ?, ?)`,
		sale.ID,
		sale.ProductID,
		sale.Quantity,
		sale.TotalPrice,
		saleDateStr, // Usar el string formateado
	)
	return err
}

func (r *SQLiteSaleRepository) Update(ctx context.Context, sale *entities.Sale) error {
	saleDateStr := sale.SaleDate.Format(time.RFC3339)

	result, err := r.db.ExecContext(ctx, `
		UPDATE sales 
		SET product_id = ?, quantity = ?, total_price = ?, sale_date = ?
		WHERE id = ?`,
		sale.ProductID,
		sale.Quantity,
		sale.TotalPrice,
		saleDateStr,
		sale.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("sale not found")
	}

	return nil
}

func (r *SQLiteSaleRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM sales WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("sale not found")
	}

	return nil
}

func RunSQLiteMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sales (
			id TEXT PRIMARY KEY,
			product_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			total_price REAL NOT NULL,
			sale_date TEXT NOT NULL
		)`)
	return err
}
