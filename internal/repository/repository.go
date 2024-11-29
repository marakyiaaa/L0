package repository

import (
	"gorm.io/gorm"
	"l0/internal/model"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetOrders() ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error
	return orders, err
}

func (r *Repository) CreateOrder(order *model.Order) error {
	return r.db.Create(&order).Error
}
