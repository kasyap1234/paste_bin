package models

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"  validate:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
