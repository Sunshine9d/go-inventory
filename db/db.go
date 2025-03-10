package db

import "github.com/Sunshine9d/go-inventory/db/mysql"

type ProductDB interface {
	GetProductByID(id int64) (*mysql.Product, error)
	CreateProduct(name string, price float64, description string, stock int32) (*mysql.Product, error)
	UpdateProduct(id int64, name string, price float64, description string, stock int32) (*mysql.Product, error)
	DeleteProduct(id int64) error
	GetProducts(limit, offset int, name string) ([]*mysql.Product, error)
}
