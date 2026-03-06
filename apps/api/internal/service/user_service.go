package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	request "github.com/fahmihidayah/go-api-orchestrator/internal/data/request"
	"github.com/fahmihidayah/go-api-orchestrator/internal/data/response"
	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/mail"
	"github.com/fahmihidayah/go-api-orchestrator/internal/repository"
	"github.com/fahmihidayah/go-api-orchestrator/internal/security"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/go-playground/validator/v10"
)

// IUserService defines the interface for user service operations
type IUserService interface {
	Register(ctx context.Context, user *request.CreateUserRequest) (*response.UserResponse, error)
	Create(ctx context.Context, user *request.CreateUserRequest) error
	Login(ctx context.Context, req *request.LoginUserRequest) (*response.UserResponse, error)
	Logout(ctx context.Context, token string, userID string) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, req *request.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	GetAll(ctx context.Context, filter *request.FilterUserRequest) ([]domain.User, int64, error)
	Verify(ctx context.Context, req *request.VerifyUserRequest) error
	HandleFailedLogin(ctx context.Context, userID string) error
	HandleSuccessfulLogin(ctx context.Context, userID string) error
	InitiatePasswordReset(ctx context.Context, req *request.InitialResetPasswordUser) (string, error)
	CompletePasswordReset(ctx context.Context, req *request.ResetPasswordUser) error
	IsAccountLocked(user *domain.User) bool
	// React Admin specific methods
	GetAllReactAdmin(ctx context.Context, limit, offset int, sortField, sortOrder string, filters map[string]interface{}) ([]domain.User, int64, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.User, error)
	GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]domain.User, *utils.PaginateInfo, error)
	ChangePassword(ctx context.Context, req *request.ChangePasswordRequest) error
}

// UserServiceImpl implements IUserService
type UserServiceImpl struct {
	userRepository           repository.IUserRepository
	tokenBlacklistRepository repository.ITokenBlacklistRepository
	mailer                   *mail.Mailer
	validate                 *validator.Validate
	config                   config.Config
}

// UserServiceProvider creates a new instance of UserServiceImpl
func UserServiceProvider(userRepo repository.IUserRepository, tokenBlacklistRepo repository.ITokenBlacklistRepository, mailer *mail.Mailer, config *config.Config) IUserService {
	return &UserServiceImpl{
		userRepository:           userRepo,
		tokenBlacklistRepository: tokenBlacklistRepo,
		mailer:                   mailer,
		validate:                 validator.New(),
		config:                   *config,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, user *request.CreateUserRequest) (*response.UserResponse, error) {
	// Validate request
	if err := s.validate.Struct(user); err != nil {
		return nil, err
	}

	// Sanitize email
	user.Email = utils.SanitizeEmail(user.Email)

	// Validate email uniqueness
	existingUser, err := s.userRepository.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Hash password
	password, err := security.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Generate verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	// Create new user
	newUser := &domain.User{
		ID:                      utils.GenerateUUID(),
		Name:                    user.Name,
		Email:                   user.Email,
		HashedPassword:          password,
		IsVerified:              false,
		VerificationCode:        verificationToken,
		VerificationTokenExpiry: utils.AddDuration(24 * time.Hour),
		VerificationHash:        verificationToken,
		VerificationKind:        "email",

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user in database first
	if err := s.userRepository.Create(ctx, newUser); err != nil {
		return nil, err
	}
	s.mailer.SendVerifyEmail(newUser)
	// Then send verification email
	// if err := s.mailer.SendVerifyEmail(newUser); err != nil {
	// 	// User is created but email failed - this is acceptable
	// 	// User can request a new verification email
	// 	return err
	// }

	userWithToken, err := s.Login(ctx, &request.LoginUserRequest{
		Email:    user.Email,
		Password: user.Password,
	})

	return userWithToken, err
}

// Create creates a new user
func (s *UserServiceImpl) Create(ctx context.Context, user *request.CreateUserRequest) error {
	// Validate request
	if err := s.validate.Struct(user); err != nil {
		return err
	}

	// Sanitize email
	user.Email = utils.SanitizeEmail(user.Email)

	// Validate email uniqueness
	existingUser, err := s.userRepository.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already in use")
	}

	// Hash password
	password, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Create new user (admin created, so verified by default)
	newUser := &domain.User{
		ID:             utils.GenerateUUID(),
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: password,
		IsVerified:     true, // Admin-created users are auto-verified
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.userRepository.Create(ctx, newUser)
}

// Login authenticates a user and returns the user if successful
func (s *UserServiceImpl) Login(ctx context.Context, req *request.LoginUserRequest) (*response.UserResponse, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Sanitize email
	email := utils.SanitizeEmail(req.Email)

	// Get user by email
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if account is locked
	if s.IsAccountLocked(user) {
		timeUntil := utils.TimeUntilExpiry(user.LockUntil)
		return nil, errors.New("account is locked due to too many failed login attempts. Try again in " + timeUntil.String())
	}

	// Verify password
	if !security.VerifyPassword(user.HashedPassword, req.Password) {
		// Handle failed login
		_ = s.HandleFailedLogin(ctx, user.ID)
		return nil, errors.New("invalid email or password")
	}

	// Handle successful login (reset login attempts)
	_ = s.HandleSuccessfulLogin(ctx, user.ID)

	// Generate JWT token
	token, err := security.GenerateJWT(user.ID, user.Name, user.Email, s.config.JWTSecret, s.config.JWTExpirationHour)

	return &response.UserResponse{
		User:  *user,
		Token: token,
		Exp:   time.Now().Add(time.Duration(s.config.JWTExpirationHour) * time.Hour).Unix(),
	}, nil
}

// Logout invalidates the user's JWT token by adding it to the blacklist
func (s *UserServiceImpl) Logout(ctx context.Context, token string, userID string) error {
	// Validate the token first to get expiration time
	claims, err := security.ValidateJWT(token, s.config.JWTSecret)
	if err != nil && err != security.ErrExpiredToken {
		return errors.New("invalid token")
	}

	// If token is already expired, no need to blacklist
	if err == security.ErrExpiredToken {
		return errors.New("token has already expired")
	}

	// Create blacklist entry
	tokenBlacklist := &domain.TokenBlacklist{
		ID:        utils.GenerateUUID(),
		Token:     token,
		UserID:    userID,
		ExpiresAt: claims.ExpiresAt.Unix(),
		Reason:    "logout",
	}

	// Add token to blacklist
	if err := s.tokenBlacklistRepository.Create(tokenBlacklist); err != nil {
		return errors.New("failed to logout: " + err.Error())
	}

	return nil
}

// GetByID retrieves a user by ID
func (s *UserServiceImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

// Update updates user information
func (s *UserServiceImpl) Update(ctx context.Context, req *request.UpdateUserRequest) (*domain.User, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Get existing user
	user, err := s.userRepository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.Email != "" {
		// Sanitize and validate email uniqueness if email is being changed
		sanitizedEmail := utils.SanitizeEmail(req.Email)
		if sanitizedEmail != user.Email {
			existingUser, err := s.userRepository.GetByEmail(ctx, sanitizedEmail)
			if err == nil && existingUser != nil {
				return nil, errors.New("email already in use")
			}
		}
		user.Email = sanitizedEmail
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes a user
func (s *UserServiceImpl) Delete(ctx context.Context, id string) error {
	return s.userRepository.Delete(ctx, id)
}

// DeleteAll deletes multiple users by their IDs
func (s *UserServiceImpl) DeleteAll(ctx context.Context, ids []string) error {
	// Check if IDs array is empty
	if len(ids) == 0 {
		return errors.New("IDs array cannot be empty")
	}

	// Validate that all IDs are non-empty
	for _, id := range ids {
		if id == "" {
			return errors.New("invalid ID: empty ID in array")
		}
	}

	return s.userRepository.DeleteAll(ctx, ids)
}

// GetAll retrieves a paginated list of users with total count
func (s *UserServiceImpl) GetAll(ctx context.Context, filter *request.FilterUserRequest) ([]domain.User, int64, error) {
	// Calculate offset from page
	offset := (filter.Page - 1) * filter.Limit

	// Get filtered users from repository
	users, err := s.userRepository.GetAllWithFilter(ctx, filter.Name, filter.Email, filter.Limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count with filter
	count, err := s.userRepository.CountWithFilter(ctx, filter.Name, filter.Email)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// Verify verifies a user by verification code
func (s *UserServiceImpl) Verify(ctx context.Context, req *request.VerifyUserRequest) error {
	user, err := s.userRepository.FindByEmailAndVerificationCode(ctx, req.Email, req.VerificationCode)
	if err != nil {
		return err
	}

	// Check if verification token is expired
	if user.VerificationTokenExpiry > 0 && time.Now().Unix() > user.VerificationTokenExpiry {
		return errors.New("verification code expired")
	}

	// Check if already verified
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// Set user as verified and clear verification data
	user.IsVerified = true
	user.VerificationCode = ""
	user.VerificationHash = ""
	user.VerificationTokenExpiry = 0

	return s.userRepository.Update(ctx, user)
}

// HandleFailedLogin increments login attempts and locks account if necessary
func (s *UserServiceImpl) HandleFailedLogin(ctx context.Context, userID string) error {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	attempts := user.LoginAttempts + 1
	var lockUntil int64 = 0

	// Lock account after 5 failed attempts for 30 minutes
	if attempts >= 5 {
		lockUntil = time.Now().Add(30 * time.Minute).Unix()
	}

	return s.userRepository.UpdateLoginAttempts(ctx, userID, attempts, lockUntil)
}

// HandleSuccessfulLogin resets login attempts
func (s *UserServiceImpl) HandleSuccessfulLogin(ctx context.Context, userID string) error {
	return s.userRepository.UpdateLoginAttempts(ctx, userID, 0, 0)
}

// InitiatePasswordReset generates a reset token and returns it
func (s *UserServiceImpl) InitiatePasswordReset(ctx context.Context, req *request.InitialResetPasswordUser) (string, error) {
	if err := s.validate.Struct(req); err != nil {
		return "", err
	}

	email := req.Email

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists
		return "", errors.New("if the email exists, a reset link will be sent")
	}

	// Generate cryptographically secure reset token
	resetToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	user.ResetPasswordToken = resetToken
	user.ResetPasswordTokenExpiry = utils.AddDuration(1 * time.Hour)

	if err := s.userRepository.Update(ctx, user); err != nil {
		return "", err
	}

	if err := s.mailer.SendResetPassword(user); err != nil {
		return "", err
	}

	return resetToken, nil
}

// CompletePasswordReset validates token and updates password
func (s *UserServiceImpl) CompletePasswordReset(ctx context.Context, req *request.ResetPasswordUser) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}
	existingUser, err := s.userRepository.FindByResetPasswordToken(ctx, req.Token)
	if err != nil || existingUser == nil {
		return errors.New("invalid or expired reset token")
	}

	// Check if token is expired
	if existingUser.ResetPasswordTokenExpiry > 0 && time.Now().Unix() > existingUser.ResetPasswordTokenExpiry {
		return errors.New("reset token has expired")
	}

	password, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	existingUser.HashedPassword = password

	// Clear reset token after use
	existingUser.ResetPasswordToken = ""
	existingUser.ResetPasswordTokenExpiry = 0

	return s.userRepository.Update(ctx, existingUser)
}

// IsAccountLocked checks if a user account is currently locked
func (s *UserServiceImpl) IsAccountLocked(user *domain.User) bool {
	if user.LockUntil == 0 {
		return false
	}
	return !utils.IsExpired(user.LockUntil)
}

// GetByEmail retrieves a user by email
func (s *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	sanitizedEmail := utils.SanitizeEmail(email)
	return s.userRepository.GetByEmail(ctx, sanitizedEmail)
}

// GetAllReactAdmin retrieves users with React Admin parameters
func (s *UserServiceImpl) GetAllReactAdmin(ctx context.Context, limit, offset int, sortField, sortOrder string, filters map[string]interface{}) ([]domain.User, int64, error) {
	// Extract name and email filters if present
	name := ""
	email := ""
	if nameVal, ok := filters["name"].(string); ok {
		name = nameVal
	}
	if emailVal, ok := filters["email"].(string); ok {
		email = emailVal
	}

	// Get filtered users from repository
	users, err := s.userRepository.GetAllWithFilter(ctx, name, email, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count with filter
	count, err := s.userRepository.CountWithFilter(ctx, name, email)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// GetByIDs retrieves multiple users by their IDs
func (s *UserServiceImpl) GetByIDs(ctx context.Context, ids []string) ([]domain.User, error) {
	if len(ids) == 0 {
		return []domain.User{}, nil
	}

	users := make([]domain.User, 0, len(ids))
	for _, id := range ids {
		user, err := s.userRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip users that don't exist
		}
		users = append(users, *user)
	}

	return users, nil
}

// GetWithQueryParams retrieves users with React Admin parameters
func (s *UserServiceImpl) GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]domain.User, *utils.PaginateInfo, error) {
	count, err := s.userRepository.CountByQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	paginateInfo := queryParams.ToPaginateInfo(count)
	users, err := s.userRepository.GetWithQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	return users, paginateInfo, nil
}

func (s *UserServiceImpl) ChangePassword(ctx context.Context, req *request.ChangePasswordRequest) error {
	log.Printf("[UserService.ChangePassword] Starting password change for user ID: %s", req.ID)

	users, err := s.userRepository.GetByID(ctx, req.ID)
	if err != nil {
		log.Printf("[UserService.ChangePassword] User not found for ID: %s - Error: %v", req.ID, err)
		return errors.New("user not found")
	}

	log.Printf("[UserService.ChangePassword] User found: %s, verifying old password", users.Email)

	// Verify old password
	if !security.VerifyPassword(users.HashedPassword, req.OldPassword) {
		log.Printf("[UserService.ChangePassword] Old password verification failed for user: %s", users.Email)
		log.Printf("[UserService.ChangePassword] Request has OldPassword length: %d", len(req.OldPassword))
		return errors.New("old password is incorrect")
	}

	log.Printf("[UserService.ChangePassword] Old password verified successfully for user: %s", users.Email)

	// Hash new password
	newHashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("[UserService.ChangePassword] Failed to hash new password for user: %s - Error: %v", users.Email, err)
		return err
	}

	log.Printf("[UserService.ChangePassword] New password hashed successfully, updating user: %s", users.Email)

	// Update user's password
	users.HashedPassword = newHashedPassword
	users.UpdatedAt = time.Now()

	err = s.userRepository.Update(ctx, users)
	if err != nil {
		log.Printf("[UserService.ChangePassword] Failed to update user password in database: %s - Error: %v", users.Email, err)
		return err
	}

	log.Printf("[UserService.ChangePassword] Password changed successfully for user: %s", users.Email)
	return nil
}
