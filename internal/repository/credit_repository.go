package repository

import (
	"bank-api/internal/models"
	"database/sql"
	"errors"
)

type CreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{db: db}
}

func (r *CreditRepository) CreateCredit(credit *models.Credit) error {
	query := `
		INSERT INTO credits (account_id, amount, interest_rate, term_months, start_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		credit.AccountID,
		credit.Amount,
		credit.InterestRate,
		credit.TermMonths,
		credit.StartDate,
		credit.Status,
	).Scan(&credit.ID, &credit.CreatedAt)
}

func (r *CreditRepository) CreatePaymentSchedule(tx *sql.Tx, payment *models.PaymentSchedule) error {
	query := `
		INSERT INTO payment_schedules (credit_id, payment_date, amount, principal, interest, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	return tx.QueryRow(
		query,
		payment.CreditID,
		payment.PaymentDate,
		payment.Amount,
		payment.Principal,
		payment.Interest,
		payment.Status,
	).Scan(&payment.ID)
}

func (r *CreditRepository) GetCreditByID(id int) (*models.Credit, error) {
	query := `
        SELECT id, account_id, amount, interest_rate, term_months, start_date, status, created_at
        FROM credits
        WHERE id = $1
    `

	credit := &models.Credit{}
	err := r.db.QueryRow(query, id).Scan(
		&credit.ID,
		&credit.AccountID,
		&credit.Amount,
		&credit.InterestRate,
		&credit.TermMonths,
		&credit.StartDate,
		&credit.Status,
		&credit.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return credit, err
}

func (r *CreditRepository) GetPaymentSchedule(creditID int) ([]*models.PaymentSchedule, error) {
	query := `
        SELECT id, credit_id, payment_date, amount, principal, interest, status, paid_at
        FROM payment_schedules
        WHERE credit_id = $1
        ORDER BY payment_date
    `

	rows, err := r.db.Query(query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.PaymentSchedule
	for rows.Next() {
		payment := &models.PaymentSchedule{}
		if err := rows.Scan(
			&payment.ID,
			&payment.CreditID,
			&payment.PaymentDate,
			&payment.Amount,
			&payment.Principal,
			&payment.Interest,
			&payment.Status,
			&payment.PaidAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}
