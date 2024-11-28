package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type OrderService interface {
	CreateOrder(order *Order) error
}

// WriteDataDB записывает данные из JSON в базу данных
func WriteDataDB(service OrderService, filePath string) error {
	// Чтение JSON-файла
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла JSON: %w", err)
	}

	// Парсинг JSON в структуру Order
	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	// Валидируем заказ перед записью
	if !validateOrder(order) {
		return fmt.Errorf("данные заказа некорректны")
	}

	// Установка связей для вложенных структур
	order.Delivery.OrderUID = order.Order_uid
	order.Payment.OrderUID = order.Order_uid
	for i := range order.Items {
		order.Items[i].OrderUID = order.Order_uid
	}

	//// Создание таблиц, если они не существуют
	//err = db.AutoMigrate(&Order{}, &Delivery{}, &Payment{}, &Items{})
	//if err != nil {
	//	return fmt.Errorf("ошибка миграции базы данных: %w", err)
	//}

	// Сохранение данных в базе данных
	if err := service.CreateOrder(&order); err != nil {
		return fmt.Errorf("ошибка записи данных через сервис: %w", err)
	}

	log.Println("Данные успешно записаны в базу данных.")
	return nil
}
