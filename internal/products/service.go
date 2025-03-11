package products

type Service struct {
	Repo Repository
}

func (s *Service) GetProducts() ([]Product, error) {
	return s.Repo.GetProducts(0, 0, "")
}
