package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email" validate:"required,email"`
	Username     string    `json:"username" validate:"required,min=3,max=50"`
	PasswordHash string    `json:"-"` // Исключаем из JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
