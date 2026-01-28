package domain

type Category struct {
	ID        string  `gorm:"primaryKey" json:"id"`
	Slug      string  `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	Name      string  `gorm:"type:varchar(255)" json:"name"`
	Posts     []*Post `gorm:"many2many:post_categories;" json:"posts,omitempty"`
	CreatedAt int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64   `gorm:"autoUpdateTime" json:"updated_at"`
}
