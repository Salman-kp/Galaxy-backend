package repository

import (
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type eventRepository struct{}

func NewEventRepository() interfaces.EventRepository {
	return &eventRepository{}
}

func (r *eventRepository) Create(event *models.Event) error {
	return config.DB.Create(event).Error
}

func (r *eventRepository) FindByID(id uint) (*models.Event, error) {
	var event models.Event
	err := config.DB.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&event).Error
	return &event, err
}

func (r *eventRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Event, error) {
	var event models.Event
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&event).Error
	return &event, err
}

func (r *eventRepository) ListAll(status string, date string) ([]models.Event, error) {
	var events []models.Event
	q := config.DB.Model(&models.Event{}).
		Where("deleted_at IS NULL")

	if status != "" {
		q = q.Where("status = ?", status)
	}

	if date != "" {
		if d, err := time.Parse("2006-01-02", date); err == nil {
			start := d
			end := d.Add(24 * time.Hour)
			q = q.Where("date >= ? AND date < ?", start, end)
		}
	}

	err := q.Order("date ASC").Find(&events).Error
	return events, err
}

func (r *eventRepository) ListAvailable(date time.Time) ([]models.Event, error) {
	var events []models.Event

	q := config.DB.Model(&models.Event{}).
		Where("deleted_at IS NULL").
		Where("status IN ?", []string{
			models.EventStatusUpcoming,
			models.EventStatusOngoing,
		})

	if !date.IsZero() {
		q = q.Where("date >= ?", date)
	}

	err := q.Order("date ASC").Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(event *models.Event) error {
	return config.DB.Save(event).Error
}

func (r *eventRepository) SoftDelete(id uint) error {
	res := config.DB.Delete(&models.Event{}, id)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}