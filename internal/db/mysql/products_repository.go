package mysql

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
)

// MySQLProductRepository handles MySQL-specific queries
type MySQLProductRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormProductRepository
}

// GetProducts fetches all products using raw SQL (native query)
func (r *MySQLProductRepository) GetProducts(limit, offset int, name string) ([]products.Product, error) {
	fmt.Println("Native query")
	// Base query
	query := "SELECT id, name, quantity, price FROM products"
	var args []interface{}

	// Add name filtering if provided
	if name != "" {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	// Add pagination
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	logger.LogQuery(query)
	// Execute query
	rows, err := r.SQLDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse results
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
