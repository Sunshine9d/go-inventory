package main

import (
	"github.com/Sunshine9d/go-inventory/internal/db"
	redisdb "github.com/Sunshine9d/go-inventory/internal/db/redis"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/pkg/config"
	"github.com/gorilla/handlers"
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

	// Load Redis configuration
	redisConfig := config.LoadRedisConfig()

	// Initialize Redis
	redisdb.InitRedis(redisConfig.Addr, redisConfig.Password, redisConfig.DB)

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
	//db.MigrateDB()

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

	// Setup CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                       // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Allowed headers
	)

	// Start server with CORS middleware
	log.Println("🚀 Server running on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", corsHandler(router)))
}
