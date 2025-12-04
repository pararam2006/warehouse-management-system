package repositories

import (
	"context"
	"database/sql"
	//"fmt"
	"time"
	"warehouse-management-system/src/models"
)

// WarehouseRepositorySQLite — реализация складского хранилища на SQLite.
// Хранит только таблицу движений stock_movements, остатки считаются на лету агрегированными запросами.
type WarehouseRepositorySQLite struct {
	db *sql.DB
}

// NewWarehouseRepository создаёт новый репозиторий складских данных.
func NewWarehouseRepository(db *sql.DB) *WarehouseRepositorySQLite {
	return &WarehouseRepositorySQLite{db: db}
}

// GetInventory возвращает текущие остатки по всем товарам.
func (r *WarehouseRepositorySQLite) GetInventory() ([]*models.StockItem, error) {
	const query = `
SELECT
    product_id,
    SUM(
        CASE type
            WHEN 'receipt'  THEN quantity
            WHEN 'write_off' THEN -quantity
            WHEN 'reserve'  THEN -quantity
        END
    ) AS quantity
FROM stock_movements
GROUP BY product_id;
`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.StockItem
	for rows.Next() {
		var it models.StockItem
		if err := rows.Scan(&it.ProductID, &it.Quantity); err != nil {
			return nil, err
		}
		result = append(result, &it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetStockByProduct возвращает текущий остаток по конкретному товару.
func (r *WarehouseRepositorySQLite) GetStockByProduct(productID string) (float64, error) {
	const query = `
SELECT
    COALESCE(SUM(
        CASE type
            WHEN 'receipt'  THEN quantity
            WHEN 'write_off' THEN -quantity
            WHEN 'reserve'  THEN -quantity
        END
    ), 0) AS quantity
FROM stock_movements
WHERE product_id = ?;
`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var qty float64
	if err := r.db.QueryRowContext(ctx, query, productID).Scan(&qty); err != nil {
		return 0, err
	}
	return qty, nil
}

// AddMovement добавляет движение товара.
func (r *WarehouseRepositorySQLite) AddMovement(m *models.StockMovement) error {
	const query = `
INSERT INTO stock_movements (id, type, product_id, supplier_id, order_id, quantity, price, expiry_date, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
`

	if m.ID == "" {
		m.ID = "w-" + time.Now().UTC().Format("20060102T150405.000000000")
	}
	now := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		m.ID,
		m.Type,
		m.ProductID,
		m.SupplierID,
		m.OrderID,
		m.Quantity,
		m.Price,
		m.ExpiryDate,
		m.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
