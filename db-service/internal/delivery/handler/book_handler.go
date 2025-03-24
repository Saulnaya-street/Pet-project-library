package handler

import (
	"awesomeProject22/db-service/internal/domain"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type BookHandler struct {
	bookService domain.BookService // поле которое хранит ссылку на сервисный слой
	router      *mux.Router
}

func NewBookHandler(bookService domain.BookService) *BookHandler {
	handler := &BookHandler{
		bookService: bookService,
		router:      mux.NewRouter(),
	}

	handler.registerRoutes()

	return handler
}

func (h *BookHandler) registerRoutes() {
	h.router.HandleFunc("/api/books", h.GetAllBooks).Methods("GET")
	h.router.HandleFunc("/api/books/{id}", h.GetBook).Methods("GET")
	h.router.HandleFunc("/api/books", h.CreateBook).Methods("POST")
	h.router.HandleFunc("/api/books/{id}", h.UpdateBook).Methods("PUT")
	h.router.HandleFunc("/api/books/{id}", h.DeleteBook).Methods("DELETE")

}

// ServeHTTP реализует интерфейс http.Handler
func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	// Получаем конкретные параметры фильтрации напрямую
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")

	books, err := h.bookService.GetAll(author, genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.bookService.Create(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = id

	if err := h.bookService.Update(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.bookService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
