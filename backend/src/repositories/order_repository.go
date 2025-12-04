package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"warehouse-management-system/src/models"
)

// OrderRepositorySQLite — реализация хранилища заказов на SQLite.
// Использует таблицы orders, order_items и order_status_history.
type OrderRepositorySQLite struct {
	db *sql.DB
}

// NewOrderRepository создаёт новый репозиторий заказов.
func NewOrderRepository(db *sql.DB) *OrderRepositorySQLite {
	return &OrderRepositorySQLite{db: db}
}

// GetAll возвращает все заказы.
func (r *OrderRepositorySQLite) GetAll() ([]*models.Order, error) {
	const queryOrders = `
SELECT id, customer, status, created_at, updated_at
FROM orders
ORDER BY created_at DESC;
`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, queryOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.Customer, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		// Подгружаем позиции и историю статусов.
		if err := r.loadItemsAndHistory(ctx, &o); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByID возвращает заказ по ID.
func (r *OrderRepositorySQLite) GetByID(id string) (*models.Order, error) {
	const query = `
SELECT id, customer, status, created_at, updated_at
FROM orders
WHERE id = ? LIMIT 1;
`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id)
	var o models.Order
	if err := row.Scan(&o.ID, &o.Customer, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := r.loadItemsAndHistory(ctx, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// Create сохраняет новый заказ, его позиции и историю статусов.
func (r *OrderRepositorySQLite) Create(order *models.Order) error {
	if order.ID == "" {
		order.ID = "o-" + time.Now().UTC().Format("20060102T150405.000000000")
	}
	now := time.Now().UTC()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	order.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const insertOrder = `
INSERT INTO orders (id, customer, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);
`
	if _, err = tx.ExecContext(ctx, insertOrder,
		order.ID,
		order.Customer,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	); err != nil {
		return err
	}

	const insertItem = `
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES (?, ?, ?, ?);
`
	for _, it := range order.Items {
		if _, err = tx.ExecContext(ctx, insertItem,
			order.ID,
			it.ProductID,
			it.Quantity,
			it.Price,
		); err != nil {
			return err
		}
	}

	const insertHist = `
INSERT INTO order_status_history (order_id, status, changed_at)
VALUES (?, ?, ?);
`
	for _, h := range order.StatusHist {
		if _, err = tx.ExecContext(ctx, insertHist,
			order.ID,
			h.Status,
			h.ChangedAt,
		); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// Update обновляет существующий заказ (статус, updated_at и историю статусов).
func (r *OrderRepositorySQLite) Update(order *models.Order) error {
	now := time.Now().UTC()
	order.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const updateOrder = `
UPDATE orders
SET customer = ?, status = ?, updated_at = ?
WHERE id = ?;
`
	res, err := tx.ExecContext(ctx, updateOrder,
		order.Customer,
		order.Status,
		order.UpdatedAt,
		order.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("order with id %s not found", order.ID)
	}

	// Добавляем последнюю запись истории статусов (если есть).
	if len(order.StatusHist) > 0 {
		last := order.StatusHist[len(order.StatusHist)-1]
		const insertHist = `
INSERT INTO order_status_history (order_id, status, changed_at)
VALUES (?, ?, ?);
`
		if _, err = tx.ExecContext(ctx, insertHist,
			order.ID,
			last.Status,
			last.ChangedAt,
		); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// loadItemsAndHistory подгружает позиции и историю статусов заказа.
func (r *OrderRepositorySQLite) loadItemsAndHistory(ctx context.Context, o *models.Order) error {
	const queryItems = `
SELECT product_id, quantity, price
FROM order_items
WHERE order_id = ?;
`
	rows, err := r.db.QueryContext(ctx, queryItems, o.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var it models.OrderItem
		if err := rows.Scan(&it.ProductID, &it.Quantity, &it.Price); err != nil {
			return err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	o.Items = items

	const queryHist = `
SELECT status, changed_at
FROM order_status_history
WHERE order_id = ?
ORDER BY changed_at;
`
	hRows, err := r.db.QueryContext(ctx, queryHist, o.ID)
	if err != nil {
		return err
	}
	defer hRows.Close()

	var hist []models.StatusEntry
	for hRows.Next() {
		var h models.StatusEntry
		if err := hRows.Scan(&h.Status, &h.ChangedAt); err != nil {
			return err
		}
		hist = append(hist, h)
	}
	if err := hRows.Err(); err != nil {
		return err
	}
	o.StatusHist = hist

	return nil
}
