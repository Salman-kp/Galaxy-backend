package worker    

import (
	"errors"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type WorkerBookingService struct {
	bookingRepo interfaces.BookingRepository
	eventRepo   interfaces.EventRepository
	userRepo    interfaces.UserRepository
}

func NewWorkerBookingService(
	bookingRepo interfaces.BookingRepository,
	eventRepo interfaces.EventRepository,
	userRepo interfaces.UserRepository,
) *WorkerBookingService {
	return &WorkerBookingService{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
		userRepo:    userRepo,
	}
}

// Worker books an event
func (s *WorkerBookingService) BookEvent(userID, eventID uint, role string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		// 🔹 Fetch user (FOR WAGE)
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return errors.New("user not found")
		}

		// 🔹 Lock event
		event, err := s.eventRepo.FindByIDForUpdate(tx, eventID)
		if err != nil {
			return err
		}

		if event.Status != models.EventStatusUpcoming {
			return errors.New("event is not open for booking")
		}

		// prevent duplicate booking
		if _, err := s.bookingRepo.FindByEventAndUser(eventID, userID); err == nil {
			return errors.New("already booked")
		}

		// 🔹 slot validation
		switch role {
		case models.RoleSubCaptain:
			if event.RemainingSubCaptains == 0 {
				return errors.New("no sub-captain slots available")
			}
			event.RemainingSubCaptains--
		case models.RoleMainBoy:
			if event.RemainingMainBoys == 0 {
				return errors.New("no main boy slots available")
			}
			event.RemainingMainBoys--
		case models.RoleJuniorBoy:
			if event.RemainingJuniors == 0 {
				return errors.New("no junior slots available")
			}
			event.RemainingJuniors--
		default:
			return errors.New("invalid role")
		}

		if err := tx.Save(event).Error; err != nil {
			return err
		}

		// 🔹 APPLY ROLE WAGE HERE (IMPORTANT)
		booking := &models.Booking{
			EventID:    eventID,
			UserID:     userID,
			Role:       role,
			Status:     models.BookingStatusBooked,
			BaseAmount: user.CurrentWage, // ✅ FIXED
		}

		return tx.Create(booking).Error
	})
}


func (s *WorkerBookingService) ListMyBookings(userID uint) ([]models.Booking, error) {
	return s.bookingRepo.ListByUser(userID)
}

func (s *WorkerBookingService) ListCompletedBookings(userID uint) ([]models.Booking, error) {
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

