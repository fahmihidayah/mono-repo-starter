package domain

import "time"

// TokenBlacklist represents a blacklisted JWT token
type TokenBlacklist struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"type:text;uniqueIndex;not null" json:"token"`       // The actual JWT token string
	UserID    string    `gorm:"type:varchar(255);index;not null" json:"user_id"`   // User who owns this token
	ExpiresAt int64     `gorm:"type:bigint;index;not null" json:"expires_at"`      // When the token expires (for cleanup)
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                  // When token was blacklisted
	Reason    string    `gorm:"type:varchar(255)" json:"reason,omitempty"`         // Optional: logout, compromised, etc.
}
