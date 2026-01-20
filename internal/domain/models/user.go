package models

import (
	"errors"
	"regexp"
	"time"

	"gorm.io/gorm"
)

const (
	RoleAdmin      = "admin"
	RoleCaptain    = "captain"
	RoleSubCaptain = "sub_captain"
	RoleMainBoy    = "main_boy"
	RoleJuniorBoy  = "junior_boy"
	
	StatusActive   = "active"
	StatusBlocked  = "blocked"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:150;not null" json:"name"`
	Phone         string         `gorm:"size:30;uniqueIndex;not null" json:"phone"`
	Password      string         `gorm:"not null" json:"-"`
	Role          string         `gorm:"size:50;not null" json:"role"`
	Branch        string         `gorm:"size:100" json:"branch"`
	StartingPoint string         `gorm:"size:150" json:"starting_point"`
	BloodGroup    string         `gorm:"size:10" json:"blood_group"`
	DOB           *time.Time     `json:"dob"`
	Photo         string         `gorm:"size:255" json:"photo"`
	JoinedAt      time.Time      `gorm:"autoCreateTime" json:"joined_at"`
	CompletedWork uint           `gorm:"default:0" json:"completed_work"`
	CurrentWage   int64          `gorm:"default:0" json:"current_wage"`
	Status        string         `gorm:"size:30;default:'active'" json:"status"`
	AdminRoleID   *uint          `json:"admin_role_id"`
    AdminRole     *AdminRole     `gorm:"foreignKey:AdminRoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"admin_role"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

var (
	nameRegex  = regexp.MustCompile(`^[A-Za-z ]{3,}$`)
	phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
)

func ValidateRole(r string) bool {
	switch r {
	case RoleAdmin, RoleCaptain, RoleSubCaptain, RoleMainBoy, RoleJuniorBoy:
		return true
	}
	return false
}

func ValidateStatus(s string) bool {
	switch s {
	case StatusActive, StatusBlocked:
		return true
	}
	return false
}
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.validateFields()
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Name != "" && !nameRegex.MatchString(u.Name) {
		return errors.New("invalid name")
	}
	if u.Phone != "" && !phoneRegex.MatchString(u.Phone) {
		return errors.New("invalid phone number")
	}
	if u.Role != "" && !ValidateRole(u.Role) {
		return errors.New("invalid role")
	}
	if u.Status != "" && !ValidateStatus(u.Status) {
		return errors.New("invalid status")
	}
	return nil
}


func (u *User) validateFields() error {
	if !nameRegex.MatchString(u.Name) {
		return errors.New("invalid name")
	}
	if !phoneRegex.MatchString(u.Phone) {
		return errors.New("invalid phone number")
	}
	if !ValidateRole(u.Role) {
		return errors.New("invalid role")
	}
	if u.Status == "" {
		u.Status = StatusActive
	} else if !ValidateStatus(u.Status) {
		return errors.New("invalid status")
	}
	if u.JoinedAt.IsZero() {
		u.JoinedAt = time.Now()
	}
	return nil
}