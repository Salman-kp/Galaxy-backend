package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	BookingStatusBooked    = "booked"
	BookingStatusPresent   = "present"
	BookingStatusAbsent    = "absent"
	BookingStatusCompleted = "completed"
)

type Booking struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	
	EventID     uint           `gorm:"not null;index;uniqueIndex:idx_event_user" json:"event_id"`
	UserID      uint           `gorm:"not null;index;uniqueIndex:idx_event_user" json:"user_id"`
	Role        string         `gorm:"size:50;not null;index" json:"role"`
	
	Status      string         `gorm:"size:20;default:'booked';index" json:"status"`
	
	BaseAmount  int64          `gorm:"default:0" json:"base_amount"`
	ExtraAmount int64          `gorm:"default:0" json:"extra_amount"` // long work
	TAAmount    int64          `gorm:"default:0" json:"ta_amount"`
	BonusAmount int64          `gorm:"default:0" json:"bonus_amount"`
	FineAmount  int64          `gorm:"default:0" json:"fine_amount"`
	TotalAmount int64          `gorm:"default:0" json:"total_amount"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}