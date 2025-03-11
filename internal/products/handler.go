package products

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Service *Service
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.GetProducts()
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}
