package dto

// Payload to make transactions with account - deposit or withdrawal
// swagger:model depositRequest
type DepositRequest struct {
	// ID of the target account
	// example: 42
	AccountID int64 `json:"account_id"`
	// Amount of funds
	// example: 1000.1
	Amount float64 `json:"amount"`
}

// Payload to make funds transfer
// swagger:model depositRequest
type TransferRequest struct {
	// ID of sender account
	// example: 42
	FromAccountID int64 `json:"from_account_id"`
	// ID of the recipient account
	// example: 24
	ToAccountID int64 `json:"to_account_id"`
	// Amount of funds
	// example: 1000.1
	Amount float64 `json:"amount"`
}
