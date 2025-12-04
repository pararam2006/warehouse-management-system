package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"warehouse-management-system/src/services"
)

// WarehouseController обрабатывает HTTP-запросы складских операций.
type WarehouseController struct {
	warehouseService *services.WarehouseService
}

// NewWarehouseController — конструктор контроллера складских операций.
func NewWarehouseController(warehouseService *services.WarehouseService) *WarehouseController {
	return &WarehouseController{warehouseService: warehouseService}
}

// receiptRequest описывает тело запроса для приёмки товара.
type receiptRequest struct {
	ProductID  string  `json:"product_id"`
	SupplierID string  `json:"supplier_id"`
	Quantity   float64 `json:"quantity"`
	Price      float64 `json:"price"`
	ExpiryDate string  `json:"expiry_date"` // ISO8601, опционально
}

// writeOffRequest описывает тело запроса для списания товара.
type writeOffRequest struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

// reserveRequest описывает тело запроса для резервирования товара под заказ.
type reserveRequest struct {
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  float64 `json:"quantity"`
}

// Receipt — приёмка товара на склад.
func (c *WarehouseController) Receipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req receiptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	var expiry *time.Time
	if req.ExpiryDate != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiryDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid expiry_date format, expected RFC3339"})
			return
		}
		expiry = &t
	}

	if err := c.warehouseService.Receipt(req.ProductID, req.SupplierID, req.Quantity, req.Price, expiry); err != nil {
		if err == services.ErrInvalidOperation {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// WriteOff — списание товара.
func (c *WarehouseController) WriteOff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req writeOffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if err := c.warehouseService.WriteOff(req.ProductID, req.Quantity); err != nil {
		if err == services.ErrInvalidOperation {
			w.WriteHeader(http.StatusBadRequest)
		} else if err == services.ErrInsufficientStock {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Reserve — резервирование товара под заказ.
func (c *WarehouseController) Reserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req reserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if err := c.warehouseService.Reserve(req.ProductID, req.OrderID, req.Quantity); err != nil {
		if err == services.ErrInvalidOperation {
			w.WriteHeader(http.StatusBadRequest)
		} else if err == services.ErrInsufficientStock {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetInventory — получение текущих остатков.
func (c *WarehouseController) GetInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	items, err := c.warehouseService.GetInventory()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}


