package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"gorm.io/gorm"
	"log"
	"strings"
)

// PostgresProductRepository handles PostgreSQL-specific queries
type PostgresProductRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormProductRepository
}

// GetProducts fetches all products using raw SQL (native query)
func (r *PostgresProductRepository) GetProducts(limit, offset int, name string) (map[string]interface{}, error) {
	query := "SELECT id, name, sku, price FROM products"
	countQuery := "SELECT COUNT(*) FROM products"
	var args []interface{}
	var conditions []string // Store conditions dynamically

	// Add name filtering if provided
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+name+"%")
	}

	// If any conditions exist, add WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Channels for concurrency
	countChan := make(chan int, 1)
	productsChan := make(chan []products.Product, 1)
	errChan := make(chan error, 2)
	defer close(countChan)
	defer close(productsChan)
	defer close(errChan)

	// Goroutine: Fetch total count
	go func() {
		var totalCount int
		err := r.SQLDB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&totalCount) // Remove limit & offset from args
		if err != nil {
			errChan <- err
			return
		}
		countChan <- totalCount
	}()

	// Goroutine: Fetch paginated products
	go func() {
		rows, err := r.SQLDB.QueryContext(context.Background(), query, args...)
		if err != nil {
			errChan <- err
			return
		}
		defer rows.Close()

		var productsList []products.Product
		for rows.Next() {
			var p products.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Sku, &p.Price); err != nil {
				errChan <- err
				return
			}
			productsList = append(productsList, p)
		}

		if err = rows.Err(); err != nil {
			errChan <- err
			return
		}

		productsChan <- productsList
	}()

	// Collect results
	var totalCount int
	var productsList []products.Product
	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			log.Printf("[DB_ERROR] %v", err)
			return nil, err
		case totalCount = <-countChan:
		case productsList = <-productsChan:
		}
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit

	// Return response
	return map[string]interface{}{
		"products":   productsList,
		"totalPages": totalPages,
	}, nil
}
