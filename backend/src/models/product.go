package models

import "time"

// UnitOfMeasure описывает единицу измерения товара (штуки, килограммы и т.п.).
// Используем string, чтобы не привязываться к конкретному типу в БД.
type UnitOfMeasure string

const (
	UnitPiece UnitOfMeasure = "pcs"  // штуки
	UnitKg    UnitOfMeasure = "kg"   // килограммы
	UnitLitre UnitOfMeasure = "l"    // литры
	UnitBox   UnitOfMeasure = "box"  // коробки/упаковки
)

// Product представляет доменную модель товара.
// Здесь нет деталей хранения (таблицы, индексы и т.п.), только бизнес-сущность.
type Product struct {
	ID          string        `json:"id"`
	SKU         string        `json:"sku"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	CategoryID  string        `json:"category_id"`
	SupplierID  string        `json:"supplier_id,omitempty"`
	Unit        UnitOfMeasure `json:"unit"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// NewProduct — фабричный метод создания товара на доменном уровне.
// ID оставляем пустым — его устанавливает репозиторий/БД.
func NewProduct(sku, name, description, categoryID, supplierID string, unit UnitOfMeasure) *Product {
	now := time.Now().UTC()
	return &Product{
		ID:          "",
		SKU:         sku,
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
		SupplierID:  supplierID,
		Unit:        unit,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}


