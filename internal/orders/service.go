package orders

type Service struct {
	Repo Repository
}

func (s *Service) GetOrders(limit, offset int) ([]Order, error) {
	return s.Repo.GetOrders(limit, offset)
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
