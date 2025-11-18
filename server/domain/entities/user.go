package entities

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLogin time.Time
	Username  string
	Password  string
	Email     string
	Role      UserRole
	ID        int
	IsActive  bool
}

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}

func (r UserRole) String() string {
	return string(r)
}

func (r UserRole) DisplayName() string {
	switch r {
	case RoleAdmin:
		return "Administrator"
	case RoleUser:
		return "User"
	default:
		return "Unknown"
	}
}

func ParseUserRole(role string) (UserRole, error) {
	userRole := UserRole(strings.ToLower(strings.TrimSpace(role)))
	if !userRole.IsValid() {
		return "", errors.New("invalid user role")
	}
	return userRole, nil
}

func GetAllValidRoles() []UserRole {
	return []UserRole{RoleAdmin, RoleUser}
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) CanLogin() bool {
	return u.IsActive
}

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *User) GeneratePasswordResetToken() (string, error) {
	bytes := make([]byte, 32)
	n, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}
	if n != len(bytes) {
		return "", fmt.Errorf("short random read: got %d want %d", n, len(bytes))
	}
	return hex.EncodeToString(bytes), nil
}

func (u *User) ChangeRole(newRole UserRole) error {
	if !newRole.IsValid() {
		return errors.New("invalid user role")
	}
	u.Role = newRole
	u.UpdatedAt = time.Now()
	return nil
}
