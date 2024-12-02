package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"l0/cmd/kafka"
	"l0/internal/cache"
	"l0/internal/handler"
	"l0/internal/repository"
	"l0/internal/service"
	"l0/migrations"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Загрузка переменных окружения
	err := godotenv.Load(".env")
	if err != nil {
		logrus.WithError(err).Fatal("Ошибка загрузки .env файла")
	}

	// Получение переменных окружения для базы данных
	dsn := os.Getenv("POSTGRES_CONN")
	if dsn == "" {
		logrus.Fatal("Переменная окружения POSTGRES_CONN не задана")
	}

	// Подключение к базе данных
	db, err := migrations.ConnectDB(dsn)
	if err != nil {
		logrus.WithError(err).Fatal("Ошибка подключения к базе данных")
	}

	// Заполнение бд данными
	orderIDs, err := migrations.SeedDB(db)
	if err != nil {
		logrus.WithError(err).Fatal("Ошибка заполнения базы данных данными")
	}

	if orderIDs != nil {
		logrus.WithField("orderIDs", orderIDs).Info("Созданы заказы")
	}

	//Инициализация кэша,сервиса
	orderCache := cache.NewCache()
	orderService := service.New(repository.NewRepository(db), orderCache)

	// Настройки Kafka
	broker := os.Getenv("KAFKA_BROKER")
	topic := "orders"

	// Инициализация Kafka Producer
	producer := kafka.InitProducer(broker, topic)
	defer producer.Close()

	// Извлечение данных из базы через репозиторий и отправка в Kafka
	orders, err := orderService.GetOrders()
	if err != nil {
		logrus.WithError(err).Fatal("Ошибка извлечения данных из базы")
	}

	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			logrus.WithError(err).WithField("orderUID", order.Order_uid).Warn("Ошибка сериализации заказа")
			continue
		}

		err = producer.SendMessage(order.Order_uid, string(orderJSON))
		if err != nil {
			logrus.WithError(err).WithField("orderUID", order.Order_uid).Error("Ошибка отправки сообщения в Kafka")
		}
	}

	// Запуск Kafka Consumer
	go kafka.ConsumeMessages(broker, topic, orderService)

	// Создаем обработчик HTTP
	orderHandler := handler.NewHandler(orderService)

	// Получаем адрес сервера из переменной окружения
	serverAddress := os.Getenv("SERVER_ADDRESS")

	// Регистрируем маршруты
	http.HandleFunc("/", orderHandler.RenderHTML)
	http.HandleFunc("/orders", orderHandler.GetOrder)

	// Запуск HTTP сервера
	go func() {
		log.Printf("Запуск HTTP сервера на %s\n", serverAddress)
		if err := http.ListenAndServe(serverAddress, nil); err != nil {
			logrus.WithError(err).Fatal("Ошибка при запуске HTTP сервера")
		}
	}()

	// Обработка прерывания работы сервиса
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала на завершение
	<-sigChan

	logrus.Info("Приложение завершило работу")
}
