package orders

type Repository interface {
	GetOrders() ([]Order, error)
}
