package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"warehouse-management-system/src/models"
	"warehouse-management-system/src/services"
)

// AuthController обрабатывает HTTP-запросы, связанные с аутентификацией.
// Он ничего не знает о БД и HTTP-роутере — только о входящих/исходящих DTO и сервисе.
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController — конструктор контроллера аутентификации.
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// registerRequest описывает тело запроса на регистрацию.
type registerRequest struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     models.Role `json:"role"`
}

// loginRequest описывает тело запроса на вход.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// authResponse — стандартный ответ при успешной аутентификации/регистрации.
type authResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Register — HTTP-обработчик регистрации пользователя.
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || len(req.Password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "email must be set and password must be at least 6 characters"})
		return
	}

	user, err := c.authService.Register(req.Email, req.Password, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyInUse):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

// Login — HTTP-обработчик входа пользователя.
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "email and password are required"})
		return
	}

	user, token, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(authResponse{
		Token: token,
		User:  user,
	})
}

// GetMe — HTTP-обработчик, возвращающий текущего пользователя по информации из контекста.
// Ожидается, что middleware аутентификации положит userID в контекст запроса.
func (c *AuthController) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// В middleware мы будем использовать этот ключ для сохранения ID пользователя в контекст.
	const contextUserIDKey = "userID"

	rawID := r.Context().Value(contextUserIDKey)
	userID, ok := rawID.(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	user, err := c.authService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
