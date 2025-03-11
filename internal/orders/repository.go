package orders

type Repository interface {
	GetOrders(limit, offset int) ([]Order, error)
	GetOrderByID(id int) (Order, error)
	CreateOrder(order *Order) error
	UpdateOrder(order *Order) error
	DeleteOrder(id int) error
}
