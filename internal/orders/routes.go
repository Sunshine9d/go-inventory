package orders

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, handler *Handler) {
	router.HandleFunc("/orders", handler.GetOrders).Methods("GET")
}
