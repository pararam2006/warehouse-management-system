package services

import (
	"errors"
	"strings"
	"warehouse-management-system/src/models"
)

// SupplierRepository описывает поведение хранилища поставщиков для слоя сервисов.
type SupplierRepository interface {
	GetAll() ([]*models.Supplier, error)
	GetByID(id string) (*models.Supplier, error)
	Create(supplier *models.Supplier) error
	Update(supplier *models.Supplier) error
	Delete(id string) error
}

// SupplierService инкапсулирует бизнес-логику работы с поставщиками.
type SupplierService struct {
	repo SupplierRepository
}

// NewSupplierService — конструктор сервиса поставщиков.
func NewSupplierService(repo SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

var (
	ErrSupplierNotFound = errors.New("supplier not found")
	ErrInvalidSupplier  = errors.New("invalid supplier data")
)

// ListSuppliers возвращает список всех поставщиков.
func (s *SupplierService) ListSuppliers() ([]*models.Supplier, error) {
	return s.repo.GetAll()
}

// GetSupplier возвращает поставщика по ID.
func (s *SupplierService) GetSupplier(id string) (*models.Supplier, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidSupplier
	}

	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, ErrSupplierNotFound
	}
	return supplier, nil
}

// CreateSupplier создаёт нового поставщика.
func (s *SupplierService) CreateSupplier(name, address, phone, email string) (*models.Supplier, error) {
	name = strings.TrimSpace(name)
	address = strings.TrimSpace(address)
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)

	if name == "" {
		return nil, ErrInvalidSupplier
	}

	supplier := models.NewSupplier(name, address, phone, email)

	if err := s.repo.Create(supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// UpdateSupplier обновляет данные поставщика.
func (s *SupplierService) UpdateSupplier(id, name, address, phone, email string) (*models.Supplier, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidSupplier
	}

	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, ErrSupplierNotFound
	}

	name = strings.TrimSpace(name)
	address = strings.TrimSpace(address)
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)

	if name == "" {
		return nil, ErrInvalidSupplier
	}

	supplier.Name = name
	supplier.Address = address
	supplier.Phone = phone
	supplier.Email = email

	if err := s.repo.Update(supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// DeleteSupplier удаляет поставщика по ID.
func (s *SupplierService) DeleteSupplier(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrInvalidSupplier
	}

	return s.repo.Delete(id)
}
