package entities

type SignUp struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=255"`
	LastName  string `json:"last_name" validate:"required,min=1,max=255"`
	Email     string `json:"email" validate:"required,max=255,email"`
	Password  string `json:"password" validate:"required,min=8,max=72,printascii"`
}
