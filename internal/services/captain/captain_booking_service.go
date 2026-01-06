package captain

import (
	"errors"
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CaptainBookingResponse struct {
	Event     models.Event `json:"event"`
	MyBooking BookingDTO  `json:"my_booking"`
}

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

// ======================= INTERNAL HELPER =======================
func mapBookingResponses(bookings []models.Booking) []CaptainBookingResponse {
	res := make([]CaptainBookingResponse, 0, len(bookings))

	for _, b := range bookings {
		res = append(res, CaptainBookingResponse{
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

// ======================= TODAY BOOKINGS =======================
func (s *CaptainBookingService) ListTodayBookings(userID uint) ([]CaptainBookingResponse, error) {
	var bookings []models.Booking
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	err := config.DB.
		Preload("Event").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(`
			bookings.user_id = ?
			AND events.date >= ?
			AND events.date < ?
			AND bookings.deleted_at IS NULL
		`, userID, start, end).
		Order("events.date ASC").
		Find(&bookings).Error

	return mapBookingResponses(bookings), err
}

// ======================= UPCOMING BOOKINGS =======================
func (s *CaptainBookingService) ListUpcomingBookings(userID uint) ([]CaptainBookingResponse, error) {
	var bookings []models.Booking
	today := time.Now().Truncate(24 * time.Hour)

	err := config.DB.
		Preload("Event").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(`
			bookings.user_id = ?
			AND events.date > ?
			AND bookings.deleted_at IS NULL
		`, userID, today).
		Order("events.date ASC").
		Find(&bookings).Error

	return mapBookingResponses(bookings), err
}

// ======================= COMPLETED BOOKINGS =======================
func (s *CaptainBookingService) ListCompletedBookings(userID uint) ([]CaptainBookingResponse, error) {
	var bookings []models.Booking

	err := config.DB.
		Preload("Event").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where(`
			bookings.user_id = ?
			AND events.status = ?
			AND bookings.deleted_at IS NULL
		`, userID, models.EventStatusCompleted).
		Order("events.date DESC").
		Find(&bookings).Error

	return mapBookingResponses(bookings), err
}

// ======================= LIST EVENT BOOKINGS =======================
func (s *CaptainBookingService) ListEventBookings(captainID, eventID uint) ([]AttendanceRowResponse, error) {
	var count int64
	if err := config.DB.Model(&models.Booking{}).
		Where("event_id=? AND user_id=? AND role=? AND deleted_at IS NULL",
			eventID, captainID, models.RoleCaptain).
		Count(&count).Error; err != nil || count == 0 {
		return nil, errors.New("you are not authorized to view attendance for this event")
	}

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
		Where("bookings.event_id = ? AND bookings.deleted_at IS NULL", eventID).
		Order("bookings.role ASC").
		Scan(&rows).Error

	return rows, err
}

// ======================= UPDATE ATTENDANCE =======================
func (s *CaptainBookingService) UpdateAttendance(
	captainID uint,
	bookingID uint,
	status string,
	ta, bonus, fine int64,
) error {

	return config.DB.Transaction(func(tx *gorm.DB) error {

		var booking models.Booking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", bookingID).
			First(&booking).Error; err != nil {
			return errors.New("booking not found")
		}

		var event models.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", booking.EventID).
			First(&event).Error; err != nil {
			return errors.New("event not found")
		}

		if event.Status != models.EventStatusOngoing {
			return errors.New("attendance can be updated only during ongoing events")
		}

		var captainBooking models.Booking
		if err := tx.Where(
			"event_id=? AND user_id=? AND role=? AND deleted_at IS NULL",
			event.ID, captainID, models.RoleCaptain,
		).First(&captainBooking).Error; err != nil {
			return errors.New("you are not authorized to update attendance")
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

		if status == models.BookingStatusAbsent {
			booking.BaseAmount = 0
			booking.ExtraAmount = 0
			booking.TAAmount = 0
			booking.BonusAmount = 0
			booking.FineAmount = 0
			booking.TotalAmount = 0
			return tx.Save(&booking).Error
		}

		if ta < 0 || bonus < 0 || fine < 0 {
			return errors.New("amounts cannot be negative")
		}

		booking.TAAmount = ta
		booking.BonusAmount = bonus
		booking.FineAmount = fine
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

		return tx.Save(&booking).Error
	})
}
// ======================= FILTER BY STATUS =======================
func (s *CaptainBookingService) ListEventBookingsByStatus(
	captainID uint,
	eventID uint,
	status string,
) ([]AttendanceRowResponse, error) {

	switch status {
	case models.BookingStatusBooked,
		models.BookingStatusPresent,
		models.BookingStatusAbsent,
		models.BookingStatusCompleted:
	default:
		return []AttendanceRowResponse{}, errors.New("invalid booking status")
	}

	if err := s.verifyCaptain(captainID, eventID); err != nil {
		return []AttendanceRowResponse{}, err
	}

	rows := make([]AttendanceRowResponse, 0)

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
		Where(`
			bookings.event_id = ?
			AND bookings.status = ?
			AND bookings.deleted_at IS NULL
		`, eventID, status).
		Order("users.name ASC").
		Scan(&rows).Error

	return rows, err
}

// ======================= SEARCH BY NAME =======================
func (s *CaptainBookingService) SearchEventBookingsByName(
	captainID uint,
	eventID uint,
	name string,
) ([]AttendanceRowResponse, error) {
	if err := s.verifyCaptain(captainID, eventID); err != nil {
		return []AttendanceRowResponse{}, err
	}

	rows := make([]AttendanceRowResponse, 0)

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
		Where(`
			bookings.event_id = ?
			AND users.name ILIKE ?
			AND bookings.deleted_at IS NULL
		`, eventID, "%"+name+"%").
		Order("users.name ASC").
		Scan(&rows).Error

	return rows, err
}

// ======================= EVENT WAGE SUMMERY =======================
func (s *CaptainBookingService) GetEventWageSummary(
	eventID uint,
) (*EventWageSummary, error) {

	var summary EventWageSummary

	err := config.DB.
		Table("bookings").
		Select(`
			COUNT(id) as total_workers,
			COALESCE(SUM(base_amount),0)  as base_total,
			COALESCE(SUM(extra_amount),0) as extra_total,
			COALESCE(SUM(ta_amount),0)    as ta_total,
			COALESCE(SUM(bonus_amount),0) as bonus_total,
			COALESCE(SUM(fine_amount),0)  as fine_total,
			COALESCE(SUM(total_amount),0) as grand_total
		`).
		Where(`
			event_id = ?
			AND deleted_at IS NULL
			AND status != ?
		`, eventID, models.BookingStatusAbsent).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}
// ======================= INTERNAL =======================
func (s *CaptainBookingService) verifyCaptain(captainID, eventID uint) error {
	var count int64
	if err := config.DB.Model(&models.Booking{}).
		Where(
			"event_id=? AND user_id=? AND role=? AND deleted_at IS NULL",
			eventID, captainID, models.RoleCaptain,
		).
		Count(&count).Error; err != nil || count == 0 {
		return errors.New("you are not authorized for this event")
	}
	return nil
}