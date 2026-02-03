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

// ---------------- CAPTAIN ----------------
func (r *eventRepository) ListAvailableForCaptain(userID uint, date time.Time) ([]models.Event, error) {
	return r.listAvailableByRole(userID, date, "remaining_captains > 0")
}

// ---------------- SUB CAPTAIN ----------------
func (r *eventRepository) ListAvailableForSubCaptain(userID uint, date time.Time) ([]models.Event, error) {
	return r.listAvailableByRole(userID, date, "remaining_sub_captains > 0")
}

// ---------------- MAIN BOY ----------------
func (r *eventRepository) ListAvailableForMainBoy(userID uint, date time.Time) ([]models.Event, error) {
	return r.listAvailableByRole(userID, date, "remaining_main_boys > 0")
}

// ---------------- JUNIOR ----------------
func (r *eventRepository) ListAvailableForJunior(userID uint, date time.Time) ([]models.Event, error) {
	return r.listAvailableByRole(userID, date, "remaining_juniors > 0")
}

// ---------------- COMMON INTERNAL QUERY ----------------
func (r *eventRepository) listAvailableByRole(
	userID uint,
	date time.Time,
	roleCondition string,
) ([]models.Event, error) {

	var events []models.Event

    q := config.DB.Model(&models.Event{}).
        Where("events.deleted_at IS NULL").
        Where("events.status = ?", models.EventStatusUpcoming).
        Where(roleCondition)

    q = q.Where(`
        NOT EXISTS (
            SELECT 1 FROM bookings 
            WHERE bookings.event_id = events.id 
            AND bookings.user_id = ? 
            AND bookings.deleted_at IS NULL
        )
    `, userID)

	q = q.Where(`
        NOT EXISTS (
            SELECT 1 FROM bookings b
            JOIN events e ON b.event_id = e.id
            WHERE b.user_id = ? 
            AND b.deleted_at IS NULL 
            AND e.deleted_at IS NULL
            AND DATE(e.date) = DATE(events.date)
        )
    `, userID)

    if !date.IsZero() {
        q = q.Where("DATE(events.date) >= DATE(?)", date)
    }

    err := q.Order("events.date ASC").Find(&events).Error
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