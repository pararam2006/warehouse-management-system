package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"warehouse-management-system/src/services"

	"github.com/gorilla/mux"
)

// CategoryController обрабатывает HTTP-запросы, связанные с категориями товаров.
type CategoryController struct {
	categoryService *services.CategoryService
}

// NewCategoryController — конструктор контроллера категорий.
func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

// categoryRequest описывает тело запроса для создания/обновления категории.
type categoryRequest struct {
	Name string `json:"name"`
}

// GetCategories — получение списка категорий.
func (c *CategoryController) GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	categories, err := c.categoryService.ListCategories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(categories)
}

// GetCategory — получение категории по ID.
func (c *CategoryController) GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "id is required"})
		return
	}

	category, err := c.categoryService.GetCategory(id)
	if err != nil {
		if err == services.ErrCategoryNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err == services.ErrInvalidCategory {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(category)
}

// CreateCategory — создание новой категории.
func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	category, err := c.categoryService.CreateCategory(req.Name)
	if err != nil {
		if err == services.ErrInvalidCategory {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(category)
}

// UpdateCategory — обновление категории по ID.
func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "id is required"})
		return
	}

	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	category, err := c.categoryService.UpdateCategory(id, req.Name)
	if err != nil {
		if err == services.ErrInvalidCategory {
			w.WriteHeader(http.StatusBadRequest)
		} else if err == services.ErrCategoryNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(category)
}

// DeleteCategory — удаление категории по ID.
func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "id is required"})
		return
	}

	if err := c.categoryService.DeleteCategory(id); err != nil {
		if err == services.ErrInvalidCategory {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
