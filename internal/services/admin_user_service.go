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

func (s *AdminUserService) UpdateUser(user *models.User) error {
	old, err := s.repo.FindByID(user.ID)
	if err != nil {
		return err
	}

	if old.Name == user.Name &&
		old.Phone == user.Phone &&
		old.Role == user.Role &&
		old.Branch == user.Branch &&
		old.StartingPoint == user.StartingPoint &&
		old.BloodGroup == user.BloodGroup &&
		old.Status == user.Status &&
		(old.DOB == nil || user.DOB == nil || old.DOB.Equal(*user.DOB)) &&
		user.Photo == "" {
		return errors.New("no changes detected")
	}

	if user.Phone != old.Phone {
		existing, err := s.repo.FindByPhone(user.Phone)
		if err == nil && existing.ID != 0 {
			return errors.New("phone already exists")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if user.Role != old.Role {
		if wage, ok := s.roleWages[user.Role]; ok {
			user.CurrentWage = wage
		}
	}

	user.CreatedAt = old.CreatedAt
	return s.repo.Update(user)
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
