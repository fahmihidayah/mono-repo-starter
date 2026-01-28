package domain

import "time"

type Account struct {
	ID                  string    `gorm:"primaryKey" json:"id"`
	UserID              string    `gorm:"type:varchar(255);index;not null" json:"user_id"` // Foreign key to User
	Name                string    `gorm:"type:varchar(255)" json:"name"`
	Picture             string    `gorm:"type:varchar(500)" json:"picture,omitempty"`
	IssuerName          string    `gorm:"type:varchar(100)" json:"issuer_name,omitempty"`
	Scope               string    `gorm:"type:text" json:"scope,omitempty"`
	Sub                 string    `gorm:"type:varchar(255)" json:"sub,omitempty"`
	AccessToken         string    `gorm:"type:text" json:"-"`
	PasskeyCredentialID string    `gorm:"type:varchar(255)" json:"passkey_credential_id,omitempty"`
	PasskeyPublicKey    string    `gorm:"type:text" json:"passkey_public_key,omitempty"`
	PasskeyCounter      uint32    `gorm:"type:int;default:0" json:"passkey_counter,omitempty"`
	PasskeyTransports   string    `gorm:"type:varchar(255)" json:"passkey_transports,omitempty"`
	PasskeyDeviceType   string    `gorm:"type:varchar(50)" json:"passkey_device_type,omitempty"`
	PasskeyBackedUp     bool      `gorm:"type:boolean;default:false" json:"passkey_backed_up"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
