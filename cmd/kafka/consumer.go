package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"l0/internal/model"
	"l0/internal/service"
	"time"
)

// Чтение сообщений из Kafka
func ConsumeMessages(broker string, topic string, orderService *service.Service) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		GroupID:        "group-1",
		CommitInterval: time.Second,
	})
	defer r.Close()

	logrus.Info("consumer Kafka запущен")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			logrus.Info("Ошибка при чтении сообщения: %v", err)
			continue
		}
		logrus.Info("Сообщение получено: %s", string(m.Key))

		if !json.Valid(m.Value) {
			logrus.Info("Некорректный формат JSON: %s", string(m.Value))
			continue
		}

		var order model.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			logrus.Info("Ошибка при декодировании данных: %v", err)
			continue
		}

		// Создание заказа через OrderService
		if err := orderService.CreateOrder(order); err != nil {
			logrus.Info("Ошибка при сохранении заказа: %v", err)
			continue
		}

		logrus.Info("Заказ успешно обработан: %s", order.Order_uid)
	}
}
