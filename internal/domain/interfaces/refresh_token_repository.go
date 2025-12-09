package interfaces

import "event-management-backend/internal/domain/models"

type RefreshTokenRepository interface {
	Save(token *models.RefreshToken) error
	FindByUserID(userID uint) (*models.RefreshToken, error)
	FindByHashedToken(hashed string) (*models.RefreshToken, error)
	DeleteByUserID(userID uint) error
}
