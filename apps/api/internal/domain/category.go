package domain

type Category struct {
	ID        string  `gorm:"primaryKey" json:"id"`
	Slug      string  `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	Title     string  `gorm:"type:varchar(255)" json:"title"`
	Posts     []*Post `gorm:"many2many:post_categories;constraint:OnDelete:CASCADE;" json:"posts,omitempty"`
	CreatedAt int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64   `gorm:"autoUpdateTime" json:"updated_at"`
}
