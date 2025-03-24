package handler

import (
	"awesomeProject22/db-service/internal/domain"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
	userService domain.UserService
	router      *mux.Router
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	handler := &UserHandler{
		userService: userService,
		router:      mux.NewRouter(),
	}

	handler.registerRoutes()

	return handler
}

func (h *UserHandler) registerRoutes() {
	h.router.HandleFunc("/api/users", h.GetAllUsers).Methods("GET")
	h.router.HandleFunc("/api/users/{id}", h.GetUser).Methods("GET")
	h.router.HandleFunc("/api/users", h.CreateUser).Methods("POST")
	h.router.HandleFunc("/api/users/{id}", h.UpdateUser).Methods("PUT")
	h.router.HandleFunc("/api/users/{id}", h.DeleteUser).Methods("DELETE")
	h.router.HandleFunc("/api/auth/login", h.Login).Methods("POST")

}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &domain.User{
		Username: userInput.Username,
		Email:    userInput.Email,
		IsAdmin:  userInput.IsAdmin,
	}

	if err := h.userService.Create(user, userInput.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var userInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		IsAdmin  bool   `json:"is_admin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUser, err := h.userService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	currentUser.Username = userInput.Username
	currentUser.Email = userInput.Email
	currentUser.IsAdmin = userInput.IsAdmin

	var passwordUpdated bool
	if userInput.Password != "" {
		passwordUpdated = true
	}

	if err := h.userService.Update(currentUser, passwordUpdated, userInput.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.userService.Authenticate(loginInput.Username, loginInput.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Очищаем хеши паролей перед отправкой
	for _, user := range users {
		user.PasswordHash = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
