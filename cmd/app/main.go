package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"l0/cmd/kafka"
	"l0/internal/handler"
	"l0/internal/model"
	"l0/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Флаг для выбора режима работы
	writeData := flag.Bool("write-data", false, "Write data to database from JSON file")

	// Дополнительный флаг для указания пути к файлу (необязательно)
	filePath := flag.String("file", "", "Path to the file to write")

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

	// Формирование строки подключенияп
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

	// Инициализация сервиса
	orderService := service.NewOrderService(db)

	// Выполнение миграций через сервис
	err = orderService.Migrate()
	if err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	log.Println("Миграции выполнены успешно")

	// Если выбран режим записи данных
	if *writeData {
		if *filePath == "" {
			fmt.Print("Введите путь к файлу для записи данных: ")
			_, err := fmt.Scanln(filePath)
			if err != nil {
				log.Fatalf("Ошибка ввода пути к файлу: %v", err)
			}
		}
		// Вызываем метод для записи данных в базу данных
		err = model.WriteDataDB(orderService, *filePath)
		if err != nil {
			log.Fatalf("Ошибка записи данных в базу: %v", err)
		}
		log.Println("Данные успешно записаны в базу.")
		return
	}

	// Настройки Kafka
	broker := os.Getenv("KAFKA_BROKER")
	topic := "orders"

	// Создание топика, если он не существует
	err = kafka.CreateTopicIfNotExist(broker, topic)
	if err != nil {
		log.Fatalf("Ошибка при создании топика: %v", err)
	}

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

	// Создаем обработчик HTTP
	orderHandler := handler.NewOrderHandler(orderService)

	// Получение адреса сервера из переменной окружения или установка значения по умолчанию
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "127.0.0.1:8080" // Значение по умолчанию
	}

	// Регистрируем маршрут для получения заказа
	http.HandleFunc("/order", orderHandler.GetOrder)

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

	// Сохранение кэша в базу данных при завершении работы
	err = orderService.SaveCacheToDB()
	if err != nil {
		log.Printf("Ошибка при сохранении кэша в базу данных: %v", err)
	} else {
		log.Println("Кэш успешно сохранен в базу данных.")
	}

	log.Println("Приложение завершило работу")
	select {}
}
