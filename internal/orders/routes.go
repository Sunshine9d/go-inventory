package orders

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, handler *Handler) {
	router.HandleFunc("/orders", handler.GetOrders).Methods("GET").Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}")
	router.HandleFunc("/order/{id:[0-9]+}", handler.GetOrderByID).Methods("GET")
	router.HandleFunc("/order", handler.CreateOrder).Methods("POST")
	router.HandleFunc("/order/{id:[0-9]+}", handler.UpdateOrder).Methods("PUT")
	router.HandleFunc("/order/{id:[0-9]+}", handler.DeleteOrder).Methods("DELETE")
}
