package dto

// createCardRequest represents a card creation request
// swagger:model createCardRequest
type CreateCardRequest struct {
	// ID of the account to associate the credit with
	// example: 42
	AccountID int64 `json:"account_id"`
}
