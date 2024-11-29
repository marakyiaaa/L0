package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

//import (
//	"context"
//	"encoding/json"
//	"github.com/segmentio/kafka-go"
//	"l0/internal/model"
//	"l0/internal/service"
//	"log"
//	"time"
//)
//
////чтение сообщений
//
//func ConsumeMessages(broker string, topic string, orderService *service.OrderService) {
//	r := kafka.NewReader(kafka.ReaderConfig{
//		Brokers:        []string{broker},
//		Topic:          topic,
//		GroupID:        "group-1",
//		CommitInterval: time.Second,
//	})
//	defer r.Close()
//
//	log.Println("consumer Kafka запущен")
//
//	for {
//		m, err := r.ReadMessage(context.Background())
//		if err != nil {
//			log.Printf("Ошибка при чтении сообщения: %v", err)
//			continue
//		}
//		log.Printf("Сообщение получено: key=%s, value=%s", string(m.Key), string(m.Value))
//
//		if !json.Valid(m.Value) {
//			log.Printf("Некорректный формат JSON: %s", string(m.Value))
//			continue
//		}
//
//		var order model.Order
//		if err := json.Unmarshal(m.Value, &order); err != nil {
//			log.Printf("Ошибка при декодировании данных: %v", err)
//			continue
//		}
//
//		//// Создание заказа в БД и кэше через OrderService
//		if err := orderService.CreateOrder(&order); err != nil {
//			log.Printf("Ошибка при сохранении заказа: %v", err)
//			continue
//		}
//
//		log.Printf("Заказ успешно обработан: %s", order.Order_uid)
//	}
//}

type Consumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewConsumer(brocker []string, topic string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brocker, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (c *Consumer) Consume() {
	partition, err := c.consumer.Partitions(c.topic)
	if err != nil {
		log.Fatalf("Ошибка получения партиций: %v", err)
	}
	for _, partition := range partition {
		partitionConsumer, _ := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Не удалось подключиться к партиции %d: %v", partition, err)
			continue
		}
		defer partitionConsumer.Close()

		go func(pc sarama.PartitionConsumer) {
			for message := range partitionConsumer.Messages() {
				log.Println("Сообщение получено: %s", string(message.Value))
			}
		}(partitionConsumer)

	}
}
