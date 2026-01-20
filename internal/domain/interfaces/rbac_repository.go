package interfaces

import "event-management-backend/internal/domain/models"

type RoleRepository interface {
    CreateRole(role *models.AdminRole) error
    FindAllRoles() ([]models.AdminRole, error)
    FindRoleByID(id uint) (*models.AdminRole, error)
    UpdateRole(role *models.AdminRole) error
    DeleteRole(id uint) error
}

type PermissionRepository interface {
    CreatePermission(perm *models.Permission) error
    FindAllPermissions() ([]models.Permission, error)
    FindPermissionsByIDs(ids []uint) ([]models.Permission, error)
}