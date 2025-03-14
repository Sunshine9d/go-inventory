package mysql

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
	"log"
)

type MySQLOrderRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormOrderRepository
}

func (r *MySQLOrderRepository) GetOrders(limit, offset int, id *int, customerName *string) (map[string]interface{}, error) {
	// Base query for fetching orders
	query := "SELECT id, customer_name, total_price FROM orders WHERE 1=1"
	var args []interface{}
	argIndex := 1 // PostgreSQL placeholders start from $1

	// Filtering conditions
	if id != nil {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, *id)
		argIndex++
	}
	if customerName != nil {
		query += fmt.Sprintf(" AND customer_name ILIKE $%d", argIndex)
		args = append(args, "%"+*customerName+"%")
		argIndex++
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Debug logging
	fmt.Println("[DEBUG] Query:", query)
	fmt.Println("[DEBUG] Args:", args)
	logger.LogQuery(query)

	// Channels for concurrency
	orderChan := make(chan []orders.Order, 1)
	countChan := make(chan int, 1)
	errChan := make(chan error, 2)

	// Fetch order list (Goroutine)
	go func() {
		defer close(orderChan)
		defer close(errChan)

		rows, err := r.SQLDB.Query(query, args...)
		if err != nil {
			errChan <- err
			return
		}
		defer rows.Close()

		var orderList []orders.Order
		for rows.Next() {
			var order orders.Order
			if err := rows.Scan(&order.ID, &order.CustomerName, &order.TotalPrice); err != nil {
				errChan <- err
				return
			}
			orderList = append(orderList, order)
		}

		orderChan <- orderList
	}()

	// Count total records (Goroutine)
	go func() {
		defer close(countChan)

		countQuery := "SELECT COUNT(*) FROM orders WHERE 1=1"
		var countArgs []interface{}
		countIndex := 1

		if id != nil {
			countQuery += fmt.Sprintf(" AND id = $%d", countIndex)
			countArgs = append(countArgs, *id)
			countIndex++
		}
		if customerName != nil {
			countQuery += fmt.Sprintf(" AND customer_name ILIKE $%d", countIndex)
			countArgs = append(countArgs, "%"+*customerName+"%")
			countIndex++
		}

		var total int
		err := r.SQLDB.QueryRow(countQuery, countArgs...).Scan(&total)
		if err != nil {
			errChan <- err
			return
		}
		countChan <- total
	}()

	// Collect results
	var totalCount int
	var ordersResult []orders.Order
	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			log.Printf("[DB_ERROR] %v", err)
			return nil, err
		case totalCount = <-countChan:
		case ordersResult = <-orderChan:
		}
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit

	// Return response
	return map[string]interface{}{
		"orders":     ordersResult,
		"totalPages": totalPages,
	}, nil
}
