package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

//отправка сообщений

var kafkaWriter *kafka.Writer

func InitProducer(broker string, topic string) {
	kafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll, // Подтверждение всех реплик
	}
	log.Println("Продюсер Kafka инициализирован")
}

func CloseProducer() {
	if kafkaWriter != nil {
		kafkaWriter.Close()
	}
}

func SendMessage(key, message string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	}

	err := kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.Printf("Сообщение отправлено: %s", message)
	return nil
}
