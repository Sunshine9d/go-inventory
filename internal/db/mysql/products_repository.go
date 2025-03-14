package mysql

import (
	"context"
	"database/sql"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"gorm.io/gorm"
	"log"
	"strings"
)

// MySQLProductRepository handles MySQL-specific queries
type MySQLProductRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormProductRepository
}

func (r *MySQLProductRepository) GetProducts(limit, offset int, name string) (map[string]interface{}, error) {
	query := "SELECT id, name, sku, price FROM products"
	countQuery := "SELECT COUNT(*) FROM products"
	var args []interface{}
	var conditions []string

	// Add name filtering if provided
	if name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}

	// If any conditions exist, add WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add pagination
	query += " LIMIT ? OFFSET ?"
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
