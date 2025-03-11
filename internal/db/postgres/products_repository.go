package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
)

// PostgresProductRepository handles PostgreSQL-specific queries
type PostgresProductRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormProductRepository
}

// GetProducts fetches all products using raw SQL (native query)
func (r *PostgresProductRepository) GetProducts(limit, offset int, name string) ([]products.Product, error) {
	// Base query
	query := "SELECT id, name, quantity, price FROM products"
	var args []interface{}
	argCount := 1 // PostgreSQL uses $1, $2, ... for placeholders

	// Add name filtering if provided
	if name != "" {
		query += fmt.Sprintf(" WHERE name ILIKE $%d", argCount) // ILIKE for case-insensitive search
		args = append(args, "%"+name+"%")
		argCount++
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
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
