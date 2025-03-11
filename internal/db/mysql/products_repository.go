package mysql

import (
	"database/sql"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"gorm.io/gorm"
)

// MySQLProductRepository handles MySQL-specific queries
type MySQLProductRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
}

// GetProducts fetches all products using raw SQL (native query)
func (r *MySQLProductRepository) GetProducts() ([]products.Product, error) {
	query := "SELECT id, name, quantity, price FROM products"
	rows, err := r.SQLDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsList []products.Product
	for rows.Next() {
		var p products.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, err
		}
		productsList = append(productsList, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return productsList, nil
}

// GetProductByID fetches a product using GORM
func (r *MySQLProductRepository) GetProductByID(id int) (*products.Product, error) {
	var p products.Product
	err := r.DB.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}
