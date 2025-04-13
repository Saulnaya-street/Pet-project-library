package main

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/controller"
	"awesomeProject22/db-service/internal/kafka"
	"awesomeProject22/db-service/internal/repository"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {

	dbConfig := repository.Config{
		Host:     getEnvOrDefault("DB_HOST", "db"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		Username: getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "123"),
		DBName:   getEnvOrDefault("DB_NAME", "Library"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	redisConfig := cache.RedisConfig{
		Host:     getEnvOrDefault("REDIS_HOST", "redis"),
		Port:     getEnvOrDefault("REDIS_PORT", "6379"),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       0,
	}

	kafkaConfig := kafka.KafkaConfig{
		Brokers: strings.Split(getEnvOrDefault("KAFKA_BROKERS", "kafka:9092"), ","),
		Topic:   getEnvOrDefault("KAFKA_TOPIC", "library-events"),
		GroupID: getEnvOrDefault("KAFKA_GROUP_ID", "library-service"),
	}

	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	redisClient, err := cache.NewRedisClient(redisConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %s", err.Error())
	}
	defer redisClient.Close()
	log.Println("Successfully connected to Redis")

	kafkaClient, err := kafka.NewKafkaClient(kafkaConfig, true, true)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %s", err.Error())
	}
	defer kafkaClient.Close()
	log.Println("Successfully connected to Kafka")

	ctrl := controller.NewControllerWithKafka(controller.ControllerWithKafkaOptions{
		DB:          db,
		RedisClient: redisClient,
		KafkaClient: kafkaClient,
	})

	ctrl.StartKafkaConsumer()

	srv := ctrl.GetServer()
	port := getEnvOrDefault("PORT", "8080")

	go func() {
		if err := srv.Run(port); err != nil {
			log.Fatalf("Error occurred while running http server: %s", err.Error())
		}
	}()

	log.Printf("Server started on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Server shutting down...")

	ctrl.StopKafkaConsumer()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error occurred on server shutting down: %s", err.Error())
	}

	ctrl.CloseConnections()

	log.Print("Server successfully stopped")
}
