package dto

type RegisterRequest struct {
	Username string `json:"Username" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"Username" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
