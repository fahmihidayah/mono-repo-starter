package repository

import (
	"context"
	"testing"
	"time"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserRepositoryTestSuite defines the test suite for UserRepository
type UserRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository IUserRepository
}

// SetupSuite runs once before all tests in the suite
func (suite *UserRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&domain.User{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = UserRepositoryProvider(db)
}

// SetupTest runs before each test
func (suite *UserRepositoryTestSuite) SetupTest() {
	// Clean up the database before each test
	suite.db.Exec("DELETE FROM users")
}

// TearDownSuite runs once after all tests in the suite
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Helper function to create a test user
func (suite *UserRepositoryTestSuite) createTestUser(id, email, name string) *domain.User {
	user := &domain.User{
		ID:             id,
		Email:          email,
		Name:           name,
		HashedPassword: "hashed_password_123",
		IsVerified:     true,
	}
	suite.db.Create(user)
	return user
}

// TestCreate tests the Create method
func (suite *UserRepositoryTestSuite) TestCreate() {
	// Test data
	ctx := context.Background()
	user := &domain.User{
		ID:             "user-1",
		Email:          "test@example.com",
		Name:           "Test User",
		HashedPassword: "hashed_password",
		IsVerified:     false,
	}

	// Execute
	err := suite.repository.Create(ctx, user)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the user was created
	var savedUser domain.User
	suite.db.First(&savedUser, "id = ?", user.ID)
	assert.Equal(suite.T(), user.ID, savedUser.ID)
	assert.Equal(suite.T(), user.Email, savedUser.Email)
	assert.Equal(suite.T(), user.Name, savedUser.Name)
	assert.Equal(suite.T(), user.HashedPassword, savedUser.HashedPassword)
	assert.Equal(suite.T(), user.IsVerified, savedUser.IsVerified)
}

// TestCreate_DuplicateEmail tests creating a user with duplicate email
func (suite *UserRepositoryTestSuite) TestCreate_DuplicateEmail() {
	// Create first user
	ctx := context.Background()
	user1 := suite.createTestUser("user-1", "test@example.com", "User 1")

	// Try to create user with same email
	user2 := &domain.User{
		ID:             "user-2",
		Email:          user1.Email, // Same email
		Name:           "User 2",
		HashedPassword: "hashed_password",
	}

	// Execute
	err := suite.repository.Create(ctx, user2)

	// Assert - should return error due to unique constraint
	assert.Error(suite.T(), err)
}

// TestGetByID tests the GetByID method
func (suite *UserRepositoryTestSuite) TestGetByID() {
	// Create test user
	ctx := context.Background()
	expectedUser := suite.createTestUser("user-1", "test@example.com", "Test User")

	// Execute
	user, err := suite.repository.GetByID(ctx, expectedUser.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), expectedUser.ID, user.ID)
	assert.Equal(suite.T(), expectedUser.Email, user.Email)
	assert.Equal(suite.T(), expectedUser.Name, user.Name)
}

// TestGetByID_NotFound tests getting a non-existent user
func (suite *UserRepositoryTestSuite) TestGetByID_NotFound() {
	// Execute
	ctx := context.Background()
	user, err := suite.repository.GetByID(ctx, "non-existent-id")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// TestGetByEmail tests the GetByEmail method
func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	// Create test user
	ctx := context.Background()
	expectedUser := suite.createTestUser("user-1", "test@example.com", "Test User")

	// Execute
	user, err := suite.repository.GetByEmail(ctx, expectedUser.Email)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), expectedUser.ID, user.ID)
	assert.Equal(suite.T(), expectedUser.Email, user.Email)
	assert.Equal(suite.T(), expectedUser.Name, user.Name)
}

// TestGetByEmail_NotFound tests getting a user with non-existent email
func (suite *UserRepositoryTestSuite) TestGetByEmail_NotFound() {
	// Execute
	ctx := context.Background()
	user, err := suite.repository.GetByEmail(ctx, "nonexistent@example.com")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// TestGetAll tests the GetAll method
func (suite *UserRepositoryTestSuite) TestGetAll() {
	// Create multiple test users
	ctx := context.Background()
	suite.createTestUser("user-1", "user1@example.com", "User 1")
	suite.createTestUser("user-2", "user2@example.com", "User 2")
	suite.createTestUser("user-3", "user3@example.com", "User 3")

	// Execute
	users, err := suite.repository.GetAll(ctx, 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 3)
}

// TestGetAll_WithPagination tests GetAll with limit and offset
func (suite *UserRepositoryTestSuite) TestGetAll_WithPagination() {
	// Create 5 test users
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.User{
			ID:    "user-" + string(rune(i)),
			Email: "user" + string(rune(i)) + "@example.com",
			Name:  "User " + string(rune(i)),
		})
	}

	// Execute - Get 2 users, skip first 2
	users, err := suite.repository.GetAll(ctx, 2, 2)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2)
}

// TestGetAll_Empty tests GetAll when no users exist
func (suite *UserRepositoryTestSuite) TestGetAll_Empty() {
	// Execute
	ctx := context.Background()
	users, err := suite.repository.GetAll(ctx, 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), users)
}

// TestGetAllWithFilter tests GetAllWithFilter with name filter
func (suite *UserRepositoryTestSuite) TestGetAllWithFilter_NameFilter() {
	// Create test users
	ctx := context.Background()
	suite.createTestUser("user-1", "john@example.com", "John Doe")
	suite.createTestUser("user-2", "jane@example.com", "Jane Smith")
	suite.createTestUser("user-3", "johnny@example.com", "Johnny Depp")

	// Execute - Filter by name containing "john"
	users, err := suite.repository.GetAllWithFilter(ctx, "john", "", 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2) // John Doe and Johnny Depp
}

// TestGetAllWithFilter tests GetAllWithFilter with email filter
func (suite *UserRepositoryTestSuite) TestGetAllWithFilter_EmailFilter() {
	// Create test users
	ctx := context.Background()
	suite.createTestUser("user-1", "john@example.com", "John Doe")
	suite.createTestUser("user-2", "jane@gmail.com", "Jane Smith")
	suite.createTestUser("user-3", "johnny@example.com", "Johnny Depp")

	// Execute - Filter by email containing "example"
	users, err := suite.repository.GetAllWithFilter(ctx, "", "example", 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2) // john@example.com and johnny@example.com
}

// TestGetAllWithFilter tests GetAllWithFilter with both filters
func (suite *UserRepositoryTestSuite) TestGetAllWithFilter_BothFilters() {
	// Create test users
	ctx := context.Background()
	suite.createTestUser("user-1", "john@example.com", "John Doe")
	suite.createTestUser("user-2", "jane@example.com", "Jane Smith")
	suite.createTestUser("user-3", "johnny@gmail.com", "Johnny Depp")

	// Execute - Filter by name "john" AND email "example"
	users, err := suite.repository.GetAllWithFilter(ctx, "john", "example", 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1) // Only john@example.com
	assert.Equal(suite.T(), "john@example.com", users[0].Email)
}

// TestGetAllWithFilter_NoFilters tests GetAllWithFilter without filters
func (suite *UserRepositoryTestSuite) TestGetAllWithFilter_NoFilters() {
	// Create test users
	ctx := context.Background()
	suite.createTestUser("user-1", "user1@example.com", "User 1")
	suite.createTestUser("user-2", "user2@example.com", "User 2")

	// Execute - No filters
	users, err := suite.repository.GetAllWithFilter(ctx, "", "", 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2) // All users
}

// TestUpdate tests the Update method
func (suite *UserRepositoryTestSuite) TestUpdate() {
	// Create test user
	ctx := context.Background()
	user := suite.createTestUser("user-1", "original@example.com", "Original Name")

	// Modify the user
	user.Name = "Updated Name"
	user.Email = "updated@example.com"
	user.IsVerified = true

	// Execute
	err := suite.repository.Update(ctx, user)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the update
	updatedUser, _ := suite.repository.GetByID(ctx, user.ID)
	assert.Equal(suite.T(), "Updated Name", updatedUser.Name)
	assert.Equal(suite.T(), "updated@example.com", updatedUser.Email)
	assert.True(suite.T(), updatedUser.IsVerified)
}

// TestDelete tests the Delete method
func (suite *UserRepositoryTestSuite) TestDelete() {
	// Create test user
	ctx := context.Background()
	user := suite.createTestUser("user-1", "test@example.com", "Test User")

	// Execute
	err := suite.repository.Delete(ctx, user.ID)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the user is deleted
	deletedUser, err := suite.repository.GetByID(ctx, user.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), deletedUser)
}

// TestDelete_NotFound tests deleting a non-existent user
func (suite *UserRepositoryTestSuite) TestDelete_NotFound() {
	// Execute
	ctx := context.Background()
	err := suite.repository.Delete(ctx, "non-existent-id")

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// TestDeleteAll tests the DeleteAll method
func (suite *UserRepositoryTestSuite) TestDeleteAll() {
	// Create multiple test users
	ctx := context.Background()
	user1 := suite.createTestUser("user-1", "user1@example.com", "User 1")
	user2 := suite.createTestUser("user-2", "user2@example.com", "User 2")
	user3 := suite.createTestUser("user-3", "user3@example.com", "User 3")

	// Execute - Delete first two users
	ids := []string{user1.ID, user2.ID}
	err := suite.repository.DeleteAll(ctx, ids)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify only user3 remains
	users, _ := suite.repository.GetAll(ctx, 10, 0)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), user3.ID, users[0].ID)
}

// TestDeleteAll_NoUsersFound tests DeleteAll with non-existent IDs
func (suite *UserRepositoryTestSuite) TestDeleteAll_NoUsersFound() {
	// Execute
	ctx := context.Background()
	ids := []string{"non-existent-1", "non-existent-2"}
	err := suite.repository.DeleteAll(ctx, ids)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "no users found with the provided IDs", err.Error())
}

// TestDeleteAll_PartialMatch tests DeleteAll with mix of existing and non-existing IDs
func (suite *UserRepositoryTestSuite) TestDeleteAll_PartialMatch() {
	// Create test user
	ctx := context.Background()
	user := suite.createTestUser("user-1", "user1@example.com", "User 1")

	// Execute - Mix of existing and non-existing IDs
	ids := []string{user.ID, "non-existent-id"}
	err := suite.repository.DeleteAll(ctx, ids)

	// Assert - Should succeed because at least one ID exists
	assert.NoError(suite.T(), err)

	// Verify the existing user was deleted
	users, _ := suite.repository.GetAll(ctx, 10, 0)
	assert.Len(suite.T(), users, 0)
}

// TestCount tests the Count method
func (suite *UserRepositoryTestSuite) TestCount() {
	// Create multiple test users
	ctx := context.Background()
	suite.createTestUser("user-1", "user1@example.com", "User 1")
	suite.createTestUser("user-2", "user2@example.com", "User 2")
	suite.createTestUser("user-3", "user3@example.com", "User 3")

	// Execute
	count, err := suite.repository.Count(ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count)
}

// TestCount_Empty tests Count when no users exist
func (suite *UserRepositoryTestSuite) TestCount_Empty() {
	// Execute
	ctx := context.Background()
	count, err := suite.repository.Count(ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}

// TestCountWithFilter tests CountWithFilter with name filter
func (suite *UserRepositoryTestSuite) TestCountWithFilter_NameFilter() {
	// Create test users
	ctx := context.Background()
	suite.createTestUser("user-1", "john@example.com", "John Doe")
	suite.createTestUser("user-2", "jane@example.com", "Jane Smith")
	suite.createTestUser("user-3", "johnny@example.com", "Johnny Depp")

	// Execute - Count users with name containing "john"
	count, err := suite.repository.CountWithFilter(ctx, "john", "")

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count) // John Doe and Johnny Depp
}

// TestCountWithFilter tests CountWithFilter with email filter
func (suite *UserRepositoryTestSuite) TestCountWithFilter_EmailFilter() {
	// Create test users
	ctx := context.Background()
	suite.createTestUser("user-1", "john@example.com", "John Doe")
	suite.createTestUser("user-2", "jane@gmail.com", "Jane Smith")
	suite.createTestUser("user-3", "johnny@example.com", "Johnny Depp")

	// Execute - Count users with email containing "example"
	count, err := suite.repository.CountWithFilter(ctx, "", "example")

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count) // john@example.com and johnny@example.com
}

// TestCountWithFilter tests CountWithFilter with both filters
func (suite *UserRepositoryTestSuite) TestCountWithFilter_BothFilters() {
	// Create test users
	ctx := context.Background()
	suite.createTestUser("user-1", "john@example.com", "John Doe")
	suite.createTestUser("user-2", "jane@example.com", "Jane Smith")
	suite.createTestUser("user-3", "johnny@gmail.com", "Johnny Depp")

	// Execute - Count users with name "john" AND email "example"
	count, err := suite.repository.CountWithFilter(ctx, "john", "example")

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count) // Only john@example.com
}

// TestFindByResetPasswordToken tests FindByResetPasswordToken method
func (suite *UserRepositoryTestSuite) TestFindByResetPasswordToken() {
	// Create test user with reset token
	ctx := context.Background()
	token := "reset-token-123"
	user := &domain.User{
		ID:                       "user-1",
		Email:                    "test@example.com",
		Name:                     "Test User",
		ResetPasswordToken:       token,
		ResetPasswordTokenExpiry: time.Now().Add(1 * time.Hour).Unix(),
	}
	suite.db.Create(user)

	// Execute
	foundUser, err := suite.repository.FindByResetPasswordToken(ctx, token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), user.ID, foundUser.ID)
	assert.Equal(suite.T(), token, foundUser.ResetPasswordToken)
}

// TestFindByResetPasswordToken_NotFound tests finding user with non-existent token
func (suite *UserRepositoryTestSuite) TestFindByResetPasswordToken_NotFound() {
	// Execute
	ctx := context.Background()
	user, err := suite.repository.FindByResetPasswordToken(ctx, "non-existent-token")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// TestFindByEmailAndVerificationCode tests FindByEmailAndVerificationCode method
func (suite *UserRepositoryTestSuite) TestFindByEmailAndVerificationCode() {
	// Create test user with verification code
	ctx := context.Background()
	email := "test@example.com"
	code := "123456"
	user := &domain.User{
		ID:               "user-1",
		Email:            email,
		Name:             "Test User",
		VerificationCode: code,
		IsVerified:       false,
	}
	suite.db.Create(user)

	// Execute
	foundUser, err := suite.repository.FindByEmailAndVerificationCode(ctx, email, code)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), user.ID, foundUser.ID)
	assert.Equal(suite.T(), email, foundUser.Email)
	assert.Equal(suite.T(), code, foundUser.VerificationCode)
}

// TestFindByEmailAndVerificationCode_NotFound tests finding with wrong code
func (suite *UserRepositoryTestSuite) TestFindByEmailAndVerificationCode_NotFound() {
	// Create test user
	ctx := context.Background()
	suite.createTestUser("user-1", "test@example.com", "Test User")

	// Execute - Wrong verification code
	user, err := suite.repository.FindByEmailAndVerificationCode(ctx, "test@example.com", "wrong-code")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// TestUpdateLoginAttempts tests UpdateLoginAttempts method
func (suite *UserRepositoryTestSuite) TestUpdateLoginAttempts() {
	// Create test user
	ctx := context.Background()
	user := suite.createTestUser("user-1", "test@example.com", "Test User")

	// Execute - Update login attempts
	lockUntil := time.Now().Add(30 * time.Minute).Unix()
	err := suite.repository.UpdateLoginAttempts(ctx, user.ID, 3, lockUntil)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the update
	updatedUser, _ := suite.repository.GetByID(ctx, user.ID)
	assert.Equal(suite.T(), 3, updatedUser.LoginAttempts)
	assert.Equal(suite.T(), lockUntil, updatedUser.LockUntil)
}

// TestUpdateLoginAttempts_NotFound tests updating login attempts for non-existent user
func (suite *UserRepositoryTestSuite) TestUpdateLoginAttempts_NotFound() {
	// Execute
	ctx := context.Background()
	err := suite.repository.UpdateLoginAttempts(ctx, "non-existent-id", 3, 0)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// TestResetPassword tests ResetPassword method
func (suite *UserRepositoryTestSuite) TestResetPassword() {
	// Create test user with reset token
	ctx := context.Background()
	user := &domain.User{
		ID:                       "user-1",
		Email:                    "test@example.com",
		Name:                     "Test User",
		HashedPassword:           "old_hashed_password",
		ResetPasswordToken:       "reset-token-123",
		ResetPasswordTokenExpiry: time.Now().Add(1 * time.Hour).Unix(),
	}
	suite.db.Create(user)

	// Execute - Reset password
	newHashedPassword := "new_hashed_password"
	err := suite.repository.ResetPassword(ctx, user.ID, newHashedPassword)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the update
	updatedUser, _ := suite.repository.GetByID(ctx, user.ID)
	assert.Equal(suite.T(), newHashedPassword, updatedUser.HashedPassword)
	assert.Equal(suite.T(), "", updatedUser.ResetPasswordToken)
	assert.Equal(suite.T(), int64(0), updatedUser.ResetPasswordTokenExpiry)
}

// TestResetPassword_NotFound tests resetting password for non-existent user
func (suite *UserRepositoryTestSuite) TestResetPassword_NotFound() {
	// Execute
	ctx := context.Background()
	err := suite.repository.ResetPassword(ctx, "non-existent-id", "new_password")

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// TestSuite entry point
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestGetAllByQueryParams_Success() {
	ctx := context.Background()
	// Create a test user

	suite.createTestUser("user-1", "user1@example.com", "User 1")
	suite.createTestUser("user-2", "user2@example.com", "User 2")
	suite.createTestUser("user-3", "user3@example.com", "User 3")

	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"name", "asc"},
		Filter: map[string]interface{}{
			"ids": []string{"user-1", "user-2"},
		},
	}

	users, err := suite.repository.GetWithQuery(ctx, queryParameter)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2)
}

func (suite *UserRepositoryTestSuite) TestCountByQueryParams_Success() {
	ctx := context.Background()
	// Create a test user

	suite.createTestUser("user-1", "user1@example.com", "User 1")
	suite.createTestUser("user-2", "user2@example.com", "User 2")
	suite.createTestUser("user-3", "user3@example.com", "User 3")

	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"name", "asc"},
		Filter: map[string]interface{}{},
	}

	count, err := suite.repository.CountByQuery(ctx, queryParameter)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, int64(3))
}
