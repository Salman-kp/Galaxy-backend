package admin

import (
	"errors"

	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
)

type WageService struct {
	bookingRepo interfaces.BookingRepository
}

func NewWageService(bookingRepo interfaces.BookingRepository) *WageService {
	return &WageService{
		bookingRepo: bookingRepo,
	}
}

//
// ---------------- CALCULATE FINAL AMOUNT ----------------
// Used internally after attendance updates
//
func (s *WageService) CalculateFinalAmount(
	booking *models.Booking,
	baseAmount int64,
	event *models.Event,
) error {

	// ABSENT → everything zero
	if booking.Status == models.BookingStatusAbsent {
		booking.BaseAmount = 0
		booking.ExtraAmount = 0
		booking.TAAmount = 0
		booking.BonusAmount = 0
		booking.FineAmount = 0
		booking.TotalAmount = 0
		return nil
	}

	// Base amount validation
	if baseAmount < 0 {
		return errors.New("base amount cannot be negative")
	}

	booking.BaseAmount = baseAmount

	// Extra wage (long work)
	if event.LongWork {
		booking.ExtraAmount = event.ExtraWageAmount
	} else {
		booking.ExtraAmount = 0
	}

	// Final calculation
	booking.TotalAmount =
		booking.BaseAmount +
			booking.ExtraAmount +
			booking.TAAmount +
			booking.BonusAmount -
			booking.FineAmount

	// Safety guard
	if booking.TotalAmount < 0 {
		booking.TotalAmount = 0
	}

	return nil
}

//
// ---------------- OVERRIDE WAGE (ADMIN) ----------------
// Admin can override ONLY: TA, Bonus, Fine
// Zero is allowed, negative is NOT allowed
//
func (s *WageService) OverrideWage(
	bookingID uint,
	taAmount int64,
	bonusAmount int64,
	fineAmount int64,
) error {

	// Validation (ZERO ALLOWED)
	if taAmount < 0 || bonusAmount < 0 || fineAmount < 0 {
		return errors.New("amounts cannot be negative")
	}

	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	booking.TAAmount = taAmount
	booking.BonusAmount = bonusAmount
	booking.FineAmount = fineAmount

	booking.TotalAmount =
		booking.BaseAmount +
			booking.ExtraAmount +
			booking.TAAmount +
			booking.BonusAmount -
			booking.FineAmount

	if booking.TotalAmount < 0 {
		booking.TotalAmount = 0
	}

	return s.bookingRepo.Update(booking)
}