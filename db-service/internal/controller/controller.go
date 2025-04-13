package controller

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/delivery/handler"
	"awesomeProject22/db-service/internal/kafka"
	"awesomeProject22/db-service/internal/repository"
	"awesomeProject22/db-service/internal/service"
	"awesomeProject22/db-service/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ControllerOptions struct {
	DB          *pgxpool.Pool
	RedisClient cache.IRedisClient
	KafkaClient kafka.IKafkaClient
}

type Controller struct {
	db            *pgxpool.Pool
	redisClient   cache.IRedisClient
	kafkaClient   kafka.IKafkaClient
	bookService   service.IBookService
	userService   service.IUserService
	server        *pkg.Server
	bookHandler   handler.IBookHandler
	userHandler   handler.IUserHandler
	eventProducer kafka.IEventProducer
}

func NewController(opts ControllerOptions) *Controller {
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

	return &Controller{
		db:            opts.DB,
		redisClient:   opts.RedisClient,
		kafkaClient:   opts.KafkaClient,
		bookService:   bookService,
		userService:   userService,
		server:        server,
		bookHandler:   bookHandler,
		userHandler:   userHandler,
		eventProducer: eventProducer,
	}
}

func (c *Controller) GetServer() *pkg.Server {
	return c.server
}

func (c *Controller) CloseConnections() {
	if c.kafkaClient != nil {
		if err := c.kafkaClient.Close(); err != nil {
			log.Printf("Error closing connection to Kafka: %v", err)
		}
	}

	if c.redisClient != nil {
		if err := c.redisClient.Close(); err != nil {
			log.Printf("Error closing connection to Redis: %v", err)
		}
	}

	if c.db != nil {
		c.db.Close()
	}

	log.Println("All connections are closed")
}
