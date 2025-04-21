package persistence

import (
	domain "Taller2/Sale/internal/query/domain/models"
	"database/sql"
	"errors"
	"time"
)

type SQLiteSaleReadRepository struct {
	db *sql.DB
}

func NewSQLiteSaleReadRepository(db *sql.DB) *SQLiteSaleReadRepository {
	return &SQLiteSaleReadRepository{db: db}
}

func (r *SQLiteSaleReadRepository) Create(sale map[string]interface{}) error {
	saleDate, err := time.Parse(time.RFC3339, sale["sale_date"].(string))
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO sales (id, product_id, quantity, total_price, sale_date) 
		VALUES (?, ?, ?, ?, ?)`,
		sale["id"],
		sale["product_id"],
		sale["quantity"],
		sale["total_price"],
		saleDate,
	)
	return err
}

func (r *SQLiteSaleReadRepository) Update(sale map[string]interface{}) error {
	saleDate, err := time.Parse(time.RFC3339, sale["sale_date"].(string))
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		UPDATE sales 
		SET product_id = ?, quantity = ?, total_price = ?, sale_date = ?
		WHERE id = ?`,
		sale["product_id"],
		sale["quantity"],
		sale["total_price"],
		saleDate,
		sale["id"],
	)
	return err
}

func (r *SQLiteSaleReadRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM sales WHERE id = ?", id)
	return err
}

func (r *SQLiteSaleReadRepository) GetByID(id string) (*domain.SaleRead, error) {
	row := r.db.QueryRow("SELECT id, product_id, quantity, total_price, sale_date FROM sales WHERE id = ?", id)

	var sale domain.SaleRead
	var saleDateStr string
	err := row.Scan(
		&sale.ID,
		&sale.ProductID,
		&sale.Quantity,
		&saleDateStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &sale, nil
}

func (r *SQLiteSaleReadRepository) List(page, limit int) ([]*domain.SaleRead, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM sales").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	rows, err := r.db.Query(`
		SELECT id, product_id, quantity, total_price, sale_date 
		FROM sales 
		ORDER BY sale_date DESC 
		LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sales []*domain.SaleRead
	for rows.Next() {
		var sale domain.SaleRead
		var saleDateStr string
		err := rows.Scan(
			&sale.ID,
			&sale.ProductID,
			&sale.Quantity,
			&saleDateStr,
		)
		if err != nil {
			return nil, 0, err
		}

		sales = append(sales, &sale)
	}

	return sales, total, nil
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
