package domain

type Media struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Alt       string `gorm:"type:varchar(255)" json:"alt"`
	Url       string `gorm:"type:varchar(500)" json:"url"`
	Path      string `gorm:"type:varchar(500);uniqueIndex" json:"path"`
	FileName  string `gorm:"type:varchar(255)" json:"file_name"`
	MimeType  string `gorm:"type:varchar(100)" json:"mime_type"`
	FileSize  int64  `json:"file_size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
}
