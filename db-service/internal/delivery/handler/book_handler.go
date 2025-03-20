package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/service"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/books", h.GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", h.GetBook).Methods("GET")
	r.HandleFunc("/books", h.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", h.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", h.DeleteBook).Methods("DELETE")
}

type BookResponse struct {
	ID     string `json:"id"`
	Genre  string `json:"genre"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type BookRequest struct {
	Genre  string `json:"genre"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {

	filters := make(map[string]interface{})

	if author := r.URL.Query().Get("author"); author != "" {
		filters["author"] = author
	}

	if genre := r.URL.Query().Get("genre"); genre != "" {
		filters["genre"] = genre
	}

	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			filters["year"] = year
		}
	}

	books, err := h.bookService.GetAll(filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch books")
		return
	}

	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = BookResponse{
			ID:     book.ID.String(),
			Genre:  book.Genre,
			Name:   book.Name,
			Author: book.Author,
			Year:   book.Year,
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := h.bookService.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Book not found")
		return
	}

	response := BookResponse{
		ID:     book.ID.String(),
		Genre:  book.Genre,
		Name:   book.Name,
		Author: book.Author,
		Year:   book.Year,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var bookReq BookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book := &domain.Book{
		ID:     uuid.New(),
		Genre:  bookReq.Genre,
		Name:   bookReq.Name,
		Author: bookReq.Author,
		Year:   bookReq.Year,
	}

	if err := h.bookService.Create(book, userID); err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			respondWithError(w, http.StatusForbidden, "Only admins can create books")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create book")
		return
	}

	response := BookResponse{
		ID:     book.ID.String(),
		Genre:  book.Genre,
		Name:   book.Name,
		Author: book.Author,
		Year:   book.Year,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var bookReq BookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book := &domain.Book{
		ID:     id,
		Genre:  bookReq.Genre,
		Name:   bookReq.Name,
		Author: bookReq.Author,
		Year:   bookReq.Year,
	}

	if err := h.bookService.Update(book, userID); err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			respondWithError(w, http.StatusForbidden, "Only admins can update books")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Book not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update book")
		return
	}

	response := BookResponse{
		ID:     book.ID.String(),
		Genre:  book.Genre,
		Name:   book.Name,
		Author: book.Author,
		Year:   book.Year,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	if err := h.bookService.Delete(id, userID); err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			respondWithError(w, http.StatusForbidden, "Only admins can delete books")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Book not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete book")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
