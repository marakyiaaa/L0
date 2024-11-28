package main

import (
	"encoding/json"
	"flag"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"l0/cmd/kafka"
	"l0/internal/model"
	"l0/internal/service"
	"log"
	"os"
)

func main() {
	// Флаг для выбора режима работы
	writeData := flag.Bool("write-data", false, "Write data to database from JSON file")
	flag.Parse()

	// Загрузка переменных окружения
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Получение переменных окружения для базы данных
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USERNAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DATABASE")

	// Формирование строки подключения
	dsn := "host=" + dbHost +
		" port=" + dbPort +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" sslmode=disable"

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Если выбран режим записи данных
	if *writeData {
		err = model.WriteDataDB(db, "internal/model/model.json")
		if err != nil {
			log.Fatalf("Ошибка записи данных в базу: %v", err)
		}
		log.Println("Данные успешно записаны в базу.")
		return
	}

	// Инициализация сервиса
	orderService := service.NewOrderService(db)

	// Настройки Kafka
	broker := os.Getenv("KAFKA_BROKER")
	topic := "orders"

	// Инициализация Kafka Producer
	kafka.InitProducer(broker, topic)
	defer kafka.CloseProducer()

	// Извлечение данных из базы и отправка в Kafka
	var orders []model.Order
	if err := db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error; err != nil {
		log.Fatalf("Ошибка извлечения данных из базы: %v", err)
	}

	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Ошибка сериализации заказа: %v", err)
			continue
		}
		err = kafka.SendMessage(order.Order_uid, string(orderJSON))
		if err != nil {
			log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		} else {
			log.Printf("Заказ %s успешно отправлен в Kafka", order.Order_uid)
		}
	}

	// Запуск Kafka Consumer
	go kafka.ConsumeMessages(broker, topic, orderService)

	log.Println("Приложение запущено")
	select {}
}
