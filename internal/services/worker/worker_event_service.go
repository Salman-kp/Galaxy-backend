package worker

import (
	"errors"
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
func (s *WorkerEventService) ListAvailableEvents(userID uint, role string) ([]models.Event, error) {
	today := time.Now().Truncate(24 * time.Hour)

	switch role {
	case models.RoleSubCaptain:
		return s.repo.ListAvailableForSubCaptain(userID, today)

	case models.RoleMainBoy:
		return s.repo.ListAvailableForMainBoy(userID, today)

	case models.RoleJuniorBoy:
		return s.repo.ListAvailableForJunior(userID, today)

	default:
		return nil, errors.New("invalid worker role")
	}
}

func (s *WorkerEventService) GetEvent(id uint) (*models.Event, error) {
	return s.repo.FindByID(id)
}