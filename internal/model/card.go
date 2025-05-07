package model

type Card struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`
}

// CardPlain represents a freshly created virtual card with plain CVV
// swagger:model CardPlain
type CardPlain struct {
	// Unique ID of the card
	// example: 102
	ID int64 `json:"id"`
	// Account ID this card is linked to
	// example: 42
	AccountID int64 `json:"account_id"`
	// Plaintext card number (shown once)
	// example: 4276123456789012
	Number string `json:"number"`
	// Expiry date in MM/YY format
	// example: 12/28
	Expiry string `json:"expiry"`
	// CVV code (shown only once after creation)
	// example: 319
	CVV string `json:"cvv"`
}

// CardRaw represents an encrypted virtual card stored in the system
// swagger:model CardRaw
type CardRaw struct {
	// Unique ID of the card
	// example: 101
	ID int64 `json:"id"`
	// Account ID this card is linked to
	// example: 42
	AccountID int64 `json:"account_id"`
	// Encrypted card number
	// example: 4276123456789012
	Number string `json:"number"`
	// Encrypted expiry date
	// example: 12/28
	Expiry string `json:"expiry"`
	// HMAC signature of the card
	// example: d41d8cd98f00b204e9800998ecf8427e
	HMAC string `json:"hmac"`
}
