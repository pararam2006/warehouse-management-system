package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит настройки приложения, считываемые из переменных окружения.
// Отдельный слой конфигурации упрощает тестирование и переключение окружений.
type Config struct {
	// JWTSecret используется для подписи и проверки JWT-токенов.
	JWTSecret string
	DBPath    string
}

// LoadConfig инициализирует конфигурацию приложения.
// Все значения берутся из переменных окружения, при отсутствии — подставляются безопасные дефолты для разработки.
func LoadConfig() *Config {
	// Загружаем переменные окружения из файла .env в корне backend.
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("INFO: .env file not found or cannot be loaded: %v\n", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("WARNING: env JWT_SECRET is not set, using insecure default value for development")
		jwtSecret = "change-me-in-prod"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Println("WARNING: env DB_PATH is not set, using default backend/warehouse.db")
		dbPath = "backend/warehouse.db"
	}

	return &Config{
		JWTSecret: jwtSecret,
		DBPath:    dbPath,
	}
}
