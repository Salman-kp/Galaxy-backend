package repository

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
)

type refreshTokenRepository struct{}

func NewRefreshTokenRepository() interfaces.RefreshTokenRepository {
	return &refreshTokenRepository{}
}

func (r *refreshTokenRepository) Save(t *models.RefreshToken) error {
	err := config.DB.Where("user_id = ?", t.UserID).Delete(&models.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return config.DB.Create(t).Error
}

func (r *refreshTokenRepository) FindByUserID(userID uint) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := config.DB.Where("user_id = ?", userID).First(&token).Error
	return &token, err
}

func (r *refreshTokenRepository) FindByHashedToken(hashed string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := config.DB.Where("token_hashed = ?", hashed).First(&token).Error
	return &token, err
}

func (r *refreshTokenRepository) DeleteByUserID(userID uint) error {
	return config.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
