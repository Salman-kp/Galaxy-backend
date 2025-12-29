package admin

import (
	"errors"

	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
)

type RoleWageService struct {
	repo interfaces.RoleWageRepository
}

func NewRoleWageService(repo interfaces.RoleWageRepository) *RoleWageService {
	return &RoleWageService{repo: repo}
}

func (s *RoleWageService) GetAll() ([]models.RoleWage, error) {
	return s.repo.GetAll()
}

func (s *RoleWageService) Update(role string, wage int64) error {
	if wage <= 0 {
		return errors.New("wage must be greater than zero")
	}

	existing, err := s.repo.FindByRole(role)
	if err != nil {
		return errors.New("role not found")
	}

	existing.Wage = wage
	return s.repo.Update(existing)
}