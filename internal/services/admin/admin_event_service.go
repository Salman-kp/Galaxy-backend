package admin

import (
	"errors"
	"time"

	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
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

	// update allowed fields only
	old.EventName = input.EventName
	old.Date = input.Date
	old.TimeSlot = input.TimeSlot
	old.ReportingTime = input.ReportingTime
	old.WorkType = input.WorkType
	old.LocationLink = input.LocationLink
	old.LongWork = input.LongWork
	old.TransportProvided = input.TransportProvided
	old.TransportType = input.TransportType
	old.ExtraWageAmount = input.ExtraWageAmount

	// required counts can change ONLY before booking module
	old.RequiredCaptains = input.RequiredCaptains
	old.RequiredSubCaptains = input.RequiredSubCaptains
	old.RequiredMainBoys = input.RequiredMainBoys
	old.RequiredJuniors = input.RequiredJuniors

	// reset remaining = required (since booking not started)
	old.RemainingCaptains = input.RequiredCaptains
	old.RemainingSubCaptains = input.RequiredSubCaptains
	old.RemainingMainBoys = input.RequiredMainBoys
	old.RemainingJuniors = input.RequiredJuniors

	// clean optional fields
	if !old.TransportProvided {
		old.TransportType = ""
	}
	if !old.LongWork {
		old.ExtraWageAmount = 0
	}

	return s.repo.Update(old)
}

// ---------------- STATUS CONTROL ----------------

func (s *AdminEventService) StartEvent(id uint) error {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusUpcoming &&
		event.Status != models.EventStatusOngoing {
		return errors.New("event cannot be started")
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

	event.Status = models.EventStatusCompleted
	return s.repo.Update(event)
}

func (s *AdminEventService) CancelEvent(id uint) error {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if event.Status == models.EventStatusCompleted {
		return errors.New("completed event cannot be cancelled")
	}

	event.Status = models.EventStatusCancelled
	return s.repo.Update(event)
}

// ---------------- DELETE ----------------

func (s *AdminEventService) DeleteEvent(id uint) error {
	return s.repo.SoftDelete(id)
}