package domain

import "time"


type Role struct {
	ID 					string `gorm:"primaryKey" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name 				string `gorm:"type:varchar(255)" json:"name" example:"Admin"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2024-01-01T00:00:00Z"`
}