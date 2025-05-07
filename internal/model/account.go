package model

// Account contain info about created account
// swagger:model Account
type Account struct {
	// Unique ID of the account
	// example: 102
	ID int64 `json:"id"`
	// Unique userID of the account
	// example: 102
	UserID int64 `json:"user_id"`
	// Current account balance
	// example: 1000
	Balance float64 `json:"balance"`
}
