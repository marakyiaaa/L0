package main

import (
	"l0/cmd/kafka"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Настройки Kafka
	broker := os.Getenv("KAFKA_BROKER")
	topic := "orders"

	// Инициализация Kafka Producer
	kafka.InitProducer(broker, topic)
	defer kafka.CloseProducer()

	// Запуск Kafka Consumer
	go kafka.ConsumeMessages(broker, topic)

	// Пример отправки сообщения
	for i := 0; i < 5; i++ {
		message := "Order " + time.Now().Format(time.RFC3339)
		err := kafka.SendMessage("order_key", message)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	log.Println("Application is running")
}
