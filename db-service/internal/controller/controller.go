package controller

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/delivery/handler"
	"awesomeProject22/db-service/internal/repository"
	"awesomeProject22/db-service/internal/service"
	"awesomeProject22/db-service/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Controller - основной контроллер приложения
type Controller struct {
	db          *pgxpool.Pool
	redisClient *cache.RedisClient
	bookService service.IBookService
	userService service.IUserService
	server      *pkg.Server
	bookHandler handler.IBookHandler
	userHandler handler.IUserHandler
}

// NewController - конструктор контроллера
func NewController(db *pgxpool.Pool) *Controller {
	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)

	bookService := service.NewBookService(bookRepo)
	userService := service.NewUserService(userRepo)

	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)

	server := pkg.NewServer()

	deliveryRouter := handler.NewRouter(bookHandler, userHandler)

	deliveryRouter.RegisterRoutes(server.GetRouter())

	return &Controller{
		db:          db,
		bookService: bookService,
		userService: userService,
		server:      server,
		bookHandler: bookHandler,
		userHandler: userHandler,
	}
}

// NewControllerWithRedis - конструктор контроллера с поддержкой Redis
func NewControllerWithRedis(db *pgxpool.Pool, redisClient *cache.RedisClient) *Controller {
	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Создаем репозитории с кешированием
	cachedBookRepo := repository.NewCachedBookRepository(bookRepo, redisClient)
	cachedUserRepo := repository.NewCachedUserRepository(userRepo, redisClient)

	bookService := service.NewBookService(cachedBookRepo)
	userService := service.NewUserService(cachedUserRepo)

	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)

	server := pkg.NewServer()

	deliveryRouter := handler.NewRouter(bookHandler, userHandler)

	deliveryRouter.RegisterRoutes(server.GetRouter())

	return &Controller{
		db:          db,
		redisClient: redisClient,
		bookService: bookService,
		userService: userService,
		server:      server,
		bookHandler: bookHandler,
		userHandler: userHandler,
	}
}

// GetServer - возвращает HTTP сервер
func (c *Controller) GetServer() *pkg.Server {
	return c.server
}
