package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	request "github.com/fahmihidayah/go-api-orchestrator/internal/data/request/users"
	"github.com/fahmihidayah/go-api-orchestrator/internal/middleware"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
)

func AuthUserControllerProvider(userService service.IUserService) *AuthUserController {
	return &AuthUserController{
		userService: userService,
	}
}

type AuthUserController struct {
	userService service.IUserService
}

// VerifyUser verifies a user's email address
// @Summary Verify user email
// @Description Verify user's email address using the verification token sent via email. Token expires after 24 hours.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.VerifyUserRequest true "Email and verification token"
// @Success 200 {object} response.WebResponse "Email verified successfully"
// @Failure 400 {object} response.WebResponse "Invalid request, token expired, or verification failed"
// @Router /api/users/auth/verify [post]
func (c *AuthUserController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var req request.VerifyUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendBadRequest(w, "Invalid request")
		return
	}

	if err := c.userService.Verify(r.Context(), &req); err != nil {
		utils.SendBadRequest(w, "Verification failed: "+err.Error())
		return
	}

	utils.SendSuccess(w, "Email verified successfully", nil)
}

// RegisterUser handles user registration with email verification
// @Summary Register a new user
// @Description Create a new user account with email, password, and name. Sends a verification email with a 24-hour expiry link.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "User registration details"
// @Success 201 {object} response.WebResponse{data=map[string]string} "User registered successfully. Check your email for verification link."
// @Failure 400 {object} response.WebResponse "Invalid request body or validation error (e.g., email already in use)"
// @Router /api/users/auth/register [post]
func (c *AuthUserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendBadRequest(w, "Invalid request body")
		return
	}

	// Register user (with email verification)
	userWithToken, err := c.userService.Register(r.Context(), &req)

	if err != nil {
		utils.SendBadRequest(w, "Registration failed: "+err.Error())
		return
	}

	userData := map[string]interface{}{
		"id":         userWithToken.User.ID,
		"email":      userWithToken.User.Email,
		"name":       userWithToken.User.Name,
		"created_at": userWithToken.User.CreatedAt,
		"token":      userWithToken.Token,
		"exp":        userWithToken.Exp,
	}
	// Success response
	utils.SendCreated(w, "User registered successfully. Please check your email to verify your account.", userData)
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginUserRequest true "User login credentials"
// @Success 200 {object} response.WebResponse{data=map[string]interface{}} "Login successful with user data and token"
// @Failure 400 {object} response.WebResponse "Invalid request body"
// @Failure 401 {object} response.WebResponse "Invalid credentials"
// @Router /api/users/auth/login [post]
func (c *AuthUserController) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginUserRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendBadRequest(w, "Invalid request body")
		return
	}

	// Authenticate user
	userResponse, err := c.userService.Login(r.Context(), &req)
	if err != nil {
		utils.SendUnauthorized(w, err.Error())
		return
	}

	user := userResponse.User
	token := userResponse.Token

	// Success response (exclude sensitive data)
	userData := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"created_at": user.CreatedAt,
		"token":      token,
		"exp":        userResponse.Exp,
	}

	utils.SendSuccess(w, "Login successful", userData)
}

// InitialResetPassword initiates a password reset process
// @Summary Initiate password reset
// @Description Request a password reset token for the given email address
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.InitialResetPasswordUser true "Email address for password reset"
// @Success 200 {object} response.WebResponse{data=map[string]string} "Password reset initiated successfully"
// @Failure 400 {object} response.WebResponse "Invalid request body or validation error"
// @Router /api/users/auth/initial-reset-password [post]
func (c *AuthUserController) InitialResetPassword(w http.ResponseWriter, r *http.Request) {
	var req request.InitialResetPasswordUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendBadRequest(w, "Invalid request body")
		return
	}

	resetToken, err := c.userService.InitiatePasswordReset(r.Context(), &req)
	if err != nil {
		utils.SendBadRequest(w, err.Error())
		return
	}

	utils.SendSuccess(w, "Password reset initiated successfully. Check your email for the reset link.", map[string]string{"reset_token": resetToken})
}

// CompleteResetPassword completes the password reset process
// @Summary Complete password reset
// @Description Reset user password using the reset token and new password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ResetPasswordUser true "Reset token and new password"
// @Success 200 {object} response.WebResponse "Password reset successful"
// @Failure 400 {object} response.WebResponse "Invalid request body, token, or validation error"
// @Router /api/users/auth/complete-reset-password [post]
func (c *AuthUserController) CompleteResetPassword(w http.ResponseWriter, r *http.Request) {
	var req request.ResetPasswordUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendBadRequest(w, "Invalid request body")
		return
	}

	// Call service to complete password reset
	if err := c.userService.CompletePasswordReset(r.Context(), &req); err != nil {
		utils.SendBadRequest(w, err.Error())
		return
	}

	utils.SendSuccess(w, "Password reset successful. You can now login with your new password.", nil)
}

// Logout handles user logout by blacklisting the JWT token
// @Summary User logout
// @Description Logout user by invalidating the current JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.WebResponse "Logout successful"
// @Failure 401 {object} response.WebResponse "Unauthorized - missing or invalid token"
// @Failure 500 {object} response.WebResponse "Failed to logout"
// @Router /api/users/auth/logout [post]
func (c *AuthUserController) Logout(w http.ResponseWriter, r *http.Request) {
	// Get Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.SendUnauthorized(w, "Authorization header required")
		return
	}

	// Extract token from Bearer header
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.SendUnauthorized(w, "Invalid authorization header format")
		return
	}

	token := parts[1]
	if token == "" {
		utils.SendUnauthorized(w, "Token is empty")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendUnauthorized(w, "User not found in context")
		return
	}

	// Call service to logout (blacklist token)
	if err := c.userService.Logout(r.Context(), token, userID); err != nil {
		utils.SendBadRequest(w, "Failed to logout: "+err.Error())
		return
	}

	utils.SendSuccess(w, "Logout successful", nil)
}
