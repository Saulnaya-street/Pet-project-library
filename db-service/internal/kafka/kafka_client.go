package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type IKafkaClient interface {
	Close() error
	Publish(ctx context.Context, key string, value []byte) error
	Subscribe(ctx context.Context, handler func(key string, value []byte) error) error
}

type KafkaClient struct {
	writer     *kafka.Writer
	reader     *kafka.Reader
	config     KafkaConfig
	isProducer bool
	isConsumer bool
}

func NewKafkaClient(cfg KafkaConfig, isProducer, isConsumer bool) (IKafkaClient, error) {
	client := &KafkaClient{
		config:     cfg,
		isProducer: isProducer,
		isConsumer: isConsumer,
	}

	if isProducer {
		client.writer = &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		}
		log.Printf("Настроен продюсер Kafka для топика %s", cfg.Topic)
	}

	if isConsumer {
		client.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:     cfg.Brokers,
			Topic:       cfg.Topic,
			GroupID:     cfg.GroupID,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			MaxWait:     1 * time.Second,
			StartOffset: kafka.FirstOffset,
		})
		log.Printf("Настроен потребитель Kafka для топика %s (группа: %s)", cfg.Topic, cfg.GroupID)
	}

	if isProducer {
		log.Printf("Проверка подключения к Kafka брокеру %s...", cfg.Brokers[0])
		conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Brokers[0], cfg.Topic, 0)
		if err != nil {
			return nil, fmt.Errorf("ошибка подключения к Kafka: %w", err)
		}
		conn.Close()
		log.Printf("Подключение к Kafka успешно установлено")
	}

	return client, nil
}

func (k *KafkaClient) Close() error {
	var producerErr, consumerErr error

	if k.isProducer && k.writer != nil {
		log.Println("Закрытие продюсера Kafka...")
		producerErr = k.writer.Close()
		if producerErr == nil {
			log.Println("Продюсер Kafka успешно закрыт")
		}
	}

	if k.isConsumer && k.reader != nil {
		log.Println("Закрытие потребителя Kafka...")
		consumerErr = k.reader.Close()
		if consumerErr == nil {
			log.Println("Потребитель Kafka успешно закрыт")
		}
	}

	if producerErr != nil {
		return fmt.Errorf("ошибка закрытия продюсера: %w", producerErr)
	}

	if consumerErr != nil {
		return fmt.Errorf("ошибка закрытия потребителя: %w", consumerErr)
	}

	return nil
}

func (k *KafkaClient) Publish(ctx context.Context, key string, value []byte) error {
	if !k.isProducer || k.writer == nil {
		return fmt.Errorf("клиент не настроен как продюсер")
	}

	log.Printf("Отправка сообщения с ключом '%s' в топик %s", key, k.config.Topic)

	err := k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	log.Printf("Сообщение с ключом '%s' успешно отправлено", key)
	return nil
}

func (k *KafkaClient) Subscribe(ctx context.Context, handler func(key string, value []byte) error) error {
	if !k.isConsumer || k.reader == nil {
		return fmt.Errorf("клиент не настроен как потребитель")
	}

	log.Printf("Начало прослушивания топика %s (группа: %s)", k.config.Topic, k.config.GroupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("Получен сигнал остановки потребителя Kafka")
			return nil
		default:
			message, err := k.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Ошибка чтения сообщения: %v", err)
				return fmt.Errorf("ошибка чтения сообщения: %w", err)
			}

			log.Printf("Получено сообщение из топика %s (партиция: %d, смещение: %d)",
				message.Topic, message.Partition, message.Offset)

			if err := handler(string(message.Key), message.Value); err != nil {
				log.Printf("Ошибка обработки сообщения: %v", err)
				return fmt.Errorf("ошибка обработки сообщения: %w", err)
			}
		}
	}
}
