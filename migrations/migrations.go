package migrations

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"l0/internal/model"
	"log"
	"time"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	if err := db.AutoMigrate(&model.Order{}, &model.Delivery{}, &model.Payment{}, &model.Items{}); err != nil {
		return nil, fmt.Errorf("ошибка миграции базы данных: %w", err)
	}

	log.Println("Успешное подключение к базе данных")
	return db, nil
}

func SeedDB(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Order{}).Count(&count).Error; err != nil {
		return fmt.Errorf("ошибка проверки данных в базе: %w", err)
	}
	if count > 0 {
		log.Println("Данные уже существуют, пропускаем заполнение")
		return nil
	}

	orders := []model.Order{
		{
			Order_uid:    "b563feb7b2b84b6test",
			Track_number: "WBILMTESTTRACK",
			Entry:        "WBIL",
			Delivery: model.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: model.Payment{
				Transaction:  "b563feb7b2b84b6test",
				RequestId:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []model.Items{
				{
					ChrtId:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmId:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerId:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmId:              99,
			DateCreated:       time.Now(),
			OofShard:          "1",
		},
	}

	if err := db.Create(&orders).Error; err != nil {
		return fmt.Errorf("Ошибка создания записи: %v", err)
	}
	log.Println("База данных успешно заполнена данными")
	return nil
}
