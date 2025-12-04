package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"warehouse-management-system/src/models"
)

// CategoryRepositorySQLite — реализация хранилища категорий на SQLite.
type CategoryRepositorySQLite struct {
	db *sql.DB
}

// NewCategoryRepository создаёт новый репозиторий категорий.
func NewCategoryRepository(db *sql.DB) *CategoryRepositorySQLite {
	return &CategoryRepositorySQLite{db: db}
}

// GetAll возвращает все категории.
func (r *CategoryRepositorySQLite) GetAll() ([]*models.Category, error) {
	const query = `
SELECT id, name, created_at, updated_at
FROM categories
ORDER BY name;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetByID возвращает категорию по идентификатору.
func (r *CategoryRepositorySQLite) GetByID(id string) (*models.Category, error) {
	const query = `
SELECT id, name, parent_id, created_at, updated_at
FROM categories
WHERE id = ? LIMIT 1;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id)
	var c models.Category
	if err := row.Scan(
		&c.ID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// Create сохраняет новую категорию.
func (r *CategoryRepositorySQLite) Create(category *models.Category) error {
	const query = `
INSERT INTO categories (id, name, created_at, updated_at)
VALUES (?, ?, ?, ?);
`
	if category.ID == "" {
		category.ID = "c-" + time.Now().UTC().Format("20060102T150405.000000000")
	}
	now := time.Now().UTC()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	category.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		category.ID,
		category.Name,
		category.CreatedAt,
		category.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// Update обновляет существующую категорию.
func (r *CategoryRepositorySQLite) Update(category *models.Category) error {
	const query = `
UPDATE categories
SET name = ?, updated_at = ?
WHERE id = ?;
`
	now := time.Now().UTC()
	category.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query,
		category.Name,
		category.UpdatedAt,
		category.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("category with id %s not found", category.ID)
	}
	return nil
}

// Delete удаляет категорию по ID.
func (r *CategoryRepositorySQLite) Delete(id string) error {
	const query = `DELETE FROM categories WHERE id = ?;`

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
		return fmt.Errorf("category with id %s not found", id)
	}
	return nil
}
