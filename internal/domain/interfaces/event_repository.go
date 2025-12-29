package interfaces

import (
	"time"

	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *models.Event) error
	Update(event *models.Event) error
	FindByID(id uint) (*models.Event, error)
	FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Event, error)
	ListAll(status string, date string) ([]models.Event, error)
	ListAvailable(fromDate time.Time) ([]models.Event, error)
	SoftDelete(id uint) error
}
