package repository

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"

)
type roleRepository struct{}
type permissionRepository struct{}

func NewRoleRepository() interfaces.RoleRepository             { return &roleRepository{} }
func NewPermissionRepository() interfaces.PermissionRepository { return &permissionRepository{} }

// Role Methods
func (r *roleRepository) CreateRole(role *models.AdminRole) error {
	return config.DB.Create(role).Error
}

func (r *roleRepository) FindAllRoles() ([]models.AdminRole, error) {
	var roles []models.AdminRole
	err := config.DB.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindRoleByID(id uint) (*models.AdminRole, error) {
	var role models.AdminRole
	err := config.DB.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (r *roleRepository) UpdateRole(role *models.AdminRole) error {
    if err := config.DB.Model(role).Updates(models.AdminRole{Name: role.Name}).Error; err != nil {
        return err
    }
    return config.DB.Model(role).Association("Permissions").Replace(role.Permissions)
}

func (r *roleRepository) DeleteRole(id uint) error {
	return config.DB.Delete(&models.AdminRole{}, id).Error
}

// Permission Methods
func (r *permissionRepository) CreatePermission(perm *models.Permission) error {
	return config.DB.Create(perm).Error
}

func (r *permissionRepository) FindAllPermissions() ([]models.Permission, error) {
	var perms []models.Permission
	err := config.DB.Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindPermissionsByIDs(ids []uint) ([]models.Permission, error) {
	var perms []models.Permission
	err := config.DB.Where("id IN ?", ids).Find(&perms).Error
	return perms, err
}