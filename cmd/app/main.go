package app

import (
	"l0/internal/repository"
	"log"
	"os"
)

func main() {
	// Чтение переменных окружения из local.env
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Подключение к базе данных
	err = repository.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Запуск сервера (заглушка)
	log.Println("Server is running on", os.Getenv("SERVER_ADDRESS"))
}
