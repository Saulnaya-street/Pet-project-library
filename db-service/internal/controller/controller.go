package controller

import (
	"awesomeProject22/db-service/internal/delivery/handler"
	"awesomeProject22/db-service/internal/repository"
	"awesomeProject22/db-service/internal/service"
	"awesomeProject22/db-service/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Controller - основной контроллер приложения
type Controller struct {
	db          *pgxpool.Pool
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

// GetServer - возвращает HTTP сервер
func (c *Controller) GetServer() *pkg.Server {
	return c.server
}
