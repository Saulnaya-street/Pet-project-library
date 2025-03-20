package main

import (
	"awesomeProject22/db-service/pkg"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"awesomeProject22/db-service/internal/repository"
)

func main() {
	// Получаем параметры подключения из переменных окружения
	dbConfig := repository.Config{
		Host:     getEnv("DB_HOST", "db"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "123"),
		DBName:   getEnv("DB_NAME", "Library"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Подключаемся к базе данных с повторными попытками
	var db *sqlx.DB
	var err error

	maxAttempts := 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Printf("Попытка подключения к базе данных %d/%d", attempt, maxAttempts)

		db, err = repository.NewPostgresDB(dbConfig)
		if err == nil {
			break
		}

		log.Printf("Не удалось подключиться: %s", err.Error())

		if attempt < maxAttempts {
			log.Printf("Повторная попытка через 3 секунды...")
			time.Sleep(3 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("Failed to initialize db after multiple attempts: %s", err.Error())
	}

	defer db.Close()

	log.Println("Successfully connected to database")

	// Создаем и запускаем сервер
	srv := pkg.NewServer(db)

	port := getEnv("PORT", "8080")

	go func() {
		if err := srv.Run(port); err != nil {
			log.Fatalf("Error occurred while running http server: %s", err.Error())
		}
	}()

	log.Printf("Server started on port %s", port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error occurred on server shutting down: %s", err.Error())
	}

	log.Print("Server successfully stopped")
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
