package captain

import (
	"errors"
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type CaptainEventService struct {
	repo interfaces.EventRepository
}

func NewCaptainEventService(repo interfaces.EventRepository) *CaptainEventService {
	return &CaptainEventService{repo: repo}
}

// ---------------- VIEW ----------------
func (s *CaptainEventService) ListAvailableEvents(userID uint) ([]models.Event, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.repo.ListAvailableForCaptain(userID, today)
}

func (s *CaptainEventService) GetEvent(id uint) (*models.Event, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("event not found")
	}
	return event, nil
}

// ---------------- STATUS CONTROL ----------------
func (s *CaptainEventService) StartEvent(captainID, eventID uint) error {
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		return err
	}

	if !s.isCaptainOfEvent(captainID, eventID) {
		return errors.New("you are not authorized to start this event")
	}

	if event.Status != models.EventStatusUpcoming {
		return errors.New("only upcoming events can be started")
	}

	event.Status = models.EventStatusOngoing
	return s.repo.Update(event)
}

func (s *CaptainEventService) CompleteEvent(captainID, eventID uint) error {
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		return err
	}

	if !s.isCaptainOfEvent(captainID, eventID) {
		return errors.New("you are not authorized to complete this event")
	}

	if event.Status != models.EventStatusOngoing {
		return errors.New("only ongoing events can be completed")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {

		event.Status = models.EventStatusCompleted
		if err := tx.Save(event).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE users
			SET completed_work = completed_work + 1
			WHERE id IN (
				SELECT user_id
				FROM bookings
				WHERE event_id = ?
				AND status IN (?, ?)
				AND deleted_at IS NULL
			)
			AND deleted_at IS NULL
		`,
			eventID,
			models.BookingStatusPresent,
			models.BookingStatusCompleted,
		).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&models.Booking{}).
			Where("event_id = ? AND status = ?", eventID, models.BookingStatusPresent).
			Update("status", models.BookingStatusCompleted).
			Error; err != nil {
			return err
		}

		return nil
	})
}

// ---------------- INTERNAL ----------------
func (s *CaptainEventService) isCaptainOfEvent(captainID, eventID uint) bool {
	var count int64
	config.DB.
		Model(&models.Booking{}).
		Where(
			"event_id = ? AND user_id = ? AND role = ? AND deleted_at IS NULL",
			eventID,
			captainID,
			models.RoleCaptain,
		).
		Count(&count)

	return count > 0
}
