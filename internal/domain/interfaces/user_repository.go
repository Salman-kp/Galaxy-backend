package interfaces

import "event-management-backend/internal/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	ListAll(role string, status string) ([]models.User, error)
	FindAll() ([]models.User, error)
	Count() (int64, error)
	UpdateFields(id uint, updates map[string]interface{}) error
	Update(user *models.User) error
	SoftDelete(id uint) error
}