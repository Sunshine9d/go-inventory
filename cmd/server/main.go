package main

import (
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/db/mysql"
	"github.com/Sunshine9d/go-inventory/internal/db/postgres"
	"log"
	"net/http"
	"os"

	"github.com/Sunshine9d/go-inventory/internal/products"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	productRepo, err := getProductRepository()
	if err != nil {
		log.Fatal("❌", err)
	}

	// Initialize services and handlers
	productService := &products.Service{Repo: productRepo}
	productHandler := &products.Handler{Service: productService}

	// Setup routes
	router := mux.NewRouter()
	products.RegisterRoutes(router, productHandler)

	// Start server
	log.Println("🚀 Server running on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", router))
}

func getProductRepository() (products.Repository, error) {
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "mysql":
		dbConn, err := mysql.NewConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
		}
		return &mysql.MySQLProductRepository{DB: dbConn}, nil // ✅ No pointer to interface

	case "postgres":
		dbConn, err := postgres.NewConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
		return &postgres.PostgresProductRepository{DB: dbConn}, nil // ✅ No pointer to interface

	default:
		return nil, fmt.Errorf("unsupported database type. Set DB_TYPE to 'mysql' or 'postgres'")
	}
}
