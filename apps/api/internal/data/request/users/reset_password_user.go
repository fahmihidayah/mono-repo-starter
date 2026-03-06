package users

type InitialResetPasswordUser struct {
	Email string `json:"email" validate:"required,email" binding:"required,email" example:"fahmi@gmail.com"`
}

type ResetPasswordUser struct {
	Token       string `json:"token" validate:"required" binding:"required" example:"123123123"`
	NewPassword string `json:"new_password" validate:"required,min=8" binding:"required,min=8" example:"Test@1234"`
}
