package validations

import (
	"errors"
	"regexp"
)

var PhoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *LoginRequest) Validate() error {
	if !PhoneRegex.MatchString(r.Phone) {
		return errors.New("invalid phone number")
	}
	if len(r.Password) < 4 {
		return errors.New("password must be at least 4 characters")
	}
	return nil
}