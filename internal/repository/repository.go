package repository

import (
	"fmt"
	"gorm.io/gorm"
	"l0/internal/model"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Получаем все заказы из бд
func (r *Repository) GetOrders() ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error
	return orders, err
}

// создание заказа
func (r *Repository) CreateOrder(order model.Order) error {
	return r.db.Create(order).Error
}

// Получение заказа по id с загрузкой всех связанных данных
func (r *Repository) GetOrder(id string) (model.Order, error) {
	var order model.Order
	err := r.db.Preload("Delivery").Preload("Payment").Preload("Items").First(&order, "order_uid = ?", id).Error
	if err != nil {
		return model.Order{}, fmt.Errorf("заказ с id %s не найден: %w", id, err)
	}
	return order, nil
}
