package models

import "time"

// OrderStatus описывает возможные статусы заказа.
type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "new"
	OrderStatusReserved  OrderStatus = "reserved"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCanceled  OrderStatus = "canceled"
)

// OrderItem описывает позицию в заказе.
type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"` // цена продажи за единицу
}

// Order представляет доменную модель заказа.
type Order struct {
	ID         string        `json:"id"`
	Customer   string        `json:"customer"`             // для простоты — строка, без отдельной сущности
	Status     OrderStatus   `json:"status"`
	Items      []OrderItem   `json:"items"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	StatusHist []StatusEntry `json:"status_history"`
}

// StatusEntry описывает изменение статуса заказа.
type StatusEntry struct {
	Status    OrderStatus `json:"status"`
	ChangedAt time.Time   `json:"changed_at"`
}

// NewOrder — фабрика для создания нового заказа.
func NewOrder(customer string, items []OrderItem) *Order {
	now := time.Now().UTC()
	return &Order{
		ID:        "",
		Customer:  customer,
		Status:    OrderStatusNew,
		Items:     items,
		CreatedAt: now,
		UpdatedAt: now,
		StatusHist: []StatusEntry{
			{
				Status:    OrderStatusNew,
				ChangedAt: now,
			},
		},
	}
}


