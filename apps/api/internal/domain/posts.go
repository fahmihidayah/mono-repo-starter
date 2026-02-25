package domain

type Post struct {
	ID         string       `gorm:"primaryKey" json:"id"`
	Slug       string       `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	Title      string       `gorm:"type:varchar(255)" json:"title"`
	Content    string       `gorm:"type:text" json:"content"`
	Categories []*Category 	`gorm:"many2many:post_categories;constraint:OnDelete:CASCADE;" json:"categories"`
	UserID     string       `gorm:"type:varchar(255);index;not null" json:"user_id"` // Foreign key to User
	User       *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt  int64        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  int64        `gorm:"autoUpdateTime" json:"updated_at"`
}
