package oauth

import (
	"time"
)

type AccessToken struct {
	ID        int       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Token     string    `json:"token" gorm:"unique;not null"`
	UserID    string    `json:"userId" gorm:"not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
}
