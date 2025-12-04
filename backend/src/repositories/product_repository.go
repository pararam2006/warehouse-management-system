package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"warehouse-management-system/src/models"
)

// ProductRepositorySQLite — реализация хранилища товаров на SQLite.
type ProductRepositorySQLite struct {
	db *sql.DB
}

// NewProductRepository создаёт новый репозиторий товаров.
func NewProductRepository(db *sql.DB) *ProductRepositorySQLite {
	return &ProductRepositorySQLite{db: db}
}

// GetAll возвращает список всех товаров.
func (r *ProductRepositorySQLite) GetAll() ([]*models.Product, error) {
	const query = `
SELECT id, sku, name, description, category_id, supplier_id, unit, created_at, updated_at
FROM products
ORDER BY created_at DESC;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Name,
			&p.Description,
			&p.CategoryID,
			&p.SupplierID,
			&p.Unit,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetByID возвращает товар по идентификатору.
func (r *ProductRepositorySQLite) GetByID(id string) (*models.Product, error) {
	const query = `
SELECT id, sku, name, description, category_id, supplier_id, unit, created_at, updated_at
FROM products
WHERE id = ? LIMIT 1;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id)
	var p models.Product
	if err := row.Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Description,
		&p.CategoryID,
		&p.SupplierID,
		&p.Unit,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// GetBySKU возвращает товар по SKU.
func (r *ProductRepositorySQLite) GetBySKU(sku string) (*models.Product, error) {
	const query = `
SELECT id, sku, name, description, category_id, supplier_id, unit, created_at, updated_at
FROM products
WHERE sku = ? LIMIT 1;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, sku)
	var p models.Product
	if err := row.Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Description,
		&p.CategoryID,
		&p.SupplierID,
		&p.Unit,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Create сохраняет новый товар.
func (r *ProductRepositorySQLite) Create(product *models.Product) error {
	const query = `
INSERT INTO products (id, sku, name, description, category_id, supplier_id, unit, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
`
	if product.ID == "" {
		product.ID = "p-" + time.Now().UTC().Format("20060102T150405.000000000")
	}
	now := time.Now().UTC()
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	product.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		product.ID,
		product.SKU,
		product.Name,
		product.Description,
		product.CategoryID,
		product.SupplierID,
		product.Unit,
		product.CreatedAt,
		product.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// Update обновляет существующий товар.
func (r *ProductRepositorySQLite) Update(product *models.Product) error {
	const query = `
UPDATE products
SET sku = ?, name = ?, description = ?, category_id = ?, supplier_id = ?, unit = ?, updated_at = ?
WHERE id = ?;
`
	now := time.Now().UTC()
	product.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query,
		product.SKU,
		product.Name,
		product.Description,
		product.CategoryID,
		product.SupplierID,
		product.Unit,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product with id %s not found", product.ID)
	}
	return nil
}

// Delete удаляет товар по ID.
func (r *ProductRepositorySQLite) Delete(id string) error {
	const query = `DELETE FROM products WHERE id = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}
	return nil
}
