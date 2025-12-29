package validations

import (
	"errors"
	"time"

	"event-management-backend/internal/domain/models"
)

/*
|--------------------------------------------------------------------------
| CREATE EVENT REQUEST
|--------------------------------------------------------------------------
*/

type CreateEventRequest struct {
	Name                string    `json:"name"`
	Date                time.Time `json:"date"`
	TimeSlot            string    `json:"time_slot"`
	ReportingTime       string    `json:"reporting_time"`
	WorkType            string    `json:"work_type"`
	LocationLink        string    `json:"location_link"`

	RequiredCaptains    uint `json:"required_captains"`
	RequiredSubCaptains uint `json:"required_sub_captains"`
	RequiredMainBoys    uint `json:"required_main_boys"`
	RequiredJuniors     uint `json:"required_juniors"`

	LongWork          bool   `json:"long_work"`
	TransportProvided bool   `json:"transport_provided"`
	TransportType     string `json:"transport_type"`
	ExtraWageAmount   int64  `json:"extra_wage_amount"`
}

func (r *CreateEventRequest) Validate() error {
	// basic required fields
	if r.Name == "" {
		return errors.New("event name is required")
	}
	if r.Date.IsZero() {
		return errors.New("event date is required")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if r.Date.Before(today) {
		return errors.New("event date cannot be in the past")
	}

	if r.TimeSlot == "" {
		return errors.New("time slot is required")
	}
	if r.ReportingTime == "" {
		return errors.New("reporting time is required")
	}
	if r.WorkType == "" {
		return errors.New("work type is required")
	}

	// time slot validation
	if !isValidTimeSlot(r.TimeSlot) {
		return errors.New("invalid time slot")
	}

	// at least one worker required
	if r.RequiredCaptains == 0 &&
		r.RequiredSubCaptains == 0 &&
		r.RequiredMainBoys == 0 &&
		r.RequiredJuniors == 0 {
		return errors.New("at least one worker is required")
	}

	// extra wage rules
	if r.ExtraWageAmount < 0 {
		return errors.New("extra wage amount cannot be negative")
	}
	if !r.LongWork && r.ExtraWageAmount != 0 {
		return errors.New("extra wage must be zero when long work is false")
	}

	// transport rules
	if r.TransportProvided {
		if !isValidTransportType(r.TransportType) {
			return errors.New("invalid transport type")
		}
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| UPDATE EVENT REQUEST
|--------------------------------------------------------------------------
*/

type UpdateEventRequest struct {
	Name                string    `json:"name"`
	Date                time.Time `json:"date"`
	TimeSlot            string    `json:"time_slot"`
	ReportingTime       string    `json:"reporting_time"`
	WorkType            string    `json:"work_type"`
	LocationLink        string    `json:"location_link"`

	RequiredCaptains    uint `json:"required_captains"`
	RequiredSubCaptains uint `json:"required_sub_captains"`
	RequiredMainBoys    uint `json:"required_main_boys"`
	RequiredJuniors     uint `json:"required_juniors"`

	LongWork          bool   `json:"long_work"`
	TransportProvided bool   `json:"transport_provided"`
	TransportType     string `json:"transport_type"`
	ExtraWageAmount   int64  `json:"extra_wage_amount"`
}

func (r *UpdateEventRequest) Validate() error {
	// optional date check (if provided)
	if !r.Date.IsZero() {
		today := time.Now().Truncate(24 * time.Hour)
		if r.Date.Before(today) {
			return errors.New("event date cannot be in the past")
		}
	}

	// optional time slot check
	if r.TimeSlot != "" && !isValidTimeSlot(r.TimeSlot) {
		return errors.New("invalid time slot")
	}

	// extra wage rules
	if r.ExtraWageAmount < 0 {
		return errors.New("extra wage amount cannot be negative")
	}
	if !r.LongWork && r.ExtraWageAmount != 0 {
		return errors.New("extra wage must be zero when long work is false")
	}

	// transport rules
	if r.TransportProvided {
		if !isValidTransportType(r.TransportType) {
			return errors.New("invalid transport type")
		}
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| HELPERS
|--------------------------------------------------------------------------
*/

func isValidTimeSlot(slot string) bool {
	switch slot {
	case models.TimeSlotMorning,
		models.TimeSlotLunch,
		models.TimeSlotNight:
		return true
	default:
		return false
	}
}

func isValidTransportType(t string) bool {
	switch t {
	case models.TransportBus,
		models.TransportTrain:
		return true
	default:
		return false
	}
}