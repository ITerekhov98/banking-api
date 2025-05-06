package model

type Card struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`
}

type CardPlain struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"account_id"`
	Number    string `json:"number"`
	Expiry    string `json:"expiry"`
	CVV       string `json:"cvv"` // return only once after card creation
}

type CardRaw struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"account_id"`
	Number    string `json:"number"`
	Expiry    string `json:"expiry"`
	HMAC      string `json:"hmac"`
}
