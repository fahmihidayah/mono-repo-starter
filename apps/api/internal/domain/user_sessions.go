package domain

import "time"

type UserSession struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(255);index;not null" json:"user_id"` // Foreign key to User
	AccountID string    `gorm:"type:varchar(255);index" json:"account_id"`       // Foreign key to Account (which authentication method was used)
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
