package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"warehouse-management-system/src/services"

	"github.com/gorilla/mux"
)

// SupplierController обрабатывает HTTP-запросы, связанные с поставщиками.
type SupplierController struct {
	supplierService *services.SupplierService
}

// NewSupplierController — конструктор контроллера поставщиков.
func NewSupplierController(supplierService *services.SupplierService) *SupplierController {
	return &SupplierController{supplierService: supplierService}
}

// supplierRequest описывает тело запроса для создания/обновления поставщика.
type supplierRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

// GetSuppliers — получение списка поставщиков.
func (c *SupplierController) GetSuppliers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	suppliers, err := c.supplierService.ListSuppliers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(suppliers)
}

// GetSupplier — получение поставщика по ID.
func (c *SupplierController) GetSupplier(w http.ResponseWriter, r *http.Request) {
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

	supplier, err := c.supplierService.GetSupplier(id)
	if err != nil {
		if err == services.ErrSupplierNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err == services.ErrInvalidSupplier {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(supplier)
}

// CreateSupplier — создание нового поставщика.
func (c *SupplierController) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req supplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	supplier, err := c.supplierService.CreateSupplier(
		req.Name,
		req.Address,
		req.Phone,
		req.Email,
	)
	if err != nil {
		if err == services.ErrInvalidSupplier {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(supplier)
}

// UpdateSupplier — обновление данных поставщика.
func (c *SupplierController) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
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

	var req supplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	supplier, err := c.supplierService.UpdateSupplier(
		id,
		req.Name,
		req.Address,
		req.Phone,
		req.Email,
	)
	if err != nil {
		if err == services.ErrInvalidSupplier {
			w.WriteHeader(http.StatusBadRequest)
		} else if err == services.ErrSupplierNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(supplier)
}

// DeleteSupplier — удаление поставщика по ID.
func (c *SupplierController) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
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

	if err := c.supplierService.DeleteSupplier(id); err != nil {
		if err == services.ErrInvalidSupplier {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


