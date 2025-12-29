package captain   

import (
	"errors"
	"time"

	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
)

type CaptainEventService struct {
	repo interfaces.EventRepository
}

func NewCaptainEventService(repo interfaces.EventRepository) *CaptainEventService {
	return &CaptainEventService{repo: repo}
}

// ---------------- VIEW ----------------

func (s *CaptainEventService) ListAvailableEvents() ([]models.Event, error) {
	return s.repo.ListAvailable(time.Now())
}

func (s *CaptainEventService) GetEvent(id uint) (*models.Event, error) {
	return s.repo.FindByID(id)
}

// ---------------- STATUS CONTROL ----------------

func (s *CaptainEventService) StartEvent(id uint) error {
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

func (s *CaptainEventService) CompleteEvent(id uint) error {
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