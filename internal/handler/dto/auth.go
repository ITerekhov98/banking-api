package dto

// registerRequest contain user payload to register in system
// swagger:model registerRequest
type RegisterRequest struct {
	// Unique user email
	// example: example@gmail.com
	Email string `json:"email"`
	// Unique user name
	// example: user_name
	Username string `json:"username"`
	// Password
	// example: P_ass11worD
	Password string `json:"password"`
}

// RegisterResponse contain user payload after succesful register
// swagger:model RegisterResponse
type RegisterResponse struct {
	// ID of created user
	// example: 1
	Id int64 `json:"id"`
	// Unique user email
	// example: example@gmail.com
	Email string `json:"email"`
	// Username
	// example: user_name
	Username string `json:"username"`
}

// LoginRequest contain user credentials to log in and get auth access token
// swagger:model LoginRequest
type LoginRequest struct {
	// Unique user email
	// example: example@gmail.com
	Email string `json:"email"`
	// Password
	// example: P_ass11worD
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
