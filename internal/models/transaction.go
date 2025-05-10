package models

import "time"

type TransactionType string

const (
	TransactionDeposit    TransactionType = "deposit"
	TransactionWithdrawal TransactionType = "withdrawal"
	TransactionTransfer   TransactionType = "transfer"
)

type Transaction struct {
	ID          int             `json:"id"`
	AccountID   int             `json:"account_id"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}

type TransferRequest struct {
	FromAccountID int     `json:"from_account_id" validate:"required"`
	ToAccountID   int     `json:"to_account_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Description   string  `json:"description"`
}
