package repository

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type bookingRepository struct{}

func NewBookingRepository() interfaces.BookingRepository {
	return &bookingRepository{}
}

//  ---------------- CREATE ----------------

func (r *bookingRepository) Create(booking *models.Booking) error {
	return config.DB.Create(booking).Error
}

// ---------------- FIND BY ID ----------------

func (r *bookingRepository) FindByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := config.DB.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&booking).Error
	return &booking, err
}

// ---------------- FIND BY EVENT + USER ----------------

func (r *bookingRepository) FindByEventAndUser(
	eventID uint,
	userID uint,
) (*models.Booking, error) {

	var booking models.Booking
	err := config.DB.
		Where(
			"event_id = ? AND user_id = ? AND deleted_at IS NULL",
			eventID,
			userID,
		).
		First(&booking).Error

	return &booking, err
}

// ---------------- LIST BY USER ----------------

func (r *bookingRepository) ListByUser(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := config.DB.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&bookings).Error

	return bookings, err
}

//---------------- LIST BY EVENT ----------------

func (r *bookingRepository) ListByEvent(eventID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := config.DB.
		Where("event_id = ? AND deleted_at IS NULL", eventID).
		Order("created_at ASC").
		Find(&bookings).Error

	return bookings, err
}

// ---------------- FIND BY ID (FOR UPDATE) ----------------

func (r *bookingRepository) FindByIDForUpdate(
	tx *gorm.DB,
	id uint,
) (*models.Booking, error) {

	var booking models.Booking
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&booking).Error

	return &booking, err
}

//---------------- UPDATE ----------------

func (r *bookingRepository) Update(booking *models.Booking) error {
	return config.DB.Save(booking).Error
}

//---------------- HARD DELETE (ADMIN ONLY) ----------------

func (r *bookingRepository) Delete(id uint) error {
	return config.DB.
		Unscoped().
		Delete(&models.Booking{}, id).
		Error
}
