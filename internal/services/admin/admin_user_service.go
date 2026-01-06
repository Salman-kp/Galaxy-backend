package admin

import (
	"errors"
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/utils"

	"gorm.io/gorm"
)

type AdminUserService struct {
	repo         interfaces.UserRepository
	roleWageRepo interfaces.RoleWageRepository
}

func NewAdminUserService(
	repo interfaces.UserRepository,
	wagesRepo interfaces.RoleWageRepository,
) *AdminUserService {
	return &AdminUserService{
		repo:         repo,
		roleWageRepo: wagesRepo,
	}
}

// ---------------- CREATE USER ----------------
func (s *AdminUserService) CreateUser(input *models.User) error {
	existing, err := s.repo.FindByPhone(input.Phone)
	if err == nil && existing.ID != 0 && !existing.DeletedAt.Valid {
		return errors.New("phone already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if rw, err := s.roleWageRepo.FindByRole(input.Role); err == nil {
		input.CurrentWage = rw.Wage
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	input.Password = hashed
	input.JoinedAt = time.Now()

	return s.repo.Create(input)
}

// ---------------- LIST USERS ----------------
func (s *AdminUserService) ListUsers(role string, status string) ([]models.User, error) {
	return s.repo.ListAll(role, status)
}

// ---------------- GET USER ----------------
func (s *AdminUserService) GetUser(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

// ---------------- UPDATE USER ----------------
func (s *AdminUserService) UpdateUser(input *models.User) error {
	old, err := s.repo.FindByID(input.ID)
	if err != nil {
		return err
	}

	changed := false

	if input.Name != "" && input.Name != old.Name {
		old.Name = input.Name
		changed = true
	}

	if input.Phone != "" && input.Phone != old.Phone {
		existing, err := s.repo.FindByPhone(input.Phone)
		if err == nil && existing.ID != old.ID && !existing.DeletedAt.Valid {
			return errors.New("phone already exists")
		}
		old.Phone = input.Phone
		changed = true
	}

	if input.Role != "" && input.Role != old.Role {
		if !models.ValidateRole(input.Role) {
			return errors.New("invalid role")
		}
		old.Role = input.Role

		// 🔧 UPDATED: fetch latest wage on role change
		if rw, err := s.roleWageRepo.FindByRole(input.Role); err == nil {
			old.CurrentWage = rw.Wage
		}

		changed = true
	}

	if input.Branch != "" && input.Branch != old.Branch {
		old.Branch = input.Branch
		changed = true
	}

	if input.StartingPoint != "" && input.StartingPoint != old.StartingPoint {
		old.StartingPoint = input.StartingPoint
		changed = true
	}

	if input.BloodGroup != "" && input.BloodGroup != old.BloodGroup {
		old.BloodGroup = input.BloodGroup
		changed = true
	}

	if input.Status != "" && input.Status != old.Status {
		if !models.ValidateStatus(input.Status) {
			return errors.New("invalid status")
		}
		old.Status = input.Status
		changed = true
	}

	if input.DOB != nil && (old.DOB == nil || !old.DOB.Equal(*input.DOB)) {
		old.DOB = input.DOB
		changed = true
	}

	if input.Photo != "" && input.Photo != old.Photo {
		old.Photo = input.Photo
		changed = true
	}

	if !changed {
		return errors.New("no changes detected")
	}

	return s.repo.Update(old)
}

// ---------------- BLOCK USER ----------------
func (s *AdminUserService) BlockUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if user.Status == models.StatusBlocked {
		return errors.New("user already blocked")
	}

	user.Status = models.StatusBlocked
	return s.repo.Update(user)
}

// ---------------- UNBLOCK USER ----------------
func (s *AdminUserService) UnblockUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if user.Status == models.StatusActive {
		return errors.New("user already active")
	}

	user.Status = models.StatusActive
	return s.repo.Update(user)
}

// ---------------- SOFT DELETE USER ----------------
func (s *AdminUserService) SoftDeleteUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if user.DeletedAt.Valid {
		return errors.New("user already deleted")
	}
	return s.repo.SoftDelete(id)
}

// ---------------- RESET PASSWORD ----------------
func (s *AdminUserService) ResetPassword(id uint, newPassword string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if utils.CheckPasswordHash(newPassword, user.Password) {
		return errors.New("new password cannot be same as old password")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashed
	return s.repo.Update(user)
}

// ---------------- FILTER BY ROLE ----------------
func (s *AdminUserService) ListUsersByRole(role string) ([]models.User, error) {
	if !models.ValidateRole(role) {
		return []models.User{}, errors.New("invalid role")
	}

	var users []models.User
	err := config.DB.
		Where("role = ? AND deleted_at IS NULL", role).
		Order("created_at DESC").
		Find(&users).Error

	return users, err
}

// ---------------- SEARCH BY PHONE ----------------
func (s *AdminUserService) SearchUsersByPhone(phone string) ([]models.User, error) {
	var users []models.User
	err := config.DB.
		Where("phone ILIKE ? AND deleted_at IS NULL", "%"+phone+"%").
		Order("created_at DESC").
		Find(&users).Error

	return users, err
}