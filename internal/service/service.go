package service

import (
	"errors"
	"fmt"
	"l0/internal/model"
	"log"
)

type OrderCache interface {
	IsEmpty() bool
	CreateOrder(order model.Order)
	GetOrder(id string) (model.Order, error)
}

type Repository interface {
	GetOrders() ([]model.Order, error)
}

type Service struct {
	repository Repository
	orderCache OrderCache
}

func New(repository Repository, orderCache OrderCache) *Service {
	service := &Service{repository: repository, orderCache: orderCache}
	if service.orderCache.IsEmpty() {
		orders, err := service.repository.GetOrders()
		if err != nil {
			log.Println("Я не могу обновить кэsh")
		}

		for _, order := range orders {
			service.orderCache.CreateOrder(order)
		}

	}

	return service
}

func (s *Service) GetOrder(id string) (model.Order, error) {
	var order model.Order
	var err error

	order, err = s.orderCache.GetOrder(id)
	if errors.Is(err, nil) {
		return order, nil
	}

	order, err = s.repository.GetOrder(id)
	if err != nil {
		return order, fmt.Errorf(err)
	}

	return order, nil
}

//
//import (
//	"fmt"
//	"gorm.io/gorm"
//	"l0/internal/model"
//	"log"
//	"sync"
//)
//
//type OrderService struct {
//	cache sync.Map // Потокобезопасный кэш в памяти
//	//Потокобезопасность позволяет использовать кэш в многопоточных приложениях
//	//без дополнительных блокировок (например, sync.Mutex).
//	db *gorm.DB // Подключение к базе данных
//}
//
//// создаем новый экземпляр OrderService
//func NewOrderService(db *gorm.DB) *OrderService {
//	service := &OrderService{
//		cache: sync.Map{}, // Инициализация пустого кэша
//		db:    db,         // Подключение к базе передается как аргумент
//	}
//
//	// Восстановление кэша из БД при старте сервиса
//	service.RestoreCacheFromDB()
//	return service
//}
//
//func (s *OrderService) GetOrderByID(orderUID string) (*model.Order, error) {
//	// Попытка найти в кэше
//	if val, ok := s.cache.Load(orderUID); ok {
//		order, _ := val.(*model.Order)
//		return order, nil
//	}
//
//	// Если не найдено, запрос к БД
//	var order model.Order
//	if err := s.db.Preload("Delivery").Preload("Payment").Preload("Items").Where("order_uid = ?", orderUID).First(&order).Error; err != nil {
//		return nil, err
//	}
//
//	//сохраняем в кэш
//	//Метод sync.Map.Store добавляет или обновляет запись в кэше
//	s.cache.Store(orderUID, &order)
//	return &order, nil
//}
//
//// Добавить заказ в и бд, и кэш
//
//func (s *OrderService) CreateOrder(order *model.Order) error {
//	var existingOrder model.Order
//	// Проверяем, существует ли уже заказ с таким order_uid
//	err := s.db.First(&existingOrder, "order_uid = ?", order.Order_uid).Error
//
//	if err == nil {
//		// Заказ с таким order_uid уже существует, обновляем его
//		err := s.db.Model(&existingOrder).Updates(order).Error
//
//		if err != nil {
//			return err
//		}
//	} else if err == gorm.ErrRecordNotFound {
//		// Заказ не найден, создаём новый
//		err = s.db.Create(order).Error
//		if err != nil {
//			return err
//		}
//	} else {
//		return err
//	}
//
//	// Сохраняем заказ в кэш
//	s.cache.Store(order.Order_uid, order)
//	return nil
//}
//
//// Восстановление кэша из БД
//func (s *OrderService) RestoreCacheFromDB() {
//	var orders []model.Order
//	// Загружаем все заказы из БД
//	if err := s.db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error; err != nil {
//		log.Printf("Ошибка при загрузке заказов из базы: %v", err)
//		return
//	}
//
//	// Добавляем заказы в кэш
//	for _, order := range orders {
//		s.cache.Store(order.Order_uid, &order)
//	}
//}
//
//// Сохраняем кэш в базу данных (например, при завершении работы)
//func (s *OrderService) SaveCacheToDB() error {
//	// Обходим весь кэш и сохраняем в базу
//	s.cache.Range(func(key, value interface{}) bool {
//		order, ok := value.(*model.Order)
//		if ok {
//			// Обновляем заказ в базе данных или создаем новый, если его нет
//			if err := s.db.Save(order).Error; err != nil {
//				log.Printf("Ошибка при сохранении заказа в базу данных: %v", err)
//			}
//		}
//		return true
//	})
//	return nil
//}
//
//func (s *OrderService) Migrate() error {
//	// Выполняем миграцию схемы
//	err := s.db.AutoMigrate(&model.Order{}, &model.Delivery{}, &model.Payment{}, &model.Items{})
//	if err != nil {
//		return fmt.Errorf("ошибка миграции базы данных: %w", err)
//	}
//	log.Println("Миграции выполнены успешно")
//	return nil
//}
