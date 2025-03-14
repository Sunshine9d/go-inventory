package orders

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *Service
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit, offset := getLimitOffset(r)

	var id *int
	if idParam := r.URL.Query().Get("id"); idParam != "" {
		parsedID, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}
		id = &parsedID
	}

	var customerName *string
	if nameParam := r.URL.Query().Get("customer_name"); nameParam != "" {
		customerName = &nameParam
	}

	// Fetch orders
	orders, err := h.Service.GetOrders(limit, offset, id, customerName)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	// Set JSON response header and encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	order, err := h.Service.GetOrderByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func getLimitOffset(r *http.Request) (int, int) {
	// Default values
	limit := 10
	offset := 0
	// Parse limit and offset from query parameters
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	return limit, offset
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	err = h.Service.CreateOrder(&order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	var order Order
	err = json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	order.ID = id
	err = h.Service.UpdateOrder(&order)
	if err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	err = h.Service.DeleteOrder(id)
	if err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}
}
