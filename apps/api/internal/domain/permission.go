package domain

import "time"

type Permission struct {
	ID         string    `gorm:"primaryKey" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Permission string    `gorm:"type:varchar(255);uniqueIndex" json:"permission" example:"users:read"`
	Roles      []*Role   `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;" json:"roles,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
