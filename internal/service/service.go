package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"l0/internal/model"
	"log"
)

type OrderCache interface {
	IsEmpty() bool
	CreateOrder(order model.Order)
	GetOrder(id string) (model.Order, error)
}

type Repository interface {
	CreateOrder(order model.Order) error
	GetOrder(id string) (model.Order, error)
	GetOrders() ([]model.Order, error)
}

type Service struct {
	repository Repository
	orderCache OrderCache
}

func (s *Service) GetOrders() ([]model.Order, error) {
	return s.repository.GetOrders()
}

func New(repository Repository, orderCache OrderCache) *Service {
	service := &Service{repository: repository, orderCache: orderCache}
	if service.orderCache.IsEmpty() {
		orders, err := service.repository.GetOrders()
		if err != nil {
			log.Println("Не удалось обновить кэш из базы данных")
		}

		// Заполняем кэш заказами из бд
		for _, order := range orders {
			service.orderCache.CreateOrder(order)
		}
	}
	return service
}

// Получаем заказ, сначала из кэша, затем из бд
func (s *Service) GetOrder(id string) (model.Order, error) {
	var order model.Order
	var err error

	order, err = s.orderCache.GetOrder(id)
	if errors.Is(err, nil) {
		return order, nil
	}

	order, err = s.repository.GetOrder(id)
	if err != nil {
		return order, fmt.Errorf("не удалось найти заказ с id %s: %w", id, err)
	}

	// После получения заказа из бд, добавляем его в кэш
	s.orderCache.CreateOrder(order)
	return order, nil
}

func (s *Service) CreateOrder(order model.Order) error {
	if !model.ValidateOrder(order) {
		return errors.New("невалидный заказ")
	}

	//есть ли заказ в кэше
	_, err := s.orderCache.GetOrder(order.Order_uid)
	if err != nil {
		log.Printf("заказ с ID '%s' уже существует в кэше", order.Order_uid)
		return nil
	}

	//есть ли заказ в бд
	_, err = s.repository.GetOrder(order.Order_uid)
	if err != nil {
		return fmt.Errorf("заказ с ID '%s' уже существует в базе данных", order.Order_uid)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("ошибка при проверке существования заказа: %w", err)
	}

	err = s.repository.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("не удалось создать заказ в репозитории: %w", err)
	}
	s.orderCache.CreateOrder(order)
	return nil
}
