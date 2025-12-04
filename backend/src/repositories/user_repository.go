package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"warehouse-management-system/src/models"
)

// UserRepositorySQLite — реализация хранилища пользователей на SQLite.
// Она использует таблицу users, описанную в миграции MigrateSQLite.
type UserRepositorySQLite struct {
	db *sql.DB
}

// NewUserRepository создаёт новый репозиторий пользователей на основе *sql.DB.
func NewUserRepository(db *sql.DB) *UserRepositorySQLite {
	return &UserRepositorySQLite{db: db}
}

var (
	// ErrUserNotFound может использоваться вызывающим кодом при необходимости.
	ErrUserNotFound = errors.New("user not found")
)

// FindByEmail ищет пользователя по email.
func (r *UserRepositorySQLite) FindByEmail(email string) (*models.User, error) {
	const query = `
SELECT id, email, password_hash, role, created_at, updated_at
FROM users
WHERE email = ? LIMIT 1;
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, email)

	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

// FindByID ищет пользователя по ID.
func (r *UserRepositorySQLite) FindByID(id string) (*models.User, error) {
	const query = `
SELECT id, email, password_hash, role, created_at, updated_at
FROM users
WHERE id = ? LIMIT 1;
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

// Create сохраняет нового пользователя.
func (r *UserRepositorySQLite) Create(user *models.User) error {
	const query = `
INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);
`

	// В SQLite мы генерируем ID на уровне приложения (т.к. в схеме он TEXT).
	if user.ID == "" {
		user.ID = "u-" + time.Now().UTC().Format("20060102T150405.000000000")
	}

	now := time.Now().UTC()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
