package repository

import (
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	DB *gorm.DB
}

// GetProductByID fetches a product using GORM
func (r *GormProductRepository) GetProductByID(id int) (products.Product, error) {
	var p products.Product
	// Log the query manually before executing
	query := "SELECT * FROM products WHERE id = ?"
	logger.LogQuery(query, id)
	err := r.DB.First(&p, id).Error
	if err != nil {
		return p, err
	}
	return p, nil
}
