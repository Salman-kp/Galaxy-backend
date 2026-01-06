package worker

import (
	"errors"
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

// ======================= RESPONSE DTO =======================

type WorkerBookingResponse struct {
	Event     models.Event `json:"event"`
	MyBooking BookingDTO  `json:"my_booking"`
}

type BookingDTO struct {
	ID          uint   `json:"id"`
	EventID     uint   `json:"event_id"`
	Status      string `json:"status"`
	Role        string `json:"role"`
	BaseAmount  int64  `json:"base_amount"`
	ExtraAmount int64  `json:"extra_amount"`
	TAAmount    int64  `json:"ta_amount"`
	BonusAmount int64  `json:"bonus_amount"`
	FineAmount  int64  `json:"fine_amount"`
	TotalAmount int64  `json:"total_amount"`
	CreatedAt   string `json:"created_at"`
}

// ======================= SERVICE =======================

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

// ======================= BOOK EVENT =======================

func (s *WorkerBookingService) BookEvent(userID, eventID uint, role string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return errors.New("user not found")
		}

		if user.Role != role {
			return errors.New("role mismatch")
		}

		event, err := s.eventRepo.FindByIDForUpdate(tx, eventID)
		if err != nil {
			return err
		}

		if event.Status != models.EventStatusUpcoming {
			return errors.New("event is not open for booking")
		}

		if _, err := s.bookingRepo.FindByEventAndUser(eventID, userID); err == nil {
			return errors.New("already booked")
		}

		switch user.Role {
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

		booking := &models.Booking{
			EventID:    eventID,
			UserID:     userID,
			Role:       user.Role,
			Status:     models.BookingStatusBooked,
			BaseAmount: user.CurrentWage,
		}

		return tx.Create(booking).Error
	})
}

// ======================= INTERNAL MAPPER =======================

func mapWorkerBookings(bookings []models.Booking) []WorkerBookingResponse {
	res := make([]WorkerBookingResponse, 0, len(bookings))

	for _, b := range bookings {
		res = append(res, WorkerBookingResponse{
			Event: b.Event,
			MyBooking: BookingDTO{
				ID:          b.ID,
				EventID:     b.EventID,
				Status:      b.Status,
				Role:        b.Role,
				BaseAmount:  b.BaseAmount,
				ExtraAmount: b.ExtraAmount,
				TAAmount:    b.TAAmount,
				BonusAmount: b.BonusAmount,
				FineAmount:  b.FineAmount,
				TotalAmount: b.TotalAmount,
				CreatedAt:   b.CreatedAt.Format(time.RFC3339),
			},
		})
	}

	return res
}

// ======================= LIST MY BOOKINGS =======================

func (s *WorkerBookingService) ListMyBookings(userID uint) ([]WorkerBookingResponse, error) {
	var bookings []models.Booking

	err := config.DB.
		Preload("Event").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(`
			bookings.user_id = ?
			AND events.status IN (?, ?)
			AND bookings.deleted_at IS NULL
		`,
			userID,
			models.EventStatusUpcoming,
			models.EventStatusOngoing,
		).
		Order("events.date ASC").
		Find(&bookings).Error

	return mapWorkerBookings(bookings), err
}

// ======================= LIST COMPLETED BOOKINGS =======================

func (s *WorkerBookingService) ListCompletedBookings(userID uint) ([]WorkerBookingResponse, error) {
	var bookings []models.Booking

	err := config.DB.
		Preload("Event").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(`
			bookings.user_id = ?
			AND events.status = ?
			AND bookings.deleted_at IS NULL
		`,
			userID,
			models.EventStatusCompleted,
		).
		Order("events.date DESC").
		Find(&bookings).Error

	return mapWorkerBookings(bookings), err
}

// ======================= GET BOOKING DETAILS =======================
// Used by booked & completed detail pages
//
func (s *WorkerBookingService) GetBookingDetails(
	userID uint,
	bookingID uint,
) (*WorkerBookingResponse, error) {

	var booking models.Booking

	err := config.DB.
		Preload("Event").
		Where(
			"id = ? AND user_id = ? AND deleted_at IS NULL",
			bookingID,
			userID,
		).
		First(&booking).Error

	if err != nil {
		return nil, errors.New("booking not found or unauthorized")
	}

	res := WorkerBookingResponse{
		Event: booking.Event,
		MyBooking: BookingDTO{
			ID:          booking.ID,
			EventID:     booking.EventID,
			Status:      booking.Status,
			Role:        booking.Role,
			BaseAmount:  booking.BaseAmount,
			ExtraAmount: booking.ExtraAmount,
			TAAmount:    booking.TAAmount,
			BonusAmount: booking.BonusAmount,
			FineAmount:  booking.FineAmount,
			TotalAmount: booking.TotalAmount,
			CreatedAt:   booking.CreatedAt.Format(time.RFC3339),
		},
	}

	return &res, nil
}