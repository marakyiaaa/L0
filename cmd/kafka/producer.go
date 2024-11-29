package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"l0/internal/model"
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
	log.Println("producer Kafka запущен")
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

	log.Printf("Message sent: %s", message)
	return nil
}

func SendOrderMessage(order model.Order, key string) error {
	// Сериализация структуры Order в JSON
	messageBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		return err
	}

	// Отправка JSON-сообщения
	return SendMessage(key, string(messageBytes))
}

// CreateTopicIfNotExist создает топик, если он еще не существует
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
