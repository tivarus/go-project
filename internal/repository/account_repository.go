package repository

import (
	"bank-api/internal/models"
	"database/sql"
	"errors"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(account *models.Account) error {
	query := `
		INSERT INTO accounts (user_id, balance, currency)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		account.UserID,
		account.Balance,
		account.Currency,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	return err
}

func (r *AccountRepository) GetAccountByID(id int) (*models.Account, error) {
	query := `
		SELECT id, user_id, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	account := &models.Account{}
	err := r.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return account, err
}

func (r *AccountRepository) UpdateBalance(id int, amount float64) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(query, amount, id)
	return err
}
