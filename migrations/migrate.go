package migrations

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/models"
)

func Migrate() error {
	return config.DB.AutoMigrate(
		&models.User{},
		&models.RoleWage{},
		&models.RefreshToken{},
	)
}
