package mysql

import (
	"database/sql"
	"log"

	"github.com/Sunshine9d/go-inventory/internal/orders"
)

type MySQLOrderRepository struct {
	DB *sql.DB
}

func (repo *MySQLOrderRepository) GetOrders() ([]orders.Order, error) {
	query := "SELECT id, customer_name, total_price FROM orders"
	rows, err := repo.DB.Query(query)
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
