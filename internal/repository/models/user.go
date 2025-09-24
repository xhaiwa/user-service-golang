package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"not null" json:"-"` // never expose the hash
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // soft delete
	OAuthID       string         `gorm:"uniqueIndex;null" json:"oauth_id,omitempty"`
	OAuthProvider string         `json:"oauth_provider,omitempty"` // "google", "facebook", etc.
}
