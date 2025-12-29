package captain

import (
	"errors"
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type CaptainBookingService struct {
	bookingRepo interfaces.BookingRepository
	eventRepo   interfaces.EventRepository
	userRepo    interfaces.UserRepository
}

func NewCaptainBookingService(
	bookingRepo interfaces.BookingRepository,
	eventRepo interfaces.EventRepository,
	userRepo interfaces.UserRepository,
) *CaptainBookingService {
	return &CaptainBookingService{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
		userRepo:    userRepo,
	}
}

// ======================= BOOK EVENT =======================
// Captain books an event for himself
func (s *CaptainBookingService) BookEvent(userID, eventID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return errors.New("user not found")
		}
		event, err := s.eventRepo.FindByIDForUpdate(tx, eventID)
		if err != nil {
			return errors.New("event not found")
		}

		if event.Status != models.EventStatusUpcoming {
			return errors.New("event is not open for booking")
		}

		// prevent duplicate booking
		if _, err := s.bookingRepo.FindByEventAndUser(eventID, userID); err == nil {
			return errors.New("already booked")
		}

		if event.RemainingCaptains == 0 {
			return errors.New("no captain slots available")
		}

		event.RemainingCaptains--

		if err := tx.Save(event).Error; err != nil {
			return err
		}

		booking := &models.Booking{
			EventID:    eventID,
			UserID:     userID,
			Role:       models.RoleCaptain,
			Status:     models.BookingStatusBooked,
			BaseAmount: user.CurrentWage,
		}

		return tx.Create(booking).Error
	})
}

// ======================= LIST MY BOOKINGS =======================
func (s *CaptainBookingService) ListMyBookings(userID uint) ([]models.Booking, error) {
	return s.bookingRepo.ListByUser(userID)
}

// ======================= LIST TODAY BOOKINGS =======================
func (s *CaptainBookingService) ListTodayBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	today := time.Now().Truncate(24 * time.Hour)

	err := config.DB.
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(
			"bookings.user_id = ? AND events.date = ? AND bookings.deleted_at IS NULL",
			userID,
			today,
		).
		Find(&bookings).Error

	return bookings, err
}

// ======================= LIST UPCOMING BOOKINGS =======================
func (s *CaptainBookingService) ListUpcomingBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	today := time.Now().Truncate(24 * time.Hour)

	err := config.DB.
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(
			"bookings.user_id = ? AND events.date > ? AND bookings.deleted_at IS NULL",
			userID,
			today,
		).
		Order("events.date ASC").
		Find(&bookings).Error

	return bookings, err
}

// ======================= LIST COMPLETED BOOKINGS =======================
func (s *CaptainBookingService) ListCompletedBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking

	err := config.DB.
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(
			"bookings.user_id = ? AND events.status = ? AND bookings.deleted_at IS NULL",
			userID,
			models.EventStatusCompleted,
		).
		Order("events.date DESC").
		Find(&bookings).Error

	return bookings, err
}

// ======================= LIST EVENT BOOKINGS =======================
// For attendance table
func (s *CaptainBookingService) ListEventBookings(eventID uint) ([]models.Booking, error) {
	return s.bookingRepo.ListByEvent(eventID)
}

// ======================= UPDATE ATTENDANCE =======================
// Captain updates attendance & wage
func (s *CaptainBookingService) UpdateAttendance(
	captainID uint,
	bookingID uint,
	status string,
	ta int64,
	bonus int64,
	fine int64,
) error {

	return config.DB.Transaction(func(tx *gorm.DB) error {

		// Lock booking
		booking, err := s.bookingRepo.FindByIDForUpdate(tx, bookingID)
		if err != nil {
			return errors.New("booking not found")
		}

		// 🔒 OWNERSHIP CHECK (THIS IS THE LINE YOU ASKED ABOUT)
		if booking.UserID != captainID {
			return errors.New("you are not assigned to this event")
		}

		// Validate status
		switch status {
		case models.BookingStatusBooked,
			models.BookingStatusPresent,
			models.BookingStatusCompleted,
			models.BookingStatusAbsent:
		default:
			return errors.New("invalid booking status")
		}

		// Lock event
		event, err := s.eventRepo.FindByIDForUpdate(tx, booking.EventID)
		if err != nil {
			return errors.New("event not found")
		}

		booking.Status = status

		// ABSENT → wipe everything
		if status == models.BookingStatusAbsent {
			booking.BaseAmount = 0
			booking.ExtraAmount = 0
			booking.TAAmount = 0
			booking.BonusAmount = 0
			booking.FineAmount = 0
			booking.TotalAmount = 0

			return s.bookingRepo.Update(booking)
		}

		// Validate amounts
		if ta < 0 || bonus < 0 || fine < 0 {
			return errors.New("amounts cannot be negative")
		}

		booking.TAAmount = ta
		booking.BonusAmount = bonus
		booking.FineAmount = fine

		// Extra wage
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

		return s.bookingRepo.Update(booking)
	})
}
