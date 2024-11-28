package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
)

// WriteDataDB записывает данные из JSON в базу данных
func WriteDataDB(db *gorm.DB, filePath string) error {
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

	// Создание таблиц, если они не существуют
	err = db.AutoMigrate(&Order{}, &Delivery{}, &Payment{}, &Items{})
	if err != nil {
		return fmt.Errorf("ошибка миграции базы данных: %w", err)
	}

	// Сохранение данных в базе данных
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("ошибка записи данных заказа: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("ошибка записи данных в базу данных: %w", err)
	}

	log.Println("Данные успешно записаны в базу данных.")
	return nil
}
