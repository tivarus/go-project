package models

import (
	"time"
)

type CreditStatus string

const (
	CreditStatusActive  CreditStatus = "active"
	CreditStatusClosed  CreditStatus = "closed"
	CreditStatusOverdue CreditStatus = "overdue"
)

type Credit struct {
	ID           int          `json:"id"`
	AccountID    int          `json:"account_id"`
	Amount       float64      `json:"amount"`
	InterestRate float64      `json:"interest_rate"`
	TermMonths   int          `json:"term_months"`
	StartDate    time.Time    `json:"start_date"`
	Status       CreditStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
}

type PaymentSchedule struct {
	ID          int        `json:"id"`
	CreditID    int        `json:"credit_id"`
	PaymentDate time.Time  `json:"payment_date"`
	Amount      float64    `json:"amount"`
	Principal   float64    `json:"principal"`
	Interest    float64    `json:"interest"`
	Status      string     `json:"status"`
	PaidAt      *time.Time `json:"paid_at"`
}

type CreateCreditRequest struct {
	AccountID  int     `json:"account_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	TermMonths int     `json:"term_months" validate:"required,gte=1,lte=60"`
}

type CreditResponse struct {
	ID           int          `json:"id"`
	Amount       float64      `json:"amount"`
	InterestRate float64      `json:"interest_rate"`
	TermMonths   int          `json:"term_months"`
	StartDate    time.Time    `json:"start_date"`
	Status       CreditStatus `json:"status"`
}
