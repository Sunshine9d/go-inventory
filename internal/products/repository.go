package products

// Repository defines database operations for products
type Repository interface {
	GetProducts(limit, offset int, name string) ([]Product, error)
	//CreateProduct(ctx context.Context, product *Product) error
	//GetProductByID(ctx context.Context, id int) (*Product, error)
	//UpdateProduct(ctx context.Context, product *Product) error
	//DeleteProduct(ctx context.Context, id int) error
}
