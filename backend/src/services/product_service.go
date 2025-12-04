package services

import (
	"errors"
	"strings"
	"warehouse-management-system/src/models"
)

// ProductRepository описывает поведение хранилища товаров, нужное слою сервисов.
type ProductRepository interface {
	GetAll() ([]*models.Product, error)
	GetByID(id string) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id string) error
}

// ProductService инкапсулирует бизнес-логику работы с товарами.
type ProductService struct {
	repo ProductRepository
}

// NewProductService — конструктор сервиса товаров.
func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrSKUAlreadyUsed  = errors.New("product with this SKU already exists")
	ErrInvalidProduct  = errors.New("invalid product data")
)

// ListProducts возвращает все товары.
// Позже сюда можно будет добавить параметры пагинации, фильтрации и сортировки.
func (s *ProductService) ListProducts() ([]*models.Product, error) {
	return s.repo.GetAll()
}

// GetProduct возвращает товар по ID.
func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

// CreateProduct создаёт новый товар.
func (s *ProductService) CreateProduct(sku, name, description, categoryID, supplierID, unit string) (*models.Product, error) {
	sku = strings.TrimSpace(sku)
	name = strings.TrimSpace(name)
	categoryID = strings.TrimSpace(categoryID)

	if sku == "" || name == "" || categoryID == "" {
		return nil, ErrInvalidProduct
	}

	// Проверяем уникальность SKU на уровне сервиса, чтобы бизнес-правило не зависело от конкретной БД.
	existing, err := s.repo.GetBySKU(sku)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSKUAlreadyUsed
	}

	product := models.NewProduct(
		sku,
		name,
		description,
		categoryID,
		supplierID,
		models.UnitOfMeasure(unit),
	)

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProduct обновляет данные товара.
func (s *ProductService) UpdateProduct(id, sku, name, description, categoryID, supplierID, unit string) (*models.Product, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidProduct
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	sku = strings.TrimSpace(sku)
	name = strings.TrimSpace(name)
	categoryID = strings.TrimSpace(categoryID)

	if sku == "" || name == "" || categoryID == "" {
		return nil, ErrInvalidProduct
	}

	// Проверка уникальности SKU при изменении.
	if product.SKU != sku {
		existing, err := s.repo.GetBySKU(sku)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrSKUAlreadyUsed
		}
	}

	product.SKU = sku
	product.Name = name
	product.Description = description
	product.CategoryID = categoryID
	product.SupplierID = supplierID
	product.Unit = models.UnitOfMeasure(unit)
	// Обновление UpdatedAt можно сделать здесь или на уровне репозитория/БД.

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct удаляет товар по ID.
func (s *ProductService) DeleteProduct(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrInvalidProduct
	}

	return s.repo.Delete(id)
}
