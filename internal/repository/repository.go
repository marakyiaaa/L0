package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"l0/internal/model"
	"log"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	if err := db.AutoMigrate(&model.Order{}, &model.Delivery{}, &model.Payment{}, &model.Items{}); err != nil {
		return nil, fmt.Errorf("ошибка миграции базы данных: %w", err)
	}

	log.Println("Успешное подключение к базе данных")
	return db, nil
}

// Получаем все заказы из бд
func (r *Repository) GetOrders() ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error
	return orders, err
}

// создание заказа
func (r *Repository) CreateOrder(order *model.Order) error {
	return r.db.Create(&order).Error
}
