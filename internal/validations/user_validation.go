package validations

import (
	"errors"
	"event-management-backend/internal/domain/models"
	"regexp"
)

var (
	nameRegex  = regexp.MustCompile(`^[A-Za-z ]{3,}$`)
	phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
)

//
// ---------------- CREATE USER ----------------
//
type CreateUserRequest struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	Branch        string `json:"branch"`
	StartingPoint string `json:"starting_point"`
	BloodGroup    string `json:"blood_group"`
	DOB           string `json:"dob"` // YYYY-MM-DD
	AdminRoleID *uint `json:"admin_role_id"`
}

func (r *CreateUserRequest) Validate() error {
    if r.Name == "" || !nameRegex.MatchString(r.Name) {
        return errors.New("name is required and must be valid")
    }
    if r.Phone == "" || !phoneRegex.MatchString(r.Phone) {
        return errors.New("10-digit phone number is required")
    }
    if r.Role == "" || !models.ValidateRole(r.Role) {
        return errors.New("a valid role is required")
    }
    if r.Role == models.RoleAdmin && r.AdminRoleID == nil {
		return errors.New("admin role ID is required for administrative accounts")
	}
	if len(r.Password) < 4 {
		return errors.New("password must be at least 4 characters")
	}
	return nil
}

//
// ---------------- UPDATE USER (ADMIN) ----------------
//
type UpdateUserRequest struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Role          string `json:"role"`
	Branch        string `json:"branch"`
	StartingPoint string `json:"starting_point"`
	BloodGroup    string `json:"blood_group"`
	DOB           string `json:"dob"` // YYYY-MM-DD
	Status        string `json:"status"`
	AdminRoleID *uint `json:"admin_role_id"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Name != "" && !nameRegex.MatchString(r.Name) {
		return errors.New("invalid name")
	}

	if r.Phone != "" && !phoneRegex.MatchString(r.Phone) {
		return errors.New("invalid phone")
	}

	if r.Role != "" && !models.ValidateRole(r.Role) {
		return errors.New("invalid role")
	}
    if r.Role == models.RoleAdmin && r.AdminRoleID == nil {
		return errors.New("admin role ID is required when setting role to admin")
	}
	if r.Status != "" && !models.ValidateStatus(r.Status) {
		return errors.New("invalid status")
	}
	return nil
}

//
// ---------------- ADMIN SELF PROFILE UPDATE ----------------
//
type UpdateAdminSelfProfileRequest struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Branch        string `json:"branch"`
	StartingPoint string `json:"starting_point"`
	BloodGroup    string `json:"blood_group"`
	DOB           string `json:"dob"`
}

func (r *UpdateAdminSelfProfileRequest) Validate() error {
	if r.Name != "" && !nameRegex.MatchString(r.Name) {
		return errors.New("invalid name")
	}

	if r.Phone != "" && !phoneRegex.MatchString(r.Phone) {
		return errors.New("invalid phone")
	}
	return nil
}