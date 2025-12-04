package services

import (
	"errors"
	"strings"
	"time"
	"warehouse-management-system/src/models"
)

// OrderRepository описывает поведение хранилища заказов.
type OrderRepository interface {
	GetAll() ([]*models.Order, error)
	GetByID(id string) (*models.Order, error)
	Create(order *models.Order) error
	Update(order *models.Order) error
}

// OrderService инкапсулирует бизнес-логику работы с заказами,
// включая автоматическое резервирование товаров при создании/изменении заказа.
type OrderService struct {
	orderRepo     OrderRepository
	warehouseRepo WarehouseRepository
	productRepo   ProductRepository
}

// NewOrderService — конструктор сервиса заказов.
func NewOrderService(orderRepo OrderRepository, warehouseRepo WarehouseRepository, productRepo ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
	}
}

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrInvalidOrder   = errors.New("invalid order data")
	ErrOrderBadStatus = errors.New("invalid order status transition")
)

// ListOrders возвращает список заказов.
func (s *OrderService) ListOrders() ([]*models.Order, error) {
	return s.orderRepo.GetAll()
}

// GetOrder возвращает заказ по ID.
func (s *OrderService) GetOrder(id string) (*models.Order, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidOrder
	}

	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

// CreateOrder создаёт новый заказ и автоматически резервирует товары.
func (s *OrderService) CreateOrder(customer string, items []models.OrderItem) (*models.Order, error) {
	customer = strings.TrimSpace(customer)
	if customer == "" || len(items) == 0 {
		return nil, ErrInvalidOrder
	}

	// Проверяем, что товары существуют.
	for _, it := range items {
		if it.ProductID == "" || it.Quantity <= 0 {
			return nil, ErrInvalidOrder
		}
		if _, err := s.productRepo.GetByID(it.ProductID); err != nil {
			return nil, err
		}
	}

	order := models.NewOrder(customer, items)

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Автоматическое резервирование товаров под заказ.
	for _, it := range order.Items {
		if err := s.warehouseRepo.AddMovement(&models.StockMovement{
			ID:        "",
			Type:      models.MovementReserve,
			ProductID: it.ProductID,
			OrderID:   order.ID,
			Quantity:  it.Quantity,
			CreatedAt: time.Now().UTC(),
		}); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// UpdateOrderStatus обновляет статус заказа и фиксирует историю изменений.
func (s *OrderService) UpdateOrderStatus(id string, newStatus models.OrderStatus) (*models.Order, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidOrder
	}

	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// Простейшая проверка допустимости переходов статусов.
	switch order.Status {
	case models.OrderStatusNew, models.OrderStatusReserved:
		// эти статусы можно менять
	default:
		return nil, ErrOrderBadStatus
	}

	now := time.Now().UTC()
	order.Status = newStatus
	order.UpdatedAt = now
	order.StatusHist = append(order.StatusHist, models.StatusEntry{
		Status:    newStatus,
		ChangedAt: now,
	})

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}
