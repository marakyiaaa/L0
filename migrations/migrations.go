package migrations

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"l0/internal/model"
	"log"
	"math/rand"
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

func randomString(lenght int) string {
	const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	x := make([]byte, lenght)
	for i := range x {
		x[i] = char[rand.Intn(len(char))]
	}
	return string(x)
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + max
}

func SeedDB(db *gorm.DB) ([]string, error) {
	var count int64
	if err := db.Model(&model.Order{}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("ошибка проверки данных в базе: %w", err)
	}
	//if count > 0 {
	//	log.Println("Данные уже существуют, пропускаем заполнение")
	//	return nil, nil
	//}

	var orders []model.Order
	var createdOrderIDs []string

	for i := 0; i < 3; i++ {
		orderUID := randomString(10)
		order := model.Order{
			Order_uid:    orderUID,
			Track_number: randomString(15),
			Entry:        "WBIL",
			Delivery: model.Delivery{
				Name:    fmt.Sprintf("User %d", i+1),
				Phone:   fmt.Sprintf("+972%07d", randomInt(1000000, 9999999)),
				Zip:     fmt.Sprintf("%05d", randomInt(10000, 99999)),
				City:    "Random City",
				Address: fmt.Sprintf("Street %d", i+1),
				Region:  "Random Region",
				Email:   fmt.Sprintf("user%d@example.com", i+1),
			},
			Payment: model.Payment{
				Transaction:  randomString(10),
				RequestId:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       randomInt(1000, 5000),
				PaymentDt:    randomInt(1, 100000),
				Bank:         "Random Bank",
				DeliveryCost: randomInt(100, 500),
				GoodsTotal:   randomInt(500, 1500),
				CustomFee:    randomInt(0, 50),
			},
			Items: []model.Items{
				{
					ChrtId:      randomInt(1000000, 9999999),
					TrackNumber: randomString(15),
					Price:       randomInt(100, 500),
					Rid:         randomString(10),
					Name:        fmt.Sprintf("Item %d", i+1),
					Sale:        randomInt(10, 50),
					Size:        "L",
					TotalPrice:  randomInt(500, 1500),
					NmId:        randomInt(100000, 999999),
					Brand:       "Random Brand",
					Status:      randomInt(200, 400),
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerId:        randomString(6),
			DeliveryService:   "meest",
			Shardkey:          fmt.Sprintf("%d", randomInt(1, 10)),
			SmId:              randomInt(10, 100),
			DateCreated:       time.Now(),
			OofShard:          "1",
		}
		orders = append(orders, order)
		createdOrderIDs = append(createdOrderIDs, orderUID)
	}

	// Добавляем в базу данных
	if err := db.Create(&orders).Error; err != nil {
		return nil, fmt.Errorf("Ошибка создания записи: %v", err)
	}

	log.Println("База данных успешно заполнена")
	return createdOrderIDs, nil
}
