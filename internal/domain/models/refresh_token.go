package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	TokenHashed  string         `gorm:"not null;unique" json:"-"`
	ExpiresAt    time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}