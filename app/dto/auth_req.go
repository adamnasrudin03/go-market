package dto

type LoginReq struct {
	Username string `json:"username" validate:"required,min=4"` // phone or email
	Password string `json:"password" validate:"required,min=4"`
}

type LoginRes struct {
	Token string `json:"token"`
}

type RegisterReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=6"`
	Address  string `json:"address" validate:"required,min=12"`
}
