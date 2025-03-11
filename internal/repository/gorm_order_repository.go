package repository

import (
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	DB *gorm.DB
}

func (r *GormOrderRepository) CreateOrder(order *orders.Order) error {
	return r.DB.Create(order).Error
}

func (r *GormOrderRepository) UpdateOrder(order *orders.Order) error {
	return r.DB.Save(order).Error
}

func (r *GormOrderRepository) DeleteOrder(id int) error {
	return r.DB.Delete(&orders.Order{}, id).Error
}

func (r *GormOrderRepository) GetOrderByID(id int) (orders.Order, error) {
	var order orders.Order
	err := r.DB.First(&order, id).Error
	return order, err
}
