package admin

import (
	"errors"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type RoleWageService struct {
	repo     interfaces.RoleWageRepository
	userRepo interfaces.UserRepository
}

func NewRoleWageService(
	repo interfaces.RoleWageRepository,
	userRepo interfaces.UserRepository,
) *RoleWageService {
	return &RoleWageService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *RoleWageService) GetAll() ([]models.RoleWage, error) {
	return s.repo.GetAll()
}

func (s *RoleWageService) Update(role string, wage int64) error {
	if !models.ValidateRole(role) {
		return errors.New("invalid role")
	}

	if wage <= 0 {
		return errors.New("wage must be greater than zero")
	}

	existing, err := s.repo.FindByRole(role)
	if err != nil {
		return errors.New("role not found")
	}

	if existing.Wage == wage {
		return errors.New("no changes detected")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		existing.Wage = wage
		if err := tx.Save(existing).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&models.User{}).
			Where("role = ? AND deleted_at IS NULL", role).
			Update("current_wage", wage).Error; err != nil {
			return err
		}

		return nil
	})
}