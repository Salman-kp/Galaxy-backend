package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	EventStatusUpcoming  = "upcoming"
	EventStatusOngoing   = "ongoing"
	EventStatusCompleted = "completed"
	EventStatusCancelled = "cancelled"

	TimeSlotMorning = "morning"
	TimeSlotLunch   = "lunch"
	TimeSlotNight   = "night"

	TransportBus   = "bus"
	TransportTrain = "train"
)

type Event struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	EventName            string         `gorm:"size:200;not null" json:"name"`
	Date                 time.Time      `gorm:"not null;index" json:"date"`
	TimeSlot             string         `gorm:"size:20;not null" json:"time_slot"`
	ReportingTime        string         `gorm:"size:10;not null" json:"reporting_time"`
	WorkType             string         `gorm:"size:50;not null" json:"work_type"`
	LocationLink         string         `gorm:"size:255" json:"location_link"`
	
	RequiredCaptains     uint           `gorm:"default:0" json:"required_captains"`
	RequiredSubCaptains  uint           `gorm:"default:0" json:"required_sub_captains"`
	RequiredMainBoys     uint           `gorm:"default:0" json:"required_main_boys"`
	RequiredJuniors      uint           `gorm:"default:0" json:"required_juniors"`
	
	RemainingCaptains    uint           `gorm:"not null;default:0" json:"remaining_captains"`
	RemainingSubCaptains uint           `gorm:"not null;default:0" json:"remaining_sub_captains"`
	RemainingMainBoys    uint           `gorm:"not null;default:0" json:"remaining_main_boys"`
	RemainingJuniors     uint           `gorm:"not null;default:0" json:"remaining_juniors"`
	
	LongWork             bool           `gorm:"default:false" json:"long_work"`
	TransportProvided    bool           `gorm:"default:false" json:"transport_provided"`
	TransportType        string         `gorm:"size:20" json:"transport_type"`
	ExtraWageAmount      int64          `gorm:"default:0" json:"extra_wage_amount"`
	
	Status               string         `gorm:"size:20;default:'upcoming';index" json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}