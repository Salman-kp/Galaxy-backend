package admin

import (
	"errors"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WageService struct {
	bookingRepo interfaces.BookingRepository
	eventRepo   interfaces.EventRepository
}

func NewWageService(
	bookingRepo interfaces.BookingRepository,
	eventRepo interfaces.EventRepository,
) *WageService {
	return &WageService{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

//
// ---------------- OVERRIDE WAGE (ADMIN) ----------------
// RULES:
// - ONLY COMPLETED EVENTS
// - Only TA / Bonus / Fine
// - ZERO allowed
// - Negative NOT allowed
//
func (s *WageService) OverrideWage(
	bookingID uint,
	taAmount int64,
	bonusAmount int64,
	fineAmount int64,
) error {

	// ---------------- VALIDATION ----------------
	if taAmount < 0 || bonusAmount < 0 || fineAmount < 0 {
		return errors.New("amounts cannot be negative")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {

		// ---------------- LOCK BOOKING ----------------
		var booking models.Booking
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", bookingID).
			First(&booking).Error; err != nil {
			return errors.New("booking not found")
		}

		// ---------------- LOCK EVENT ----------------
		event, err := s.eventRepo.FindByIDForUpdate(tx, booking.EventID)
		if err != nil {
			return err
		}

		// ---------------- COMPLETED ONLY ----------------
		if event.Status != models.EventStatusCompleted {
			return errors.New("wage override allowed only for completed events")
		}

		// ---------------- ABSENT SAFETY ----------------
		if booking.Status == models.BookingStatusAbsent {
			return errors.New("cannot override wage for absent booking")
		}

		// ---------------- APPLY OVERRIDE ----------------
		booking.TAAmount = taAmount
		booking.BonusAmount = bonusAmount
		booking.FineAmount = fineAmount

		// ---------------- RECALCULATE TOTAL ----------------
		booking.TotalAmount =
			booking.BaseAmount +
				booking.ExtraAmount +
				booking.TAAmount +
				booking.BonusAmount -
				booking.FineAmount

		if booking.TotalAmount < 0 {
			booking.TotalAmount = 0
		}

		return tx.Save(&booking).Error
	})
}