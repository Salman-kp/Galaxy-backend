package admin

import (
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
)

type RoleService struct {
	roleRepo interfaces.RoleRepository
	permRepo interfaces.PermissionRepository
}

func NewRoleService(r interfaces.RoleRepository, p interfaces.PermissionRepository) *RoleService {
	return &RoleService{roleRepo: r, permRepo: p}
}

func (s *RoleService) CreatePermission(slug, desc string) error {
	perm := &models.Permission{Slug: slug, Description: desc}
	return s.permRepo.CreatePermission(perm)
}

func (s *RoleService) ListPermissions() ([]models.Permission, error) {
	return s.permRepo.FindAllPermissions()
}

func (s *RoleService) CreateRole(name string, permIDs []uint) error {
	perms, err := s.permRepo.FindPermissionsByIDs(permIDs)
	if err != nil { return err }

	role := &models.AdminRole{
		Name:        name,
		Permissions: perms,
	}
	return s.roleRepo.CreateRole(role)
}

func (s *RoleService) ListRoles() ([]models.AdminRole, error) {
	return s.roleRepo.FindAllRoles()
}

func (s *RoleService) GetRoleDetails(id uint) (*models.AdminRole, error) {
	return s.roleRepo.FindRoleByID(id)
}

func (s *RoleService) UpdateRole(id uint, name string, permIDs []uint) error {
    role, err := s.roleRepo.FindRoleByID(id)
    if err != nil {
        return err
    }

    perms, err := s.permRepo.FindPermissionsByIDs(permIDs)
    if err != nil {
        return err
    }

    role.Name = name
    role.Permissions = perms

    return s.roleRepo.UpdateRole(role)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.roleRepo.DeleteRole(id)
}