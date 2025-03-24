package pkg

import (
	"awesomeProject22/db-service/internal/delivery/handler"
	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/repository"
	"awesomeProject22/db-service/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type Controller struct {
	db          *pgxpool.Pool
	bookService domain.BookService
	userService domain.UserService
	server      *Server
	bookHandler http.Handler
	userHandler http.Handler
}

func NewController(db *pgxpool.Pool) *Controller {

	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)

	bookService := service.NewBookService(bookRepo)
	userService := service.NewUserService(userRepo)

	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)

	server := NewServer()

	return &Controller{
		db:          db,
		bookService: bookService,
		userService: userService,
		server:      server,
		bookHandler: bookHandler,
		userHandler: userHandler,
	}
}

func (c *Controller) InitRoutes() {

	handlers := map[string]http.Handler{
		"/api/books": c.bookHandler,
		"/api/users": c.userHandler,
		"/api/auth":  c.userHandler,
	}

	c.server.SetupRoutes(handlers)
}

func (c *Controller) GetServer() *Server {
	return c.server
}
