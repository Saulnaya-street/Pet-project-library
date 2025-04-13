package pkg

import (
	"context"
	"fmt"
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

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        s.router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error starting HTTP server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error while stopping server: %w", err)
	}
	return nil
}
