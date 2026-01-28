package domain

import "time"

// User represents a user in the system
type User struct {
	ID                       string    `gorm:"primaryKey" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name                     string    `gorm:"type:varchar(255)" json:"name" example:"John Doe"`
	Email                    string    `gorm:"type:varchar(255);uniqueIndex" json:"email" example:"user@example.com"`
	ResetPasswordToken       string    `gorm:"type:varchar(255);column:reset_password_token" json:"-"`
	ResetPasswordTokenExpiry int64     `gorm:"type:bigint" json:"-"`
	HashedPassword           string    `gorm:"type:text" json:"-"`              // Consolidated password hash
	HashSalt                 string    `gorm:"type:varchar(255)" json:"-"`      // Consolidated salt
	HashIterations           int       `gorm:"type:int;default:10000" json:"-"` // Hash iteration count for PBKDF2 or similar
	LoginAttempts            int       `gorm:"type:int;default:0" json:"-"`
	LockUntil                int64     `gorm:"type:bigint;default:0" json:"-"`
	IsSuperUser              bool      `gorm:"type:boolean;default:false" json:"is_super_user" example:"false"`
	IsVerified               bool      `gorm:"type:boolean;default:false" json:"is_verified" example:"false"`
	VerificationCode         string    `gorm:"type:varchar(255)" json:"-"`
	VerificationHash         string    `gorm:"type:varchar(255)" json:"-"`
	VerificationTokenExpiry  int64     `gorm:"type:bigint" json:"-"`
	VerificationKind         string    `gorm:"type:varchar(50)" json:"-"`
	CreatedAt                time.Time `gorm:"autoCreateTime" json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt                time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
