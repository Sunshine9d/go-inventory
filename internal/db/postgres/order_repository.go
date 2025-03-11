package postgres

import (
	"database/sql"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
)

type PostgresOrderRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormOrderRepository
}

func (r *PostgresOrderRepository) GetOrders(limit, offset int) ([]orders.Order, error) {
	query := "SELECT id, customer_name, total_price FROM orders LIMIT $1 OFFSET $2"
	logger.LogQuery(query)
	rows, err := r.SQLDB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderList []orders.Order
	for rows.Next() {
		var order orders.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.TotalPrice); err != nil {
			return nil, err
		}
		orderList = append(orderList, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orderList, nil
}
