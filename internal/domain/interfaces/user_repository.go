package interfaces

import "event-management-backend/internal/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	ListAll(role string, status string) ([]models.User, error)
    FindAll() ([]models.User, error)
	Count() (int64, error)
	Update(user *models.User) error
	SoftDelete(id uint) error
}