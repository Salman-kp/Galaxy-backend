package captain

import (
	"errors"

	"event-management-backend/internal/domain/models"
)

type WageService struct{}

func NewWageService() *WageService {
	return &WageService{}
}

func (s *WageService) CalculateFinalAmount(
	booking *models.Booking,
	baseAmount int64,
	event *models.Event,
) error {

	if booking.Status == models.BookingStatusAbsent {
		booking.BaseAmount = 0
		booking.ExtraAmount = 0
		booking.TAAmount = 0
		booking.BonusAmount = 0
		booking.FineAmount = 0
		booking.TotalAmount = 0
		return nil
	}

	if baseAmount < 0 {
		return errors.New("invalid base amount")
	}

	booking.BaseAmount = baseAmount

	if event.LongWork {
		booking.ExtraAmount = event.ExtraWageAmount
	} else {
		booking.ExtraAmount = 0
	}

	booking.TotalAmount =
		booking.BaseAmount +
			booking.ExtraAmount +
			booking.TAAmount +
			booking.BonusAmount -
			booking.FineAmount

	if booking.TotalAmount < 0 {
		booking.TotalAmount = 0
	}

	return nil
}