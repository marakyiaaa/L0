package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"l0/internal/model"
)

// отправка сообщений
type Producer struct {
	writer *kafka.Writer
}

func InitProducer(broker string, topic string) *Producer {
	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	logrus.Info("producer Kafka запущен")
	return &Producer{
		writer: kafkaWriter,
	}
}

func (p *Producer) Close() {
	if p.writer != nil {
		p.writer.Close()
	}
}
func (p *Producer) SendMessage(key, message string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	}

	err := p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) SendOrderMessage(order model.Order, key string) error {
	// Сериализация структуры Order в JSON
	messageBytes, err := json.Marshal(order)
	if err != nil {
		logrus.Info("Ошибка сериализации JSON: %v", err)
		return err
	}
	// Отправка JSON-сообщения
	return p.SendMessage(key, string(messageBytes))
}
