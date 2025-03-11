package repository

import (
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	DB *gorm.DB
}

// GetProductByID fetches a product using GORM
//func (r *GormProductRepository) GetProductByID(id int) (products.Product, error) {
//	fmt.Println("GORM query")
//	var p products.Product
//	err := r.DB.First(&p, id).Error
//	if err != nil {
//		return p, err
//	}
//	return p, nil
//}
