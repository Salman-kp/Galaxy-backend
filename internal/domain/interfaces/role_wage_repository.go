package interfaces

import "event-management-backend/internal/domain/models"

type RoleWageRepository interface {
	Create(wage *models.RoleWage) error
	Update(wage *models.RoleWage) error
	Delete(role string) error
	GetAll() ([]models.RoleWage, error)
	FindByRole(role string) (*models.RoleWage, error)
}