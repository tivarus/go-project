package models

import "time"

type Card struct {
	ID         int       `json:"id"`
	AccountID  int       `json:"account_id"`
	Number     string    `json:"number"`
	ExpiryDate string    `json:"expiry_date"`
	CVV        string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateCardRequest struct {
	AccountID int `json:"account_id" validate:"required"`
}

type CardResponse struct {
	ID         int       `json:"id"`
	AccountID  int       `json:"account_id"`
	LastFour   string    `json:"last_four"`
	ExpiryDate string    `json:"expiry_date"`
	CreatedAt  time.Time `json:"created_at"`
}
