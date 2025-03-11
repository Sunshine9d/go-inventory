package repository

import (
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	DB *gorm.DB
}

// CreateProduct inserts a new product using GORM
func (r *GormProductRepository) CreateProduct(p *products.Product) error {
	return r.DB.Create(p).Error
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

func (r *GormProductRepository) UpdateProduct(p *products.Product) error {
	return r.DB.Save(p).Error
}

// DeleteProduct removes a product using GORM
func (r *GormProductRepository) DeleteProduct(id int) error {
	return r.DB.Delete(&products.Product{}, id).Error
}
