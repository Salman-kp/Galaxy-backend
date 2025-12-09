package repository

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
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
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) ListAll(role string, status string) ([]models.User, error) {
	var users []models.User
	query := config.DB.Model(&models.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&users).Error
	return users, err
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := config.DB.Find(&users).Error
	return users, err
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Save(user).Error
}

func (r *userRepository) SoftDelete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}