package admin

import (
	"errors"
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type AdminEventService struct {
	repo interfaces.EventRepository
}

func NewAdminEventService(repo interfaces.EventRepository) *AdminEventService {
	return &AdminEventService{repo: repo}
}

// ---------------- CREATE ----------------

func (s *AdminEventService) CreateEvent(event *models.Event) error {
	today := time.Now().Truncate(24 * time.Hour)
	if event.Date.Before(today) {
		return errors.New("event date cannot be in the past")
	}

	event.Status = models.EventStatusUpcoming
	event.RemainingCaptains = event.RequiredCaptains
	event.RemainingSubCaptains = event.RequiredSubCaptains
	event.RemainingMainBoys = event.RequiredMainBoys
	event.RemainingJuniors = event.RequiredJuniors

	if !event.TransportProvided {
		event.TransportType = ""
	}
	if !event.LongWork {
		event.ExtraWageAmount = 0
	}

	return s.repo.Create(event)
}

// ---------------- READ ----------------

func (s *AdminEventService) GetEvent(id uint) (*models.Event, error) {
	return s.repo.FindByID(id)
}

func (s *AdminEventService) ListEvents(status, date string) ([]models.Event, error) {
	return s.repo.ListAll(status, date)
}

// ---------------- UPDATE ----------------

func (s *AdminEventService) UpdateEvent(input *models.Event) error {
	old, err := s.repo.FindByID(input.ID)
	if err != nil {
		return err
	}

	if old.Status != models.EventStatusUpcoming {
		return errors.New("event cannot be updated after start")
	}

	changed := false

	if input.EventName != "" && old.EventName != input.EventName {
		old.EventName = input.EventName
		changed = true
	}

	if !input.Date.IsZero() && !old.Date.Equal(input.Date) {
		old.Date = input.Date
		changed = true
	}

	if input.TimeSlot != "" && old.TimeSlot != input.TimeSlot {
		old.TimeSlot = input.TimeSlot
		changed = true
	}

	if input.ReportingTime != "" && old.ReportingTime != input.ReportingTime {
		old.ReportingTime = input.ReportingTime
		changed = true
	}

	if input.WorkType != "" && old.WorkType != input.WorkType {
		old.WorkType = input.WorkType
		changed = true
	}

	if input.LocationLink != "" && old.LocationLink != input.LocationLink {
		old.LocationLink = input.LocationLink
		changed = true
	}

	// -------- REQUIRED COUNTS (SAFE UPDATE) --------

	if input.RequiredCaptains != old.RequiredCaptains {
		booked := old.RequiredCaptains - old.RemainingCaptains
		if input.RequiredCaptains < booked {
			return errors.New("required captains less than already booked")
		}
		old.RequiredCaptains = input.RequiredCaptains
		old.RemainingCaptains = input.RequiredCaptains - booked
		changed = true
	}

	if input.RequiredSubCaptains != old.RequiredSubCaptains {
		booked := old.RequiredSubCaptains - old.RemainingSubCaptains
		if input.RequiredSubCaptains < booked {
			return errors.New("required sub captains less than already booked")
		}
		old.RequiredSubCaptains = input.RequiredSubCaptains
		old.RemainingSubCaptains = input.RequiredSubCaptains - booked
		changed = true
	}

	if input.RequiredMainBoys != old.RequiredMainBoys {
		booked := old.RequiredMainBoys - old.RemainingMainBoys
		if input.RequiredMainBoys < booked {
			return errors.New("required main boys less than already booked")
		}
		old.RequiredMainBoys = input.RequiredMainBoys
		old.RemainingMainBoys = input.RequiredMainBoys - booked
		changed = true
	}

	if input.RequiredJuniors != old.RequiredJuniors {
		booked := old.RequiredJuniors - old.RemainingJuniors
		if input.RequiredJuniors < booked {
			return errors.New("required juniors less than already booked")
		}
		old.RequiredJuniors = input.RequiredJuniors
		old.RemainingJuniors = input.RequiredJuniors - booked
		changed = true
	}

	// -------- FLAGS --------

	if old.LongWork != input.LongWork {
		old.LongWork = input.LongWork
		if !old.LongWork {
			old.ExtraWageAmount = 0
		}
		changed = true
	}

	if old.TransportProvided != input.TransportProvided {
		old.TransportProvided = input.TransportProvided
		if !old.TransportProvided {
			old.TransportType = ""
		}
		changed = true
	}

	if input.TransportType != "" && old.TransportType != input.TransportType {
		old.TransportType = input.TransportType
		changed = true
	}

	if old.ExtraWageAmount != input.ExtraWageAmount {
		old.ExtraWageAmount = input.ExtraWageAmount
		changed = true
	}

	if !changed {
		return errors.New("no changes detected")
	}

	return s.repo.Update(old)
}

// ---------------- STATUS CONTROL ----------------

func (s *AdminEventService) StartEvent(id uint) error {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusUpcoming {
		return errors.New("only upcoming events can be started")
	}

	event.Status = models.EventStatusOngoing
	return s.repo.Update(event)
}

func (s *AdminEventService) CompleteEvent(id uint) error {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return err
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
			id,
			models.BookingStatusPresent,
			models.BookingStatusCompleted,
		).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&models.Booking{}).
			Where("event_id = ? AND status = ?", id, models.BookingStatusPresent).
			Update("status", models.BookingStatusCompleted).
			Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *AdminEventService) CancelEvent(id uint) error {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusUpcoming {
		return errors.New("only upcoming events can be cancelled")
	}

	event.Status = models.EventStatusCancelled
	return s.repo.Update(event)
}

// ---------------- DELETE ----------------

func (s *AdminEventService) DeleteEvent(id uint) error {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("event not found")
	}
	if event.Status != models.EventStatusUpcoming {
		return errors.New("only upcoming events can be deleted")
	}
	return s.repo.SoftDelete(id)
}