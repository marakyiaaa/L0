package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

//чтение сообщений

func ConsumeMessages(broker string, topic string) {
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
	}
}
