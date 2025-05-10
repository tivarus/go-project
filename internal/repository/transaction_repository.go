package repository

import (
	"bank-api/internal/models"
	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(tx *sql.Tx, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (account_id, amount, type, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var err error
	if tx != nil {
		err = tx.QueryRow(
			query,
			transaction.AccountID,
			transaction.Amount,
			transaction.Type,
			transaction.Description,
		).Scan(&transaction.ID, &transaction.CreatedAt)
	} else {
		err = r.db.QueryRow(
			query,
			transaction.AccountID,
			transaction.Amount,
			transaction.Type,
			transaction.Description,
		).Scan(&transaction.ID, &transaction.CreatedAt)
	}

	return err
}

func (r *TransactionRepository) GetTransactionsByAccount(accountID int) ([]*models.Transaction, error) {
	query := `
		SELECT id, account_id, amount, type, description, created_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		if err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.Amount,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
