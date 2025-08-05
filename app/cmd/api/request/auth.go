package requests

type RegisterUserRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=200"`
	Password string `json:"password" validate:"required,min=7,max=20"`
}

type LoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
