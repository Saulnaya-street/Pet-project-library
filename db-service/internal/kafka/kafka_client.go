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
}

type KafkaClient struct {
	writer *kafka.Writer
	config KafkaConfig
}

func NewKafkaClient(cfg KafkaConfig) (IKafkaClient, error) {
	client := &KafkaClient{
		config: cfg,
	}

	client.writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
	log.Printf("Configured Kafka producer for topic %s", cfg.Topic)

	log.Printf("Checking the connection to the Kafka broker %s...", cfg.Brokers[0])
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Brokers[0], cfg.Topic, 0)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Kafka: %w", err)
	}
	conn.Close()
	log.Printf("Connection to Kafka established successfully")

	return client, nil
}

func (k *KafkaClient) Close() error {
	if k.writer != nil {
		log.Println("Kafka Producer Shutdown...")
		err := k.writer.Close()
		if err != nil {
			return fmt.Errorf("producer closing error: %w", err)
		}
		log.Println("Kafka producer successfully closed")
	}

	return nil
}

func (k *KafkaClient) Publish(ctx context.Context, key string, value []byte) error {
	if k.writer == nil {
		return fmt.Errorf("the client is not configured as a producer")
	}

	log.Printf("Sending a message with a key '%s' to topic %s", key, k.config.Topic)

	err := k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	log.Printf("Key message '%s' successfully sent", key)
	return nil
}
