package models

import "time"

// Category представляет доменную модель категории товара.
// Категория нужна для группировки товаров и дальнейшей фильтрации/отчётности.
type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewCategory — фабричный метод создания категории на доменном уровне.
// ID оставляем пустым — его сгенерирует репозиторий/БД.
func NewCategory(name string) *Category {
	now := time.Now().UTC()
	return &Category{
		ID:        "",
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
