package repository

import (
	"errors"
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
)

type roleWageRepository struct{}

func NewRoleWageRepository() interfaces.RoleWageRepository {
	return &roleWageRepository{}
}

func (r *roleWageRepository) Create(wage *models.RoleWage) error {
	return config.DB.Create(wage).Error
}

func (r *roleWageRepository) Update(wage *models.RoleWage) error {
	res := config.DB.Save(wage)
	if res.RowsAffected == 0 {
		return errors.New("no rows updated")
	}
	return res.Error
}

func (r *roleWageRepository) Delete(role string) error {
	return config.DB.Where("role = ?", role).Delete(&models.RoleWage{}).Error
}

func (r *roleWageRepository) GetAll() ([]models.RoleWage, error) {
	var wages []models.RoleWage
	err := config.DB.Find(&wages).Error
	return wages, err
}

func (r *roleWageRepository) FindByRole(role string) (*models.RoleWage, error) {
	var wage models.RoleWage
	err := config.DB.Where("role = ?", role).First(&wage).Error
	return &wage, err
}