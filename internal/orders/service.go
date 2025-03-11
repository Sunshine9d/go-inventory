package orders

type Service struct {
	Repo Repository
}

func (s *Service) GetOrders() ([]Order, error) {
	return s.Repo.GetOrders()
}
