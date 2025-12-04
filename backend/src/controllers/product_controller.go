package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"warehouse-management-system/src/services"

	"github.com/gorilla/mux"
)

// ProductController обрабатывает HTTP-запросы для управления товарами.
// Он работает только с DTO и слоем сервисов, не зная о деталях хранилища.
type ProductController struct {
	productService *services.ProductService
}

// NewProductController — конструктор контроллера товаров.
func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

// productRequest описывает тело запроса для создания/обновления товара.
type productRequest struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
	SupplierID  string `json:"supplier_id"`
	Unit        string `json:"unit"`
}

// GetProducts — обработчик получения списка товаров.
// В будущем сюда можно добавить поддержку пагинации, фильтрации и сортировки через query-параметры.
func (c *ProductController) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	products, err := c.productService.ListProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(products)
}

// GetProduct — обработчик получения товара по ID.
func (c *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
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

	product, err := c.productService.GetProduct(id)
	if err != nil {
		if err == services.ErrProductNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}

// CreateProduct — обработчик создания нового товара.
func (c *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	product, err := c.productService.CreateProduct(
		req.SKU,
		req.Name,
		req.Description,
		req.CategoryID,
		req.SupplierID,
		req.Unit,
	)
	if err != nil {
		switch err {
		case services.ErrInvalidProduct:
			w.WriteHeader(http.StatusBadRequest)
		case services.ErrSKUAlreadyUsed:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(product)
}

// UpdateProduct — обработчик обновления товара по ID.
func (c *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	product, err := c.productService.UpdateProduct(
		id,
		req.SKU,
		req.Name,
		req.Description,
		req.CategoryID,
		req.SupplierID,
		req.Unit,
	)
	if err != nil {
		switch err {
		case services.ErrInvalidProduct:
			w.WriteHeader(http.StatusBadRequest)
		case services.ErrProductNotFound:
			w.WriteHeader(http.StatusNotFound)
		case services.ErrSKUAlreadyUsed:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}

// DeleteProduct — обработчик удаления товара по ID.
func (c *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	err := c.productService.DeleteProduct(id)
	if err != nil {
		if err == services.ErrInvalidProduct {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
