package services

import (
	"errors"
	"time"
	"warehouse-management-system/src/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository описывает поведение хранилища пользователей, необходимое слою сервиса.
// Конкретная реализация будет находиться в пакете repositories.
type UserRepository interface {
	// FindByEmail возвращает пользователя по email или nil, если не найден.
	FindByEmail(email string) (*models.User, error)
	// FindByID возвращает пользователя по идентификатору или nil, если не найден.
	FindByID(id string) (*models.User, error)
	// Create сохраняет нового пользователя в хранилище.
	Create(user *models.User) error
}

// AuthService инкапсулирует бизнес-логику аутентификации и авторизации.
type AuthService struct {
	userRepo  UserRepository
	jwtSecret string
}

// NewAuthService — конструктор сервиса аутентификации.
func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Claims описывает JWT-пейлоад, который мы отдаём клиенту.
type Claims struct {
	UserID string      `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

// Register регистрирует нового пользователя: валидирует данные, хэширует пароль,
// создаёт запись пользователя и возвращает JWT-токен.
func (s *AuthService) Register(email, password string, role models.Role) (*models.User, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyInUse
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.NewUser("", email, passwordHash, role)

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login выполняет вход пользователя: проверяет email/пароль и выдаёт JWT-токен.
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := comparePassword(user.PasswordHash, password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserByID возвращает пользователя по его идентификатору.
func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// generateToken создаёт JWT-токен с данными пользователя.
func (s *AuthService) generateToken(user *models.User) (string, error) {
	now := time.Now().UTC()
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func hashPassword(password string) (string, error) {
	const cost = 12
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
