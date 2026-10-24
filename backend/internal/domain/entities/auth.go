package entities

import "time"

type SignUp struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=255"`
	LastName  string `json:"last_name" validate:"required,min=1,max=255"`
	Email     string `json:"email" validate:"required,max=255,email"`
	Password  string `json:"password" validate:"required,min=8,max=72,printascii"`
}

type VerificationRequest struct {
	Email            string `json:"email" validate:"required,max=255,email"`
	VerificationCode string `json:"verification_code" validate:"required,min=6,max=6,number"`
}

type ResendCodeRequest struct {
	Email string `json:"email" validate:"required,max=255,email"`
}

type LoginParameter struct {
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,min=8,max=72,printascii"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      LoginUser `json:"user"`
}

type LoginUser struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
