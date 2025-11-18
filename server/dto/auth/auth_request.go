package dto

import (
	"portfolio/domain"
	"strings"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"admin"`
	Password string `json:"password" validate:"required,min=6" example:"password"`
} //@name AuthRequest

func (lr *AuthRequest) Validate() error {
	if strings.TrimSpace(lr.Username) == "" {
		return domain.NewRequiredFieldError("username")
	}

	if strings.TrimSpace(lr.Password) == "" {
		return domain.NewRequiredFieldError("password")
	}

	if len(lr.Username) < 3 {
		return domain.NewValidationError("Username must be at least 3 characters long", "username", nil)
	}

	if len(lr.Password) < 6 {
		return domain.NewValidationError("Password must be at least 6 characters long", "password", nil)
	}

	lr.Sanitize()
	return nil
}

func (lr *AuthRequest) Sanitize() {
	lr.Username = strings.TrimSpace(strings.ToLower(lr.Username))
	lr.Password = strings.TrimSpace(lr.Password)
}
