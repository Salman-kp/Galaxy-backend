package validations

import (
	"errors"

	"event-management-backend/internal/domain/models"
)

type UpdateAttendanceRequest struct {
	BookingID   uint   `json:"booking_id" binding:"required"`
	Status      string `json:"status" binding:"required"`
	TAAmount    int64  `json:"ta_amount"`
	BonusAmount int64  `json:"bonus_amount"`
	FineAmount  int64  `json:"fine_amount"`
}

func (r *UpdateAttendanceRequest) Validate() error {
	// ---------------- BOOKING ID ----------------
	if r.BookingID == 0 {
		return errors.New("booking id is required")
	}

	// ---------------- STATUS ----------------
	switch r.Status {
	case models.BookingStatusBooked,
		models.BookingStatusPresent,
		models.BookingStatusCompleted,
		models.BookingStatusAbsent:
	default:
		return errors.New("invalid booking status")
	}

	// ---------------- AMOUNTS ----------------
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
	if r.Status == models.BookingStatusAbsent {
		if r.TAAmount != 0 || r.BonusAmount != 0 || r.FineAmount != 0 {
			return errors.New("amounts must be zero when status is absent")
		}
	}

	return nil
}