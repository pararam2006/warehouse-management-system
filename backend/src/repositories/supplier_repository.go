package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"warehouse-management-system/src/models"
)

// SupplierRepositorySQLite — реализация хранилища поставщиков на SQLite.
type SupplierRepositorySQLite struct {
	db *sql.DB
}

// NewSupplierRepository создаёт новый репозиторий поставщиков.
func NewSupplierRepository(db *sql.DB) *SupplierRepositorySQLite {
	return &SupplierRepositorySQLite{db: db}
}

// GetAll возвращает всех поставщиков.
func (r *SupplierRepositorySQLite) GetAll() ([]*models.Supplier, error) {
	const query = `
SELECT id, name, address, phone, email, created_at, updated_at
FROM suppliers
ORDER BY name;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Supplier
	for rows.Next() {
		var s models.Supplier
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Address,
			&s.Phone,
			&s.Email,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetByID возвращает поставщика по идентификатору.
func (r *SupplierRepositorySQLite) GetByID(id string) (*models.Supplier, error) {
	const query = `
SELECT id, name, address, phone, email, created_at, updated_at
FROM suppliers
WHERE id = ? LIMIT 1;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id)
	var s models.Supplier
	if err := row.Scan(
		&s.ID,
		&s.Name,
		&s.Address,
		&s.Phone,
		&s.Email,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// Create сохраняет нового поставщика.
func (r *SupplierRepositorySQLite) Create(supplier *models.Supplier) error {
	const query = `
INSERT INTO suppliers (id, name, address, phone, email, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);
`
	if supplier.ID == "" {
		supplier.ID = "s-" + time.Now().UTC().Format("20060102T150405.000000000")
	}
	now := time.Now().UTC()
	if supplier.CreatedAt.IsZero() {
		supplier.CreatedAt = now
	}
	supplier.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		supplier.ID,
		supplier.Name,
		supplier.Address,
		supplier.Phone,
		supplier.Email,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// Update обновляет данные поставщика.
func (r *SupplierRepositorySQLite) Update(supplier *models.Supplier) error {
	const query = `
UPDATE suppliers
SET name = ?, address = ?, phone = ?, email = ?, updated_at = ?
WHERE id = ?;
`
	now := time.Now().UTC()
	supplier.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query,
		supplier.Name,
		supplier.Address,
		supplier.Phone,
		supplier.Email,
		supplier.UpdatedAt,
		supplier.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("supplier with id %s not found", supplier.ID)
	}
	return nil
}

// Delete удаляет поставщика по ID.
func (r *SupplierRepositorySQLite) Delete(id string) error {
	const query = `DELETE FROM suppliers WHERE id = ?;`

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
		return fmt.Errorf("supplier with id %s not found", id)
	}
	return nil
}
