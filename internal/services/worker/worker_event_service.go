package worker    

import (
	"time"

	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
)

type WorkerEventService struct {
	repo interfaces.EventRepository
}

func NewWorkerEventService(repo interfaces.EventRepository) *WorkerEventService {
	return &WorkerEventService{repo: repo}
}

// ---------------- VIEW ONLY ----------------

func (s *WorkerEventService) ListAvailableEvents() ([]models.Event, error) {
	return s.repo.ListAvailable(time.Now())
}

func (s *WorkerEventService) GetEvent(id uint) (*models.Event, error) {
	return s.repo.FindByID(id)
}