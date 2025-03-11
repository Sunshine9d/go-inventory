package postgres

import (
	"database/sql"
	"github.com/Sunshine9d/go-inventory/internal/products"
)

// PostgresProductRepository handles PostgreSQL-specific queries
type PostgresProductRepository struct {
	DB *sql.DB
}

func (r *PostgresProductRepository) GetProducts(limit, offset int, name string) ([]products.Product, error) {
	query := `SELECT id, name, quantity, price FROM products WHERE name ILIKE $1 LIMIT $2 OFFSET $3`
	rows, err := r.DB.Query(query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsList []products.Product
	for rows.Next() {
		var p products.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, err
		}
		productsList = append(productsList, p)
	}

	return productsList, nil
}

// Other CRUD methods (CreateProduct, GetProductByID, UpdateProduct, DeleteProduct)
