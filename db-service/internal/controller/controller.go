package controller

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/delivery/handler"
	"awesomeProject22/db-service/internal/kafka"
	"awesomeProject22/db-service/internal/repository"
	"awesomeProject22/db-service/internal/service"
	"awesomeProject22/db-service/pkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ControllerWithKafkaOptions struct {
	DB          *pgxpool.Pool
	RedisClient cache.IRedisClient
	KafkaClient kafka.IKafkaClient
}

type ControllerWithKafka struct {
	db             *pgxpool.Pool
	redisClient    cache.IRedisClient
	kafkaClient    kafka.IKafkaClient
	bookService    service.IBookService
	userService    service.IUserService
	server         *pkg.Server
	bookHandler    handler.IBookHandler
	userHandler    handler.IUserHandler
	eventProducer  kafka.IEventProducer
	eventConsumer  *kafka.EventConsumer
	shutdownSignal chan struct{}
}

func NewControllerWithKafka(opts ControllerWithKafkaOptions) *ControllerWithKafka {
	var bookRepo repository.IBookRepository
	var userRepo repository.IUserRepository

	baseBookRepo := repository.NewBookRepository(opts.DB)
	baseUserRepo := repository.NewUserRepository(opts.DB)

	if opts.RedisClient != nil {
		bookRepo = repository.NewCachedBookRepository(baseBookRepo, opts.RedisClient)
		userRepo = repository.NewCachedUserRepository(baseUserRepo, opts.RedisClient)
	} else {
		bookRepo = baseBookRepo
		userRepo = baseUserRepo
	}

	eventProducer := kafka.NewEventProducer(opts.KafkaClient)

	bookService := service.NewBookServiceWithKafka(bookRepo, eventProducer)
	userService := service.NewUserServiceWithKafka(userRepo, eventProducer)

	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)

	server := pkg.NewServer()

	deliveryRouter := handler.NewRouter(bookHandler, userHandler)
	deliveryRouter.RegisterRoutes(server.GetRouter())

	eventConsumer := kafka.NewEventConsumer(opts.KafkaClient)

	return &ControllerWithKafka{
		db:             opts.DB,
		redisClient:    opts.RedisClient,
		kafkaClient:    opts.KafkaClient,
		bookService:    bookService,
		userService:    userService,
		server:         server,
		bookHandler:    bookHandler,
		userHandler:    userHandler,
		eventProducer:  eventProducer,
		eventConsumer:  eventConsumer,
		shutdownSignal: make(chan struct{}),
	}
}

func (c *ControllerWithKafka) GetServer() *pkg.Server {
	return c.server
}

func (c *ControllerWithKafka) StartKafkaConsumer() {

	c.registerEventHandlers()

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			<-c.shutdownSignal
			cancel()
		}()

		if err := c.eventConsumer.Start(ctx); err != nil {
			log.Printf("Ошибка запуска потребителя Kafka: %v", err)
		}
	}()

	log.Println("Потребитель Kafka запущен")
}

func (c *ControllerWithKafka) StopKafkaConsumer() {
	close(c.shutdownSignal)
	log.Println("Остановка потребителя Kafka...")
}

func (c *ControllerWithKafka) CloseConnections() {

	c.StopKafkaConsumer()

	if c.kafkaClient != nil {
		if err := c.kafkaClient.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с Kafka: %v", err)
		}
	}

	if c.redisClient != nil {
		if err := c.redisClient.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с Redis: %v", err)
		}
	}

	if c.db != nil {
		c.db.Close()
	}

	log.Println("Все соединения закрыты")
}

func (c *ControllerWithKafka) registerEventHandlers() {

	c.eventConsumer.RegisterHandler(kafka.BookCreated, func(ctx context.Context, event kafka.Event) error {
		log.Printf("Получено событие создания книги: %v", event.ID)

		var bookEvent kafka.BookEvent
		bookData, err := json.Marshal(event.Payload)
		if err != nil {
			log.Printf("Ошибка сериализации данных книги: %v", err)
			return err
		}

		if err := json.Unmarshal(bookData, &bookEvent); err != nil {
			log.Printf("Ошибка десериализации данных книги: %v", err)
			return err
		}

		log.Printf("Книга создана: %s, автор: %s", bookEvent.Book.Name, bookEvent.Book.Author)

		return nil
	})

	c.eventConsumer.RegisterHandler(kafka.BookUpdated, func(ctx context.Context, event kafka.Event) error {
		log.Printf("Получено событие обновления книги: %v", event.ID)

		var bookEvent kafka.BookEvent
		bookData, err := json.Marshal(event.Payload)
		if err != nil {
			log.Printf("Ошибка сериализации данных книги: %v", err)
			return err
		}

		if err := json.Unmarshal(bookData, &bookEvent); err != nil {
			log.Printf("Ошибка десериализации данных книги: %v", err)
			return err
		}

		log.Printf("Книга обновлена: %s, автор: %s", bookEvent.Book.Name, bookEvent.Book.Author)

		return nil
	})

	c.eventConsumer.RegisterHandler(kafka.BookDeleted, func(ctx context.Context, event kafka.Event) error {
		log.Printf("Получено событие удаления книги: %v", event.ID)

		payload, ok := event.Payload.(map[string]interface{})
		if !ok {
			return fmt.Errorf("неверный формат payload")
		}

		bookID, ok := payload["id"].(string)
		if !ok {
			return fmt.Errorf("неверный формат id книги")
		}

		log.Printf("Книга удалена: %s", bookID)

		return nil
	})

	c.eventConsumer.RegisterHandler(kafka.UserLoggedIn, func(ctx context.Context, event kafka.Event) error {
		log.Printf("Получено событие входа пользователя: %v", event.ID)

		var loginEvent kafka.LoginEvent
		loginData, err := json.Marshal(event.Payload)
		if err != nil {
			log.Printf("Ошибка сериализации данных входа: %v", err)
			return err
		}

		if err := json.Unmarshal(loginData, &loginEvent); err != nil {
			log.Printf("Ошибка десериализации данных входа: %v", err)
			return err
		}

		log.Printf("Пользователь вошел в систему: %s (ID: %s)", loginEvent.Username, loginEvent.UserID)

		return nil
	})

	log.Println("Все обработчики событий успешно зарегистрированы")
}
