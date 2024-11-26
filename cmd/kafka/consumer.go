package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"l0/internal/model"
	"l0/internal/service"
	"log"
	"time"
)

//чтение сообщений

func ConsumeMessages(broker string, topic string, orderService *service.OrderService) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		GroupID:        "group-1",
		CommitInterval: time.Second, // Интервал фиксации смещения
	})
	defer r.Close()

	log.Println("consumer Kafka запущен")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			continue
		}
		log.Printf("Сообщение получено: key=%s, value=%s", string(m.Key), string(m.Value))

		var order model.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Ошибка при декодировании данных: %v", err)
			continue
		}

		//// Создание заказа в БД и кэше через OrderService
		if err := orderService.CreateOrder(&order); err != nil {
			log.Printf("Ошибка при сохранении заказа: %v", err)
			continue
		}

		log.Printf("Заказ успешно обработан: %s", order.Order_uid)
	}
}
