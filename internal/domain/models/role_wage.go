package models

import "time"

type RoleWage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Role      string    `gorm:"size:50;uniqueIndex;not null" json:"role"`
	Wage      int64     `gorm:"not null" json:"wage"`
	CreatedAt time.Time `json:"created_at"` // added
	UpdatedAt time.Time `json:"updated_at"` // added
}