package admin

import (
	"errors"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdminBookingService struct {
	bookingRepo interfaces.BookingRepository
	eventRepo   interfaces.EventRepository
}

func NewAdminBookingService(
	bookingRepo interfaces.BookingRepository,
	eventRepo interfaces.EventRepository,
) *AdminBookingService {
	return &AdminBookingService{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

//
// ---------------- LIST EVENT BOOKINGS ----------------
//
func (s *AdminBookingService) ListEventBookings(
	eventID uint,
) ([]AttendanceRowResponse, error) {

	var rows []AttendanceRowResponse

	err := config.DB.
		Table("bookings").
		Select(`
			bookings.id AS booking_id,
			bookings.user_id,
			users.name AS user_name,
			bookings.role,
			bookings.status,
			bookings.base_amount,
			bookings.extra_amount,
			bookings.ta_amount,
			bookings.bonus_amount,
			bookings.fine_amount,
			bookings.total_amount
		`).
		Joins("JOIN users ON users.id = bookings.user_id").
		Where(
			"bookings.event_id = ? AND bookings.deleted_at IS NULL",
			eventID,
		).
		Order("bookings.role ASC").
		Scan(&rows).Error

	return rows, err
}


//
// ---------------- REMOVE USER FROM EVENT ----------------
// RULE: ONLY UPCOMING EVENTS
//
func (s *AdminBookingService) RemoveUserFromEvent(eventID, bookingID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		booking, err := s.bookingRepo.FindByIDForUpdate(tx, bookingID)
		if err != nil {
			return errors.New("booking not found")
		}

		if booking.EventID != eventID {
			return errors.New("booking does not belong to this event")
		}

		event, err := s.eventRepo.FindByIDForUpdate(tx, eventID)
		if err != nil {
			return err
		}

		if event.Status != models.EventStatusUpcoming {
			return errors.New("booking removal allowed only for upcoming events")
		}

		// restore slot
		switch booking.Role {
		case models.RoleCaptain:
			event.RemainingCaptains++
		case models.RoleSubCaptain:
			event.RemainingSubCaptains++
		case models.RoleMainBoy:
			event.RemainingMainBoys++
		case models.RoleJuniorBoy:
			event.RemainingJuniors++
		}

		if err := tx.Save(event).Error; err != nil {
			return err
		}

		return s.bookingRepo.DeleteTx(tx, booking.ID)
	})
}

//
// ---------------- UPDATE ATTENDANCE ----------------
// RULE: ONLY ONGOING EVENTS
// RETURNS: UPDATED BOOKING (WITH TOTAL AMOUNT)
//
func (s *AdminBookingService) UpdateAttendance(
	bookingID uint,
	status string,
	taAmount int64,
	bonusAmount int64,
	fineAmount int64,
) (*models.Booking, error) {

	var updatedBooking *models.Booking

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		var booking models.Booking
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", bookingID).
			First(&booking).Error; err != nil {
			return errors.New("booking not found")
		}

		event, err := s.eventRepo.FindByIDForUpdate(tx, booking.EventID)
		if err != nil {
			return err
		}

		if event.Status != models.EventStatusOngoing {
			return errors.New("attendance can be updated only during ongoing events")
		}

		switch status {
		case models.BookingStatusBooked,
			models.BookingStatusPresent,
			models.BookingStatusCompleted,
			models.BookingStatusAbsent:
		default:
			return errors.New("invalid booking status")
		}

		booking.Status = status

		// ---------------- ABSENT ----------------
		if status == models.BookingStatusAbsent {
			booking.BaseAmount = 0
			booking.ExtraAmount = 0
			booking.TAAmount = 0
			booking.BonusAmount = 0
			booking.FineAmount = 0
			booking.TotalAmount = 0

			if err := tx.Save(&booking).Error; err != nil {
				return err
			}

			updatedBooking = &booking
			return nil
		}

		// ---------------- VALIDATE AMOUNTS ----------------
		if taAmount < 0 || bonusAmount < 0 || fineAmount < 0 {
			return errors.New("amounts cannot be negative")
		}

		booking.TAAmount = taAmount
		booking.BonusAmount = bonusAmount
		booking.FineAmount = fineAmount

		// ---------------- EXTRA WAGE ----------------
		if event.LongWork {
			booking.ExtraAmount = event.ExtraWageAmount
		} else {
			booking.ExtraAmount = 0
		}

		// ---------------- TOTAL CALCULATION ----------------
		booking.TotalAmount =
			booking.BaseAmount +
				booking.ExtraAmount +
				booking.TAAmount +
				booking.BonusAmount -
				booking.FineAmount

		if booking.TotalAmount < 0 {
			booking.TotalAmount = 0
		}

		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		updatedBooking = &booking
		return nil
	})

	return updatedBooking, err
}