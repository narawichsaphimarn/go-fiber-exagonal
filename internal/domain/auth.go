package domain

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
