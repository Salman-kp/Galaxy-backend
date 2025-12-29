package interfaces

import (
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *models.Booking) error
	FindByID(id uint) (*models.Booking, error)
	FindByEventAndUser(eventID uint, userID uint) (*models.Booking, error)
	ListByUser(userID uint) ([]models.Booking, error)
	ListByEvent(eventID uint) ([]models.Booking, error)
	FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Booking, error)
	Update(booking *models.Booking) error
	Delete(id uint) error
}
