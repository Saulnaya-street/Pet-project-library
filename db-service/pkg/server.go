package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	"awesomeProject22/db-service/internal/delivery/handler"
	"awesomeProject22/db-service/internal/repository"
	"awesomeProject22/db-service/internal/service"
)

type Server struct {
	httpServer *http.Server
	db         *sqlx.DB
}

func NewServer(db *sqlx.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Run(port string) error {
	// Инициализируем репозитории
	bookRepo := repository.NewBookRepository(s.db)

	// Инициализируем сервисы (используем мок для пользовательского сервиса)
	userService := &service.MockUserService{}
	bookService := service.NewBookService(bookRepo, userService)

	// Инициализируем обработчики
	bookHandler := handler.NewBookHandler(bookService)

	// Создаем роутер
	router := mux.NewRouter()

	// Регистрируем маршруты
	bookHandler.RegisterRoutes(router)

	// Добавляем CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Создаем HTTP сервер
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        c.Handler(router),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Запускаем сервер
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
