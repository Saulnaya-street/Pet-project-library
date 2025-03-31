package handler

import (
	"github.com/gorilla/mux"
)

// Router - структура для регистрации маршрутов API
type Router struct {
	bookHandler IBookHandler
	userHandler IUserHandler
}

func NewRouter(bookHandler IBookHandler, userHandler IUserHandler) *Router {
	return &Router{
		bookHandler: bookHandler,
		userHandler: userHandler,
	}
}

func (r *Router) RegisterRoutes(router *mux.Router) {
	// Регистрируем маршруты для книг
	router.HandleFunc("/api/books", r.bookHandler.GetAllBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", r.bookHandler.GetBook).Methods("GET")
	router.HandleFunc("/api/books", r.bookHandler.CreateBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", r.bookHandler.UpdateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", r.bookHandler.DeleteBook).Methods("DELETE")

	// Регистрируем маршруты для пользователей
	router.HandleFunc("/api/users", r.userHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", r.userHandler.GetUser).Methods("GET")
	router.HandleFunc("/api/users", r.userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", r.userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", r.userHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/auth/login", r.userHandler.Login).Methods("POST")
}
