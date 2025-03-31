package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
}

func NewServer() *Server {
	router := mux.NewRouter()

	return &Server{
		router: router,
	}
}

func (s *Server) GetRouter() *mux.Router {
	return s.router
}

func (s *Server) SetupRoutes(handlers map[string]http.Handler) {
	for prefix, handler := range handlers {
		s.router.PathPrefix(prefix).Handler(handler)
	}
}

func (s *Server) Run(port string) error {
	// Создаем HTTP сервер
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        s.router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Запускаем сервер
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
