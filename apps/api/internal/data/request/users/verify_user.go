package users

type VerifyUserRequest struct {
	Email            string `json:"email" validate:"omitempty,email" binding:"omitempty,email"`
	VerificationCode string `json:"verification_code" validate:"required,min=2,max=100" binding:"omitempty,min=2,max=100"`
}
