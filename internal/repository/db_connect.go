package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"l0/internal/model"
	"log"
	"os"
)

var DB *gorm.DB

// функция подключения к бд
func ConnectDB() error {

	//С помощью библиотеки godotenv,загружаем переменные окружения из файла local.env
	dbConn := os.Getenv("POSTGRESS_CONN")
	if dbConn == "" {
		return fmt.Errorf("Переменная окружения POSTGRES_CONN не установлена")
	}

	var err error
	//Подключение к бд (1-типа используем прогрес, 2 -  базовый конфиг gorm)
	DB, err = gorm.Open(postgres.Open(dbConn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к базе данных: %w", err)
	}

	//миграция
	err = DB.AutoMigrate(&model.Order{}, &model.Delivery{}, &model.Payment{}, &model.Items{})
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных")
	return nil
}
