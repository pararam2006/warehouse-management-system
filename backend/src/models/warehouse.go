package models

import "time"

// StockMovementType описывает тип движения товара на складе.
type StockMovementType string

const (
	MovementReceipt  StockMovementType = "receipt"
	MovementWriteOff StockMovementType = "write_off"
	MovementReserve  StockMovementType = "reserve"
)

// StockItem представляет текущий остаток товара на складе.
type StockItem struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

// StockMovement описывает операцию движения товара (приёмка, списание, резервирование).
type StockMovement struct {
	ID         string            `json:"id"`
	Type       StockMovementType `json:"type"`
	ProductID  string            `json:"product_id"`
	SupplierID string            `json:"supplier_id,omitempty"` // только для приёмки
	OrderID    string            `json:"order_id,omitempty"`    // для резервирования под заказ
	Quantity   float64           `json:"quantity"`
	Price      float64           `json:"price,omitempty"`       // цена закупки (приёмка)
	ExpiryDate *time.Time        `json:"expiry_date,omitempty"` // срок годности, если есть
	CreatedAt  time.Time         `json:"created_at"`
}


