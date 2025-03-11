package services

import (
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/products"
)

// Services struct to hold all service dependencies
type Services struct {
	ProductService *products.Service
	OrderService   *orders.Service // Example: Additional service
}

// Handlers struct to hold all HTTP handlers
type Handlers struct {
	ProductHandler *products.Handler
	OrderHandler   *orders.Handler // Example: Additional handler
}

// InitializeServices creates all services and handlers
func InitializeServices(productRepo products.Repository, orderRepo orders.Repository) (*Services, *Handlers) {
	// Initialize services
	productService := &products.Service{Repo: productRepo}
	orderService := &orders.Service{Repo: orderRepo} // Example

	// Initialize handlers
	productHandler := &products.Handler{Service: productService}
	orderHandler := &orders.Handler{Service: orderService} // Example

	// Return both services and handlers
	return &Services{
			ProductService: productService,
			OrderService:   orderService,
		}, &Handlers{
			ProductHandler: productHandler,
			OrderHandler:   orderHandler,
		}
}
