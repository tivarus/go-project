package models

import "time"

type Account struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	Currency string `json:"currency" validate:"required,oneof=RUB"`
}

type AccountResponse struct {
	ID        int       `json:"id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateBalanceRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}
