package services

import (
	"errors"
	"time"

	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/utils"

	"gorm.io/gorm"
)

type AdminUserService struct {
	repo      interfaces.UserRepository
	roleWages map[string]int64
}

func NewAdminUserService(repo interfaces.UserRepository, wagesRepo interfaces.RoleWageRepository) *AdminUserService {
	wages := make(map[string]int64)
	allWages, err := wagesRepo.GetAll()
	if err != nil {
		allWages = []models.RoleWage{}
	}
	for _, w := range allWages {
		wages[w.Role] = w.Wage
	}
	return &AdminUserService{repo: repo, roleWages: wages}
}

func (s *AdminUserService) CreateUser(input *models.User) error {
	existing, err := s.repo.FindByPhone(input.Phone)
	if err == nil && existing.ID != 0 {
		return errors.New("phone already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if wage, ok := s.roleWages[input.Role]; ok {
		input.CurrentWage = wage
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	input.Password = hashed
	input.JoinedAt = time.Now()

	return s.repo.Create(input)
}

func (s *AdminUserService) ListUsers(role string, status string) ([]models.User, error) {
	return s.repo.ListAll(role, status)
}

func (s *AdminUserService) GetUser(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

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
		if err == nil && existing.ID != old.ID {
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
		if wage, ok := s.roleWages[input.Role]; ok {
			old.CurrentWage = wage
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

	// 🔥 THIS IS THE IMPORTANT PART
	if !changed {
		return errors.New("no changes detected")
	}

	return s.repo.Update(old)
}


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
