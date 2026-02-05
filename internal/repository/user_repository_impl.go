package repository

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

type userRepository struct{}

func NewUserRepository() interfaces.UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.
	    Preload("AdminRole.Permissions").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user).Error
	return &user, err
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := config.DB.
	    Preload("AdminRole.Permissions").
		Where("phone = ? AND deleted_at IS NULL", phone).
		First(&user).Error
	return &user, err
}

func (r *userRepository) ListAll(role string, status string) ([]models.User, error) {
	var users []models.User
	query := config.DB.Model(&models.User{}).
             Preload("AdminRole"). 
             Where("deleted_at IS NULL") 

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at ASC").Find(&users).Error
	return users, err
}
func (r *userRepository) ListByRole(role string) ([]models.User, error) {
    var users []models.User
    err := config.DB.
        Preload("AdminRole"). 
        Where("role = ? AND deleted_at IS NULL", role).
        Order("created_at DESC").
        Find(&users).Error
    return users, err
}
func (r *userRepository) SearchByPhone(phone string) ([]models.User, error) {
    var users []models.User
    err := config.DB.
        Preload("AdminRole"). 
        Where("phone ILIKE ? AND deleted_at IS NULL", "%"+phone+"%").
        Order("created_at DESC").
        Find(&users).Error
    return users, err
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := config.DB.
		Where("deleted_at IS NULL").
		Find(&users).Error
	return users, err
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := config.DB.
		Model(&models.User{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Save(user).Error
}
func (r *userRepository) UpdateRole(user *models.User) error {
 return config.DB.Model(user).Select("Role", "AdminRoleID", "CurrentWage", "UpdatedAt").Updates(user).Error
}
func (r *userRepository) RemovePhoto(id uint) error {
    return config.DB.Model(&models.User{}).Where("id = ?", id).Update("photo", "").Error
}
func (r *userRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return config.DB.
		Session(&gorm.Session{SkipHooks: true}).
		Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *userRepository) UpdateWageByRole(role string, wage int64) error {
	return config.DB.
		Model(&models.User{}).
		Where("role = ? AND deleted_at IS NULL", role).
		Update("current_wage", wage).Error
}

func (r *userRepository) SoftDelete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}