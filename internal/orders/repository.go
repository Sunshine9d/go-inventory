package orders

type Repository interface {
	GetOrders(limit int, offset int, id *int, customerName *string) (map[string]interface{}, error)
	GetOrderByID(id int) (Order, error)
	CreateOrder(order *Order) error
	UpdateOrder(order *Order) error
	DeleteOrder(id int) error
}
