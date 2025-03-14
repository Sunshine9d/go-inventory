package orders

type Service struct {
	Repo Repository
}

func (s *Service) GetOrders(limit int, offset int, id *int, customerName *string) (map[string]interface{}, error) {
	return s.Repo.GetOrders(limit, offset, id, customerName)
}

func (s *Service) GetOrderByID(id int) (Order, error) {
	return s.Repo.GetOrderByID(id)
}

func (s *Service) CreateOrder(order *Order) error {
	return s.Repo.CreateOrder(order)
}

func (s *Service) UpdateOrder(order *Order) error {
	return s.Repo.UpdateOrder(order)
}

func (s *Service) DeleteOrder(id int) error {
	return s.Repo.DeleteOrder(id)
}
