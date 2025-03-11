package main

import (
	"github.com/Sunshine9d/go-inventory/internal/db"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: No .env file found. Using system environment variables.")
	}

	// Initialize repositories
	productRepo, err := db.GetProductRepository()
	if err != nil {
		log.Fatal("❌", err)
	}

	orderRepo, err := db.GetOrderRepository()
	if err != nil {
		log.Fatal("❌", err)
	}
	// Run database migrations
	db.MigrateDB()
	// Initialize services
	productService := &products.Service{Repo: productRepo}
	orderService := &orders.Service{Repo: orderRepo}
	// Initialize handlers
	productHandler := &products.Handler{Service: productService}
	orderHandler := &orders.Handler{Service: orderService}
	// Setup routes
	router := mux.NewRouter()
	products.RegisterRoutes(router, productHandler)
	orders.RegisterRoutes(router, orderHandler)

	// Start server
	log.Println("🚀 Server running on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", router))
}
