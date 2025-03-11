package mysql

import (
	"database/sql"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"github.com/Sunshine9d/go-inventory/pkg/logger"
	"gorm.io/gorm"
	"log"

	"github.com/Sunshine9d/go-inventory/internal/orders"
)

type MySQLOrderRepository struct {
	DB    *gorm.DB
	SQLDB *sql.DB
	*repository.GormProductRepository
}

func (repo *MySQLOrderRepository) GetOrders() ([]orders.Order, error) {
	query := "SELECT id, customer_name, total_price FROM orders"
	rows, err := repo.SQLDB.Query(query)
	// Log query to file
	logger.LogQuery(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderList []orders.Order
	for rows.Next() {
		var order orders.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.TotalPrice); err != nil {
			log.Println("Error scanning order:", err)
			continue
		}
		orderList = append(orderList, order)
	}
	return orderList, nil
}
