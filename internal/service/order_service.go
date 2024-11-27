package service

import (
	"fmt"
	"gorm.io/gorm"
	"l0/cmd/kafka"
	"l0/internal/model"
	"sync"
)

type OrderService struct {
	cache sync.Map // Потокобезопасный кэш в памяти
	//Потокобезопасность позволяет использовать кэш в многопоточных приложениях
	//без дополнительных блокировок (например, sync.Mutex).
	db *gorm.DB // Подключение к базе данных
}

// создает новый экземпляр OrderService
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		cache: sync.Map{}, // Инициализируем пустой кэш
		db:    db,         // Подключение к базе передается как аргумент
	}
}

func (s *OrderService) GetOrderByID(orderUID string) (*model.Order, error) {
	// Попытка найти в кэше
	if val, ok := s.cache.Load(orderUID); ok {
		order, _ := val.(*model.Order)
		return order, nil
	}

	// Если не найдено, запрос к БД
	var order model.Order
	if err := s.db.Preload("Delivery").Preload("Payment").Preload("Items").Where("order_uid = ?", orderUID).First(&order).Error; err != nil {
		return nil, err
	}

	//сохраняем в кэш
	//Метод sync.Map.Store добавляет или обновляет запись в кэше
	s.cache.Store(orderUID, &order)
	return &order, nil
}

// Добавить заказ в и бд, и кэш

func (s *OrderService) CreateOrder(order *model.Order) error {
	if err := s.db.Create(order).Error; err != nil {
		//gorm.DB.Create сохраняет новый заказ в бд
		return err
	}
	s.cache.Store(order.Order_uid, order)

	// Отправляем сообщение в Kafka о создании заказа
	message := fmt.Sprintf("Order %s created", order.Order_uid)
	if err := kafka.SendMessage("order_key", message); err != nil {
		return fmt.Errorf("failed to send Kafka message: %v", err)
	}

	return nil
}
