package repository

import (
	"bank-api/internal/models"
	"database/sql"
	"errors"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) CreateCard(card *models.Card) error {
	query := `
		INSERT INTO cards (account_id, card_number, expiry_date, cvv_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		card.AccountID,
		card.Number,
		card.ExpiryDate,
		card.CVV,
	).Scan(&card.ID, &card.CreatedAt)

	return err
}

func (r *CardRepository) GetCardByID(id int) (*models.Card, error) {
	query := `
		SELECT id, account_id, card_number, expiry_date, created_at
		FROM cards
		WHERE id = $1
	`

	card := &models.Card{}
	err := r.db.QueryRow(query, id).Scan(
		&card.ID,
		&card.AccountID,
		&card.Number,
		&card.ExpiryDate,
		&card.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return card, err
}

func (r *CardRepository) GetCardsByAccount(accountID int) ([]*models.Card, error) {
	query := `
		SELECT id, account_id, card_number, expiry_date, created_at
		FROM cards
		WHERE account_id = $1
	`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		if err := rows.Scan(
			&card.ID,
			&card.AccountID,
			&card.Number,
			&card.ExpiryDate,
			&card.CreatedAt,
		); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}
