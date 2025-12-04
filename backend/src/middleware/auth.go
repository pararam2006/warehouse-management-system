package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	contextUserIDKey = "userID"
	contextRoleKey   = "role"
)

// authClaims описывает часть пейлоада JWT, которую мы используем в middleware.
// Поля должны совпадать с теми, что устанавливаются в сервисе аутентификации.
type authClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type errorResponse struct {
	Error string `json:"error"`
}

// AuthMiddleware проверяет JWT-токен в заголовке Authorization и,
// если он валиден, добавляет userID и роль в контекст запроса.
func AuthMiddleware(next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем preflight-запросы CORS.
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w, "missing Authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(w, "invalid Authorization header format")
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			unauthorized(w, "empty token")
			return
		}

		claims := &authClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что используется ожидаемый алгоритм.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			unauthorized(w, "invalid or expired token")
			return
		}

		if claims.UserID == "" {
			unauthorized(w, "invalid token payload")
			return
		}

		// Добавляем данные в контекст для последующих обработчиков.
		ctx := context.WithValue(r.Context(), contextUserIDKey, claims.UserID)
		if claims.Role != "" {
			ctx = context.WithValue(ctx, contextRoleKey, claims.Role)
		}

		next(w, r.WithContext(ctx))
	}
}

// RoleMiddleware ограничивает доступ к обработчику по ролям.
// Ожидается, что AuthMiddleware уже положил роль пользователя в контекст.
func RoleMiddleware(next http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем preflight-запросы CORS.
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		rawRole := r.Context().Value(contextRoleKey)
		role, ok := rawRole.(string)
		if !ok || role == "" {
			forbidden(w, "access denied")
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				next(w, r)
				return
			}
		}

		forbidden(w, "access denied")
	}
}

func unauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

func forbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}


