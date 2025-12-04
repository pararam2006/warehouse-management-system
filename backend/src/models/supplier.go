package models

import "time"

// Supplier представляет доменную модель поставщика.
type Supplier struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewSupplier — фабричный метод создания поставщика.
func NewSupplier(name, address, phone, email string) *Supplier {
	now := time.Now().UTC()
	return &Supplier{
		ID:        "",
		Name:      name,
		Address:   address,
		Phone:     phone,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}


