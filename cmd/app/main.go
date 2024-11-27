package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"l0/cmd/kafka"
	"l0/internal/service"
	"log"
	"os"
	"time"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

	// Инициализация базы данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Инициализация сервиса
	orderService := service.NewOrderService(db)

	// Настройки Kafka
	broker := os.Getenv("KAFKA_BROKER")
	topic := "orders"

	// Инициализация Kafka Producer
	kafka.InitProducer(broker, topic)
	defer kafka.CloseProducer()

	// Пример отправки сообщения
	for i := 0; i < 5; i++ {
		message := "order " + time.Now().Format(time.RFC3339)
		err := kafka.SendMessage("order_key", message)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	log.Println("Сообщение отправлено в Kafka")

	// Запуск Kafka Consumer
	go kafka.ConsumeMessages(broker, topic, orderService)

	log.Println("Application is running")
	select {}
}

//func main() {
//	// Загрузка переменных окружения
//	err := godotenv.Load("local.env")
//	if err != nil {
//		log.Fatalf("Error loading .env file: %v", err)
//	}
//
//	// Настройки Kafka
//	broker := os.Getenv("KAFKA_BROKER")
//	topic := "orders"
//
//	// Инициализация Kafka Producer
//	kafka.InitProducer(broker, topic)
//	defer kafka.CloseProducer()
//
//	//	// Инициализация сервиса
//	orderService := service.NewOrderService(db)
//
//	// Запуск Kafka Consumer
//	go kafka.ConsumeMessages(broker, topic, orderService)
//
//	// Пример отправки сообщения
//	for i := 0; i < 5; i++ {
//		message := "Order " + time.Now().Format(time.RFC3339)
//		err := kafka.SendMessage("order_key", message)
//		if err != nil {
//			log.Printf("Failed to send message: %v", err)
//		}
//		time.Sleep(2 * time.Second)
//	}
//
//	log.Println("Application is running")
//}
