package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
	"log"
)

type PostgresOrderRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormOrderRepository
}

func (r *PostgresOrderRepository) GetOrders(limit, offset int, id *int, customerName *string) (map[string]interface{}, error) {
	// Base query for fetching orders
	query := "SELECT id, customer_name, total_price FROM orders WHERE 1=1"
	var args []interface{}

	// Filtering conditions
	if id != nil {
		query += fmt.Sprintf(" AND id = $%d", len(args)+1)
		args = append(args, *id)
	}
	if customerName != nil {
		query += fmt.Sprintf(" AND customer_name ILIKE $%d", len(args)+1)
		args = append(args, "%"+*customerName+"%")
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, interface{}(limit)) // ✅ Ensure separate values
	args = append(args, interface{}(offset))

	// Debug logging
	fmt.Println("[DEBUG] Query:", query)
	fmt.Println("[DEBUG] Args:", args)
	logger.LogQuery(query, args)

	// Execute Query
	rows, err := r.SQLDB.Query(query, args...)
	if err != nil {
		log.Printf("[DB_ERROR] %v", err)
		return nil, err
	}
	defer rows.Close()

	var orderList []orders.Order
	for rows.Next() {
		var order orders.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.TotalPrice); err != nil {
			log.Printf("[DB_ERROR] %v", err)
			return nil, err
		}
		orderList = append(orderList, order)
	}

	// Check for row errors
	if err = rows.Err(); err != nil {
		log.Printf("[DB_ERROR] %v", err)
		return nil, err
	}

	// Fetch total count
	countQuery := "SELECT COUNT(*) FROM orders WHERE 1=1"
	var countArgs []interface{}

	if id != nil {
		countQuery += fmt.Sprintf(" AND id = $%d", len(countArgs)+1)
		countArgs = append(countArgs, *id)
	}
	if customerName != nil {
		countQuery += fmt.Sprintf(" AND customer_name ILIKE $%d", len(countArgs)+1)
		countArgs = append(countArgs, "%"+*customerName+"%")
	}

	var totalCount int
	err = r.SQLDB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		log.Printf("[DB_ERROR] %v", err)
		return nil, err
	}

	// Return response
	return map[string]interface{}{
		"orders":     orderList,
		"totalPages": (totalCount + limit - 1) / limit,
	}, nil
}
