package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"l0/internal/model"
	"log"
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
	log.Println("producer Kafka запущен")
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

	//log.Printf("Сообщение отправлено: %s", message)
	return nil
}

func (p *Producer) SendOrderMessage(order model.Order, key string) error {
	// Сериализация структуры Order в JSON
	messageBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		return err
	}

	// Отправка JSON-сообщения
	return p.SendMessage(key, string(messageBytes))
}

// создает топик, если он еще не существует
func CreateTopicIfNotExist(broker string, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		log.Printf("Ошибка подключения к Kafka: %v", err)
		return err
	}
	defer conn.Close()

	// Получаем метаданные топиков
	partitions, err := conn.ReadPartitions(topic)
	if err == nil && len(partitions) > 0 {
		log.Printf("Топик '%s' уже существует.", topic)
		return nil
	}

	// Если топик не существует, создаем его
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1, // Количество партиций
		ReplicationFactor: 1, // Количество реплик
	})
	if err != nil {
		log.Printf("Ошибка при создании топика: %v", err)
		return err
	}

	log.Printf("Топик '%s' успешно создан.", topic)
	return nil
}
