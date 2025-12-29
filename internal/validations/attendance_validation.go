package validations

import (
	"errors"

	"event-management-backend/internal/domain/models"
)

type UpdateAttendanceRequest struct {
	Status      string `json:"status" binding:"required"`
	TAAmount    int64  `json:"ta_amount"`
	BonusAmount int64  `json:"bonus_amount"`
	FineAmount  int64  `json:"fine_amount"`
}

func (r *UpdateAttendanceRequest) Validate() error {
	// ---------------- STATUS VALIDATION ----------------
	switch r.Status {
	case models.BookingStatusBooked,
		models.BookingStatusPresent,
		models.BookingStatusCompleted,
		models.BookingStatusAbsent:
		// valid
	default:
		return errors.New("invalid booking status")
	}

	// ---------------- AMOUNT VALIDATION ----------------
	if r.TAAmount < 0 {
		return errors.New("ta amount cannot be negative")
	}
	if r.BonusAmount < 0 {
		return errors.New("bonus amount cannot be negative")
	}
	if r.FineAmount < 0 {
		return errors.New("fine amount cannot be negative")
	}

	// ---------------- ABSENT RULE ----------------
	// If ABSENT → amounts must be zero
	if r.Status == models.BookingStatusAbsent {
		if r.TAAmount != 0 || r.BonusAmount != 0 || r.FineAmount != 0 {
			return errors.New("amounts must be zero when status is absent")
		}
	}

	return nil
}