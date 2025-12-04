package models

import "time"

// Role представляет роль пользователя в системе.
type Role string

const (
	RoleAdmin       Role = "admin"
	RoleManager     Role = "manager"
	RoleStorekeeper Role = "storekeeper"
)

// User — доменная модель пользователя.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // не возвращаем хэш наружу
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser — фабричный метод для создания нового пользователя на доменном уровне.
func NewUser(id, email, passwordHash string, role Role) *User {
	now := time.Now().UTC()
	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
