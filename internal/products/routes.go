package products

import "github.com/gorilla/mux"

// RegisterRoutes sets up product-related routes
func RegisterRoutes(router *mux.Router, handler *Handler) {
	router.HandleFunc("/products", handler.GetProducts).Methods("GET")
	router.HandleFunc("/product/{id:[0-9]+}", handler.GetProductByID).Methods("GET")
	//router.HandleFunc("/product", handler.CreateProduct).Methods("POST")
	//router.HandleFunc("/product/{id:[0-9]+}", handler.GetProduct).Methods("GET")
	//router.HandleFunc("/product/{id:[0-9]+}", handler.UpdateProduct).Methods("PUT")
	//router.HandleFunc("/product/{id:[0-9]+}", handler.DeleteProduct).Methods("DELETE")
}
