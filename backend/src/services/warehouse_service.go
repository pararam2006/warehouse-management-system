package services

import (
	"errors"
	"time"
	"warehouse-management-system/src/models"
)

// WarehouseRepository описывает поведение складского хранилища для слоя сервисов.
type WarehouseRepository interface {
	GetInventory() ([]*models.StockItem, error)
	GetStockByProduct(productID string) (float64, error)
	AddMovement(m *models.StockMovement) error
}

// WarehouseService инкапсулирует бизнес-логику складских операций.
type WarehouseService struct {
	warehouseRepo WarehouseRepository
	productRepo   ProductRepository
}

// NewWarehouseService — конструктор сервиса складских операций.
func NewWarehouseService(warehouseRepo WarehouseRepository, productRepo ProductRepository) *WarehouseService {
	return &WarehouseService{
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
	}
}

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidOperation  = errors.New("invalid warehouse operation data")
)

// Receipt регистрирует приёмку товара на склад.
func (s *WarehouseService) Receipt(productID, supplierID string, quantity, price float64, expiry *time.Time) error {
	if productID == "" || quantity <= 0 {
		return ErrInvalidOperation
	}

	// Убедимся, что товар существует.
	if _, err := s.productRepo.GetByID(productID); err != nil {
		return err
	}

	m := &models.StockMovement{
		ID:         "",
		Type:       models.MovementReceipt,
		ProductID:  productID,
		SupplierID: supplierID,
		Quantity:   quantity,
		Price:      price,
		ExpiryDate: expiry,
		CreatedAt:  time.Now().UTC(),
	}

	return s.warehouseRepo.AddMovement(m)
}

// WriteOff регистрирует списание товара со склада (метод FIFO/LIFO пока не учитывается, только проверка количества).
func (s *WarehouseService) WriteOff(productID string, quantity float64) error {
	if productID == "" || quantity <= 0 {
		return ErrInvalidOperation
	}

	current, err := s.warehouseRepo.GetStockByProduct(productID)
	if err != nil {
		return err
	}
	if current < quantity {
		return ErrInsufficientStock
	}

	m := &models.StockMovement{
		ID:        "",
		Type:      models.MovementWriteOff,
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: time.Now().UTC(),
	}

	return s.warehouseRepo.AddMovement(m)
}

// Reserve резервирует товар под заказ (без создания самого заказа).
func (s *WarehouseService) Reserve(productID, orderID string, quantity float64) error {
	if productID == "" || orderID == "" || quantity <= 0 {
		return ErrInvalidOperation
	}

	current, err := s.warehouseRepo.GetStockByProduct(productID)
	if err != nil {
		return err
	}
	if current < quantity {
		return ErrInsufficientStock
	}

	m := &models.StockMovement{
		ID:        "",
		Type:      models.MovementReserve,
		ProductID: productID,
		OrderID:   orderID,
		Quantity:  quantity,
		CreatedAt: time.Now().UTC(),
	}

	return s.warehouseRepo.AddMovement(m)
}

// GetInventory возвращает текущие остатки по складу.
func (s *WarehouseService) GetInventory() ([]*models.StockItem, error) {
	return s.warehouseRepo.GetInventory()
}
