package dto

import (
	"portfolio/domain/entities"
)

// @Description Response for a successful authentication
type AuthSuccess struct {
	Token     string      `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	TokenType string      `json:"token_type,omitempty" example:"Bearer"`                    // e.g., "Bearer"
	ExpiresAt string      `json:"expires_at,omitempty" example:"2025-08-18T22:13:57+02:00"` // ISO 8601 format
	IssuedAt  string      `json:"issued_at,omitempty" example:"2025-08-17T22:13:57+02:00"`  // ISO 8601 format
	User      *UserPublic `json:"user,omitempty"`
	ExpiresIn int         `json:"expires_in,omitempty" example:"86400"` // Duration in seconds
} //@name ResponseAuthSuccess

// @Description Representation of a user
type UserPublic struct {
	LastLogin string `json:"last_login,omitempty" example:"2025-08-17T20:13:35+02:00"`
	Email     string `json:"email" example:"admin@localhost"`
	Role      string `json:"role" example:"Administrator"`
	Username  string `json:"username" example:"admin"`
	CreatedAt string `json:"created_at" example:"2025-08-09T04:18:24+02:00"`
	ID        int    `json:"id" example:"1"`
	IsActive  bool   `json:"is_active" example:"true"`
} //@name User

func FromUserEntityToUserPublic(user *entities.User) *UserPublic {
	if user == nil {
		return nil
	}
	return &UserPublic{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.DisplayName(),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastLogin: user.LastLogin.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func NewAuthResponse(user *entities.User, token string, expiresAt string, expiresIn int, issuedAt string) *AuthSuccess {
	return &AuthSuccess{
		User:      FromUserEntityToUserPublic(user),
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		ExpiresIn: expiresIn,
		IssuedAt:  issuedAt,
	}
}
