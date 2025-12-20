package validations

import (
	"errors"
	"regexp"
)

var (
	nameRegex  = regexp.MustCompile(`^[A-Za-z ]{3,}$`)
	phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
)

type CreateUserRequest struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	Branch        string `json:"branch"`
	StartingPoint string `json:"starting_point"`
	BloodGroup    string `json:"blood_group"`
	DOB           string `json:"dob"`
}

func (r *CreateUserRequest) Validate() error {
	if !nameRegex.MatchString(r.Name) {
		return errors.New("invalid name")
	}
	if !phoneRegex.MatchString(r.Phone) {
		return errors.New("invalid phone")
	}
	if len(r.Password) < 4 {
		return errors.New("password must be at least 4 characters")
	}
	if r.Role == "" {
		return errors.New("role is required")
	}
	return nil
}

type UpdateUserRequest struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Role          string `json:"role"`
	Branch        string `json:"branch"`
	StartingPoint string `json:"starting_point"`
	BloodGroup    string `json:"blood_group"`
	DOB           string `json:"dob"`
	Status        string `json:"status"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Name != "" && !nameRegex.MatchString(r.Name) {
		return errors.New("invalid name")
	}
	if r.Phone != "" && !phoneRegex.MatchString(r.Phone) {
		return errors.New("invalid phone")
	}
	if r.Role == "" {
		return errors.New("role is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}