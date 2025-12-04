package services

import (
	"errors"
	"strings"
	"warehouse-management-system/src/models"
)

// CategoryRepository описывает поведение хранилища категорий для слоя сервисов.
type CategoryRepository interface {
	GetAll() ([]*models.Category, error)
	GetByID(id string) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id string) error
}

// CategoryService инкапсулирует бизнес-логику работы с категориями.
type CategoryService struct {
	repo CategoryRepository
}

// NewCategoryService — конструктор сервиса категорий.
func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidCategory  = errors.New("invalid category data")
)

// ListCategories возвращает список всех категорий.
func (s *CategoryService) ListCategories() ([]*models.Category, error) {
	return s.repo.GetAll()
}

// GetCategory возвращает категорию по ID.
func (s *CategoryService) GetCategory(id string) (*models.Category, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidCategory
	}

	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCategoryNotFound
	}
	return c, nil
}

// CreateCategory создаёт новую категорию.
func (s *CategoryService) CreateCategory(name string) (*models.Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrInvalidCategory
	}

	category := models.NewCategory(name)

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory обновляет существующую категорию.
func (s *CategoryService) UpdateCategory(id, name string) (*models.Category, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidCategory
	}

	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidCategory
	}

	category.Name = name

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory удаляет категорию по ID.
func (s *CategoryService) DeleteCategory(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrInvalidCategory
	}

	return s.repo.Delete(id)
}
