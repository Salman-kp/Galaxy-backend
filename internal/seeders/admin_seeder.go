package seeders

import (
	"time"

	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/utils"

	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var count int64

	db.Model(&models.User{}).
		Where("role = ?", models.RoleAdmin).
		Count(&count)

	if count > 0 {
		return
	}

	hash, err := utils.HashPassword("Admin@123")
	if err != nil {
		panic("password hash failed")
	}

	dob := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	admin := models.User{
		Name:          "Admin",
		Phone:         "9876543210",
		Password:      hash,
		Role:          models.RoleAdmin,
		Branch:        "Head Office",
		StartingPoint: "Main Branch",
		BloodGroup:    "O+",
		DOB:           &dob,
		Status:        models.StatusActive,
	}

	if err := db.Create(&admin).Error; err != nil {
		panic(err)
	}
}