package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware — простой middleware для логирования HTTP-запросов.
// Подключается через router.Use(middleware.LoggingMiddleware).
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("%s %s %s (%s)", r.RemoteAddr, r.Method, r.URL.Path, duration)
	})
}


