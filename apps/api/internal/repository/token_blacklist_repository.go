package repository

import (
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"gorm.io/gorm"
)

// ITokenBlacklistRepository defines the interface for token blacklist repository operations
type ITokenBlacklistRepository interface {
	Create(tokenBlacklist *domain.TokenBlacklist) error
	IsBlacklisted(token string) (bool, error)
	DeleteExpired() error // Cleanup expired tokens
	GetByToken(token string) (*domain.TokenBlacklist, error)
}

// TokenBlacklistRepositoryImpl implements ITokenBlacklistRepository
type TokenBlacklistRepositoryImpl struct {
	db *gorm.DB
}

// TokenBlacklistRepositoryProvider creates a new instance of TokenBlacklistRepositoryImpl
func TokenBlacklistRepositoryProvider(db *gorm.DB) ITokenBlacklistRepository {
	return &TokenBlacklistRepositoryImpl{db: db}
}

// Create adds a token to the blacklist
func (r *TokenBlacklistRepositoryImpl) Create(tokenBlacklist *domain.TokenBlacklist) error {
	return r.db.Create(tokenBlacklist).Error
}

// IsBlacklisted checks if a token is in the blacklist
func (r *TokenBlacklistRepositoryImpl) IsBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.TokenBlacklist{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteExpired removes expired tokens from the blacklist (cleanup)
func (r *TokenBlacklistRepositoryImpl) DeleteExpired() error {
	// Delete tokens where expires_at is less than current Unix timestamp
	result := r.db.Where("expires_at < ?", gorm.Expr("strftime('%s', 'now')")).Delete(&domain.TokenBlacklist{})
	return result.Error
}

// GetByToken retrieves a blacklisted token by its token string
func (r *TokenBlacklistRepositoryImpl) GetByToken(token string) (*domain.TokenBlacklist, error) {
	var tokenBlacklist domain.TokenBlacklist
	if err := r.db.Where("token = ?", token).First(&tokenBlacklist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found in blacklist")
		}
		return nil, err
	}
	return &tokenBlacklist, nil
}
