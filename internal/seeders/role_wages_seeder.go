package seeders

import (
	"event-management-backend/internal/domain/models"

	"gorm.io/gorm"
)

func SeedRoleWages(db *gorm.DB) {
	wages := []models.RoleWage{
		{Role: models.RoleCaptain, Wage: 1200},
		{Role: models.RoleSubCaptain, Wage: 900},
		{Role: models.RoleMainBoy, Wage: 700},
		{Role: models.RoleJuniorBoy, Wage: 500},
	}

	for _, w := range wages {
		var count int64

		db.Model(&models.RoleWage{}).
			Where("role = ?", w.Role).
			Count(&count)

		if count > 0 {
			continue
		}

		if err := db.Create(&w).Error; err != nil {
			panic(err)
		}
	}
}