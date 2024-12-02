package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"l0/cmd/kafka"
	"l0/internal/cache"
	"l0/internal/handler"
	"l0/internal/model"
	"l0/internal/repository"
	"l0/internal/service"
	"l0/migrations"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @Summary Получение информации о заказе по ID
// @Description Получить заказ по ID
// @ID get-order-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "ID заказа"
// @Success 200 {object} model.Order "Информация о заказе"
// @Failure 400 {object} ErrorResponse "Некорректный запрос"
// @Failure 404 {object} ErrorResponse "Заказ не найден"
// @Router /order/{id} [get]

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Получение переменных окружения для базы данных
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USERNAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DATABASE")

	// Формирование строки подключенияп
	dsn := "host=" + dbHost +
		" port=" + dbPort +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" sslmode=disable"

	// Подключение к базе данных
	db, err := migrations.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Заполнение базы данных тестовыми данными
	if err := migrations.SeedDB(db); err != nil {
		log.Fatalf("Ошибка заполнения базы данных данными: %v", err)
	}

	//Инициализация кэша
	orderCache := cache.NewCache()

	// Инициализация сервиса
	orderService := service.New(repository.NewRepository(db), orderCache)

	// Настройки Kafka
	broker := os.Getenv("KAFKA_BROKER")
	topic := "orders"

	// Инициализация Kafka Producer
	producer := kafka.InitProducer(broker, topic)
	defer producer.Close()

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
		err = producer.SendMessage(order.Order_uid, string(orderJSON))
		if err != nil {
			log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		} else {
			log.Printf("Заказ %s успешно отправлен в Kafka", order.Order_uid)
		}
	}

	// Запуск Kafka Consumer
	go kafka.ConsumeMessages(broker, topic, orderService)

	// Создаем обработчик HTTP
	orderHandler := handler.NewHandler(orderService)

	// Получение адреса сервера из переменной окружения или установка значения по умолчанию
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "127.0.0.1:8080" // Значение по умолчанию
	}

	// Регистрируем маршрут для получения заказа
	http.HandleFunc("/orders", orderHandler.GetOrder)

	// Запуск HTTP сервера на порту 8080
	go func() {
		log.Println("Запуск HTTP сервера на порту 8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Ошибка при запуске HTTP сервера: %v", err)
		}
	}()

	// Обработка прерывания работы сервиса
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала на завершение
	<-sigChan

	log.Println("Приложение завершило работу")
}
