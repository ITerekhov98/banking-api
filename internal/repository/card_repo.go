package repository

import (
	"context"

	"banking-api/internal/db"
	"banking-api/internal/model"
	"banking-api/internal/security"
)

type CardRepository struct{}

func (r *CardRepository) Create(
	ctx context.Context, userID, accountID int64, number,
	expiry []byte, cvvHash, hmac string) (*model.Card, error) {
	query := `
		INSERT INTO cards (user_id, account_id, number, expiry, cvv_hash, hmac)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	var id int64
	err := db.DB.QueryRowContext(ctx, query, userID, accountID, number, expiry, cvvHash, hmac).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &model.Card{ID: id, UserID: userID, AccountID: accountID}, nil
}

func (r *CardRepository) GetByUser(ctx context.Context, userID int64) ([]*model.CardRaw, error) {
	query := `
		SELECT id, account_id, number, expiry, hmac
		FROM cards
		WHERE user_id = $1
	`
	rows, err := db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*model.CardRaw
	for rows.Next() {
		var card model.CardRaw
		var numberEnc, expiryEnc []byte
		err := rows.Scan(&card.ID, &card.AccountID, &numberEnc, &expiryEnc, &card.HMAC)
		if err != nil {
			return nil, err
		}

		// encrypt
		number, err := security.DecryptPGP(numberEnc)
		if err != nil {
			return nil, err
		}
		expiry, err := security.DecryptPGP(expiryEnc)
		if err != nil {
			return nil, err
		}
		card.Number = number
		card.Expiry = expiry
		cards = append(cards, &card)
	}
	return cards, nil
}
