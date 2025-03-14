package products

type Service struct {
	Repo Repository
}

func (s *Service) GetProducts(limit int, offset int, name string) (map[string]interface{}, error) {
	return s.Repo.GetProducts(limit, offset, name)
}

func (s *Service) GetProductByID(id int) (Product, error) {
	return s.Repo.GetProductByID(id)
}
