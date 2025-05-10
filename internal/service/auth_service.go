package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"bank-api/pkg/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Проверка уникальности email
	exists, err := s.userRepo.IsEmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создание пользователя
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (string, error) {
	// Получение пользователя
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Генерация токена
	token, err := auth.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
