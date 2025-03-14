package products

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *Service
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()

	// Get "name" filter
	name := queryParams.Get("name")

	// Parse "limit" (default: 10)
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	// Parse "offset" (default: 0)
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	// Fetch products
	products, err := h.Service.GetProducts(limit, offset, name)
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	// Set response headers & encode JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	var id int = 1
	products, err := h.Service.GetProductByID(id)
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}
