package repository

import (
	"testing"
	"time"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TokenBlacklistRepositoryTestSuite defines the test suite for TokenBlacklistRepository
type TokenBlacklistRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository ITokenBlacklistRepository
}

// SetupSuite runs once before all tests in the suite
func (suite *TokenBlacklistRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&domain.TokenBlacklist{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = TokenBlacklistRepositoryProvider(db)
}

// SetupTest runs before each test
func (suite *TokenBlacklistRepositoryTestSuite) SetupTest() {
	// Clean up the database before each test
	suite.db.Exec("DELETE FROM token_blacklists")
}

// TearDownSuite runs once after all tests in the suite
func (suite *TokenBlacklistRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Helper function to create a test blacklisted token
func (suite *TokenBlacklistRepositoryTestSuite) createTestBlacklistedToken(id, token, userID string, expiresAt int64) *domain.TokenBlacklist {
	blacklistedToken := &domain.TokenBlacklist{
		ID:        id,
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Reason:    "logout",
	}
	suite.db.Create(blacklistedToken)
	return blacklistedToken
}

// TestCreate tests the Create method
func (suite *TokenBlacklistRepositoryTestSuite) TestCreate() {
	// Test data
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := &domain.TokenBlacklist{
		ID:        "token-1",
		Token:     "jwt-token-string-123",
		UserID:    "user-1",
		ExpiresAt: expiresAt,
		Reason:    "logout",
	}

	// Execute
	err := suite.repository.Create(token)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the token was created
	var savedToken domain.TokenBlacklist
	suite.db.First(&savedToken, "id = ?", token.ID)
	assert.Equal(suite.T(), token.ID, savedToken.ID)
	assert.Equal(suite.T(), token.Token, savedToken.Token)
	assert.Equal(suite.T(), token.UserID, savedToken.UserID)
	assert.Equal(suite.T(), token.ExpiresAt, savedToken.ExpiresAt)
	assert.Equal(suite.T(), token.Reason, savedToken.Reason)
}

// TestCreate_DuplicateToken tests creating a token with duplicate token string
func (suite *TokenBlacklistRepositoryTestSuite) TestCreate_DuplicateToken() {
	// Create first token
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token1 := suite.createTestBlacklistedToken("token-1", "jwt-token-123", "user-1", expiresAt)

	// Try to create token with same token string
	token2 := &domain.TokenBlacklist{
		ID:        "token-2",
		Token:     token1.Token, // Same token string
		UserID:    "user-2",
		ExpiresAt: expiresAt,
		Reason:    "logout",
	}

	// Execute
	err := suite.repository.Create(token2)

	// Assert - should return error due to unique constraint
	assert.Error(suite.T(), err)
}

// TestIsBlacklisted tests the IsBlacklisted method
func (suite *TokenBlacklistRepositoryTestSuite) TestIsBlacklisted() {
	// Create test blacklisted token
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := suite.createTestBlacklistedToken("token-1", "jwt-token-123", "user-1", expiresAt)

	// Execute
	isBlacklisted, err := suite.repository.IsBlacklisted(token.Token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isBlacklisted)
}

// TestIsBlacklisted_NotFound tests IsBlacklisted with non-blacklisted token
func (suite *TokenBlacklistRepositoryTestSuite) TestIsBlacklisted_NotFound() {
	// Execute
	isBlacklisted, err := suite.repository.IsBlacklisted("non-existent-token")

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isBlacklisted)
}

// TestGetByToken tests the GetByToken method
func (suite *TokenBlacklistRepositoryTestSuite) TestGetByToken() {
	// Create test blacklisted token
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	expectedToken := suite.createTestBlacklistedToken("token-1", "jwt-token-123", "user-1", expiresAt)

	// Execute
	token, err := suite.repository.GetByToken(expectedToken.Token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), token)
	assert.Equal(suite.T(), expectedToken.ID, token.ID)
	assert.Equal(suite.T(), expectedToken.Token, token.Token)
	assert.Equal(suite.T(), expectedToken.UserID, token.UserID)
	assert.Equal(suite.T(), expectedToken.ExpiresAt, token.ExpiresAt)
}

// TestGetByToken_NotFound tests GetByToken with non-existent token
func (suite *TokenBlacklistRepositoryTestSuite) TestGetByToken_NotFound() {
	// Execute
	token, err := suite.repository.GetByToken("non-existent-token")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), token)
	assert.Equal(suite.T(), "token not found in blacklist", err.Error())
}

// TestDeleteExpired tests the DeleteExpired method
func (suite *TokenBlacklistRepositoryTestSuite) TestDeleteExpired() {
	// Create tokens with different expiration times
	pastTime := time.Now().Add(-1 * time.Hour).Unix()
	futureTime := time.Now().Add(24 * time.Hour).Unix()

	// Create expired token
	suite.createTestBlacklistedToken("token-1", "expired-token-1", "user-1", pastTime)
	suite.createTestBlacklistedToken("token-2", "expired-token-2", "user-2", pastTime)

	// Create valid token
	suite.createTestBlacklistedToken("token-3", "valid-token", "user-3", futureTime)

	// Execute
	err := suite.repository.DeleteExpired()

	// Assert
	assert.NoError(suite.T(), err)

	// Verify expired tokens are deleted and valid token remains
	var count int64
	suite.db.Model(&domain.TokenBlacklist{}).Count(&count)
	assert.Equal(suite.T(), int64(1), count)

	// Verify the remaining token is the valid one
	var remainingToken domain.TokenBlacklist
	suite.db.First(&remainingToken)
	assert.Equal(suite.T(), "valid-token", remainingToken.Token)
}

// TestDeleteExpired_NoExpiredTokens tests DeleteExpired when no tokens are expired
func (suite *TokenBlacklistRepositoryTestSuite) TestDeleteExpired_NoExpiredTokens() {
	// Create only valid tokens
	futureTime := time.Now().Add(24 * time.Hour).Unix()
	suite.createTestBlacklistedToken("token-1", "valid-token-1", "user-1", futureTime)
	suite.createTestBlacklistedToken("token-2", "valid-token-2", "user-2", futureTime)

	// Execute
	err := suite.repository.DeleteExpired()

	// Assert
	assert.NoError(suite.T(), err)

	// Verify all tokens remain
	var count int64
	suite.db.Model(&domain.TokenBlacklist{}).Count(&count)
	assert.Equal(suite.T(), int64(2), count)
}

// TestDeleteExpired_EmptyTable tests DeleteExpired on empty table
func (suite *TokenBlacklistRepositoryTestSuite) TestDeleteExpired_EmptyTable() {
	// Execute
	err := suite.repository.DeleteExpired()

	// Assert
	assert.NoError(suite.T(), err)

	// Verify table is still empty
	var count int64
	suite.db.Model(&domain.TokenBlacklist{}).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

// TestMultipleTokensForSameUser tests blacklisting multiple tokens for same user
func (suite *TokenBlacklistRepositoryTestSuite) TestMultipleTokensForSameUser() {
	// Create multiple tokens for same user
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	suite.createTestBlacklistedToken("token-1", "jwt-token-1", "user-1", expiresAt)
	suite.createTestBlacklistedToken("token-2", "jwt-token-2", "user-1", expiresAt)
	suite.createTestBlacklistedToken("token-3", "jwt-token-3", "user-1", expiresAt)

	// Execute - Check all tokens are blacklisted
	isBlacklisted1, _ := suite.repository.IsBlacklisted("jwt-token-1")
	isBlacklisted2, _ := suite.repository.IsBlacklisted("jwt-token-2")
	isBlacklisted3, _ := suite.repository.IsBlacklisted("jwt-token-3")

	// Assert
	assert.True(suite.T(), isBlacklisted1)
	assert.True(suite.T(), isBlacklisted2)
	assert.True(suite.T(), isBlacklisted3)

	// Verify count
	var count int64
	suite.db.Model(&domain.TokenBlacklist{}).Where("user_id = ?", "user-1").Count(&count)
	assert.Equal(suite.T(), int64(3), count)
}

// TestSuite entry point
func TestTokenBlacklistRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TokenBlacklistRepositoryTestSuite))
}
