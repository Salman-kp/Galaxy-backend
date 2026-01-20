package seeders

import (
	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedRBAC(db *gorm.DB) {
	// 1. Define ALL granular permissions exactly as they appear in AdminRoutes
	permissions := []models.Permission{
		// --- USER MANAGEMENT ---
		{Slug: "user:create", Description: "Ability to create or invite new users"},
		{Slug: "user:view", Description: "Ability to view user lists, search, and details"},
		{Slug: "user:edit", Description: "Ability to update user information"},
		{Slug: "user:status", Description: "Ability to block or unblock users"},
		{Slug: "user:delete", Description: "Ability to delete user accounts"},
		{Slug: "user:password", Description: "Ability to reset user passwords"},

		// --- EVENT MANAGEMENT ---
		{Slug: "event:view", Description: "View event details, lists, and bookings"},
		{Slug: "event:create", Description: "Create new event entries"},
		{Slug: "event:edit", Description: "Update existing event information"},
		{Slug: "event:delete", Description: "Delete event entries"},
		{Slug: "event:operate", Description: "Operational access: Start, Complete, Cancel events and Attendance"},

		// --- WAGES & FINANCE ---
		{Slug: "managewages:view", Description: "Update global standard role-based wages"},
		{Slug: "wage:view", Description: "View event-specific wage summaries and reports"},
		{Slug: "wage:edit", Description: "Override individual worker wages for specific bookings"},

		// --- DASHBOARD & PROFILE ---
		{Slug: "dashboard:view", Description: "Access to view dashboard statistics and charts"},
		{Slug: "profile:edit", Description: "Ability to edit personal admin profile information"},

		// --- RBAC MANAGEMENT ---
		{Slug: "rbac:view", Description: "Full control over roles, permissions, and admin management"},
	}

	// 2. Insert or Update Permissions
	for _, p := range permissions {
		db.Where(models.Permission{Slug: p.Slug}).
			Assign(models.Permission{Description: p.Description}).
			FirstOrCreate(&p)
	}

	// 3. Fetch all newly created/existing permissions
	var allPerms []models.Permission
	db.Find(&allPerms)

	// 4. Create or Update the "Super Admin" Role
	var superAdminRole models.AdminRole
	if err := db.Where("name = ?", "Super Admin").First(&superAdminRole).Error; err != nil {
		superAdminRole = models.AdminRole{
			Name:        "Super Admin",
			Permissions: allPerms,
		}
		db.Create(&superAdminRole)
	} else {
		db.Model(&superAdminRole).Association("Permissions").Replace(allPerms)
	}

	// 5. Create or Sync the Primary System Administrator
	hashedPassword, _ := utils.HashPassword("admin123") 
	dob := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	firstAdmin := models.User{
		Name:          "System Administrator",
		Phone:         "1234567890", 
		Password:      hashedPassword,
		Role:          models.RoleAdmin, 
		Branch:        "Head Office",    
		StartingPoint: "Main Branch",    
		BloodGroup:    "O+",             
		DOB:           &dob,             
		Status:        models.StatusActive,
		AdminRoleID:   &superAdminRole.ID, 
		JoinedAt:      time.Now(),
	}

	var existingUser models.User
	if err := db.Where("phone = ?", firstAdmin.Phone).First(&existingUser).Error; err != nil {
		if err := db.Create(&firstAdmin).Error; err != nil {
			fmt.Printf("❌ Failed to create admin: %v\n", err)
		} else {
			fmt.Println("✅ Database seeded: Super Admin created (1234567890 / admin123)")
		}
	} else {
		db.Model(&existingUser).Updates(map[string]interface{}{
			"admin_role_id": superAdminRole.ID,
			"role":          models.RoleAdmin, 
			"status":        models.StatusActive,
		})
		fmt.Println("ℹ️ Seeding synced: User 1234567890 is now a Super Admin with full permissions.")
	}
}