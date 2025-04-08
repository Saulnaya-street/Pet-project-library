package controller

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/delivery/handler"
	"awesomeProject22/db-service/internal/repository"
	"awesomeProject22/db-service/internal/service"
	"awesomeProject22/db-service/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ControllerOptions struct {
	DB          *pgxpool.Pool
	RedisClient cache.IRedisClient
}

type Controller struct {
	db          *pgxpool.Pool
	redisClient cache.IRedisClient
	bookService service.IBookService
	userService service.IUserService
	server      *pkg.Server
	bookHandler handler.IBookHandler
	userHandler handler.IUserHandler
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

	bookService := service.NewBookService(bookRepo)
	userService := service.NewUserService(userRepo)

	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)

	server := pkg.NewServer()

	deliveryRouter := handler.NewRouter(bookHandler, userHandler)
	deliveryRouter.RegisterRoutes(server.GetRouter())

	return &Controller{
		db:          opts.DB,
		redisClient: opts.RedisClient,
		bookService: bookService,
		userService: userService,
		server:      server,
		bookHandler: bookHandler,
		userHandler: userHandler,
	}
}

func (c *Controller) GetServer() *pkg.Server {
	return c.server
}
