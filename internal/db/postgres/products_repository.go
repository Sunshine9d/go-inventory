package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"gorm.io/gorm"
	"log"
)

// PostgresProductRepository handles PostgreSQL-specific queries
type PostgresProductRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormProductRepository
}

// GetProducts fetches all products using raw SQL (native query)
func (r *PostgresProductRepository) GetProducts(limit, offset int, name string) ([]products.Product, error) {
	query := "SELECT id, name, sku, price FROM products"
	var args []interface{}
	argCount := 1

	// Add name filtering if provided
	if name != "" {
		query += fmt.Sprintf(" WHERE name ILIKE $%d", argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}

	// Always add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset) // ✅ Now actually appending values

	log.Printf("[DB_LOG] SQL: %s | ARGS: %+v\n", query, args) // Debugging

	// Execute query
	rows, err := r.SQLDB.Query(query, args...)
	if err != nil {
		log.Println("[DB_ERROR]", err)
		return nil, err
	}
	defer rows.Close()

	// Parse results
	var productsList []products.Product
	for rows.Next() {
		var p products.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Sku, &p.Price); err != nil {
			return nil, err
		}
		productsList = append(productsList, p)
		fmt.Println(productsList)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return productsList, nil
}
