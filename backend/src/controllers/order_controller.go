package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"warehouse-management-system/src/models"
	"warehouse-management-system/src/services"

	"github.com/gorilla/mux"
)

// OrderController обрабатывает HTTP-запросы, связанные с заказами.
type OrderController struct {
	orderService *services.OrderService
}

// NewOrderController — конструктор контроллера заказов.
func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

// orderItemRequest описывает одну позицию в заказе.
type orderItemRequest struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
}

// createOrderRequest — тело запроса на создание заказа.
type createOrderRequest struct {
	Customer string             `json:"customer"`
	Items    []orderItemRequest `json:"items"`
}

// updateStatusRequest — тело запроса на обновление статуса.
type updateStatusRequest struct {
	Status models.OrderStatus `json:"status"`
}

// GetOrders — получение списка заказов.
func (c *OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	orders, err := c.orderService.ListOrders()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
}

// GetOrder — получение заказа по ID.
func (c *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
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

	order, err := c.orderService.GetOrder(id)
	if err != nil {
		if err == services.ErrOrderNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err == services.ErrInvalidOrder {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(order)
}

// CreateOrder — создание нового заказа с автоматическим резервированием товаров.
func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	items := make([]models.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, models.OrderItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			Price:     it.Price,
		})
	}

	order, err := c.orderService.CreateOrder(req.Customer, items)
	if err != nil {
		if err == services.ErrInvalidOrder {
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
	_ = json.NewEncoder(w).Encode(order)
}

// UpdateOrderStatus — обновление статуса заказа.
func (c *OrderController) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
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

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	order, err := c.orderService.UpdateOrderStatus(id, req.Status)
	if err != nil {
		if err == services.ErrInvalidOrder {
			w.WriteHeader(http.StatusBadRequest)
		} else if err == services.ErrOrderNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err == services.ErrOrderBadStatus {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(order)
}


