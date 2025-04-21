package persistence

import (
	"Taller2/Sale/internal/query/domain/models"
	"context"
	"database/sql"
	"errors"
)

type SaleRead struct {
	ID        string
	ProductID string
	Quantity  int32
	Date      string
}

type SQLiteSaleReadRepository struct {
	db *sql.DB
}

func (r *SQLiteSaleReadRepository) FindByID(ctx context.Context, id string) (*models.SaleRead, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SQLiteSaleReadRepository) ListAll(ctx context.Context) ([]*models.SaleRead, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SQLiteSaleReadRepository) Save(ctx context.Context, sale *models.SaleRead) error {
	//TODO implement me
	panic("implement me")
}

func (r *SQLiteSaleReadRepository) ListWithPagination(ctx context.Context, limit, offset int) ([]*models.SaleRead, int, error) {
	//TODO implement me
	panic("implement me")
}

func NewSQLiteSaleReadRepository(db *sql.DB) *SQLiteSaleReadRepository {
	return &SQLiteSaleReadRepository{db: db}
}

func (r *SQLiteSaleReadRepository) Create(ctx context.Context, sale *models.SaleRead) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sales (id, product_id, quantity, date) 
		VALUES (?, ?, ?, ?)`,
		sale.ID,
		sale.ProductID,
		sale.Quantity,
		sale.Date,
	)
	return err
}

func (r *SQLiteSaleReadRepository) Update(ctx context.Context, sale *models.SaleRead) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE sales 
		SET product_id = ?, quantity = ?, date = ?
		WHERE id = ?`,
		sale.ProductID,
		sale.Quantity,
		sale.Date,
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

func (r *SQLiteSaleReadRepository) Delete(ctx context.Context, id string) error {
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

func (r *SQLiteSaleReadRepository) GetByID(ctx context.Context, id string) (*SaleRead, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, product_id, quantity, date FROM sales WHERE id = ?", id)

	var sale SaleRead
	err := row.Scan(
		&sale.ID,
		&sale.ProductID,
		&sale.Quantity,
		&sale.Date,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &sale, nil
}

func (r *SQLiteSaleReadRepository) List(ctx context.Context, page, limit int) ([]*SaleRead, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sales").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, product_id, quantity, date 
		FROM sales 
		ORDER BY date DESC 
		LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sales []*SaleRead
	for rows.Next() {
		var sale SaleRead
		err := rows.Scan(
			&sale.ID,
			&sale.ProductID,
			&sale.Quantity,
			&sale.Date,
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
			date TEXT NOT NULL
		)`)
	return err
}
