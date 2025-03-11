package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/Sunshine9d/go-inventory/internal/db"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/services"
)

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	// Initialize repositories
	productRepo, err := db.GetProductRepository()
	if err != nil {
		log.Fatal("❌", err)
	}

	orderRepo, err := db.GetOrderRepository() // Example for another service
	if err != nil {
		log.Fatal("❌", err)
	}

	// Initialize services and handlers
	_, handlers := services.InitializeServices(productRepo, orderRepo)

	// Setup routes
	router := mux.NewRouter()
	products.RegisterRoutes(router, handlers.ProductHandler)
	orders.RegisterRoutes(router, handlers.OrderHandler) // Example

	// Start server
	log.Println("🚀 Server running on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", router))
}
