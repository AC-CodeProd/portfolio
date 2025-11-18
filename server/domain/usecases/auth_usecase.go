package usecases

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	dto "portfolio/dto/auth"
	"portfolio/helpers"
	"portfolio/logger"
	"portfolio/service"
	"time"
)

type AuthUseCase struct {
	salt            string
	userRepo        interfaces.UserRepository
	revokeTokenRepo interfaces.RevokedTokenRepository
	settingUseCase  *SettingUseCase
	authService     *service.AuthService
	logger          *logger.Logger
}

func NewAuthUseCase(userRepo interfaces.UserRepository, revokeTokenRepo interfaces.RevokedTokenRepository, settingUseCase *SettingUseCase, authService *service.AuthService, logger *logger.Logger, salt string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		revokeTokenRepo: revokeTokenRepo,
		settingUseCase:  settingUseCase,
		authService:     authService,
		logger:          logger,
		salt:            salt,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, request *dto.AuthRequest) (*dto.AuthSuccess, error) {
	if request == nil {
		return nil, domain.NewValidationError("Request cannot be nil", "request", nil)
	}

	request.Sanitize()
	if err := request.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.ValidateCredentials(ctx, request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		uc.logger.Error("User account is disabled for username %s", user.Username)
		return nil, domain.NewUnauthorizedError("User account is disabled")
	}

	token, err := uc.authService.GenerateToken(user.ID)
	if err != nil {
		uc.logger.Error("Failed to generate token for user %d: %v", user.ID, err)
		return nil, domain.NewInternalError("Failed to generate token", err)
	}

	err = uc.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		uc.logger.Warn("Failed to update last login for user %d: %v", user.ID, err)
	}

	return dto.NewAuthResponse(
		user,
		token.Token,
		token.ExpiresAt.Format(time.RFC3339),
		token.ExpiresIn,
		time.Now().Format(time.RFC3339),
	), nil
}

func (uc *AuthUseCase) ValidateCredentials(ctx context.Context, username, password string) (*entities.User, error) {
	if username == "" || password == "" {
		return nil, domain.NewValidationError("credentials", "username and password are required", nil)
	}

	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		uc.logger.Error("Failed to retrieve user by username %s: %v", username, err)
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	if user == nil {
		uc.logger.Error("User not found for username %s", username)
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	if !user.IsActive {
		uc.logger.Error("User account is disabled for username %s", username)
		return nil, domain.NewUnauthorizedError("user account is disabled")
	}

	if !uc.authService.CheckPassword(password, uc.salt, user.Password) {
		uc.logger.Error("Invalid password for user %s", username)
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	err = uc.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		uc.logger.Warn("Failed to update last login for user %d: %v", user.ID, err)
	}

	return user, nil
}

func (uc *AuthUseCase) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	return uc.userRepo.GetByUsername(ctx, username)
}

func (uc *AuthUseCase) HashPassword(ctx context.Context, password, salt string) (string, error) {
	return uc.authService.HashPassword(password, salt)
}

func (uc *AuthUseCase) CheckPassword(ctx context.Context, password, salt, hash string) bool {
	return uc.authService.CheckPassword(password, salt, hash)
}

func (uc *AuthUseCase) CreateDefaultAdmin(ctx context.Context, username string) error {

	if len(uc.salt) > 60 {
		uc.logger.Fatal("Admin salt too long: %d bytes (max 60 for bcrypt safety)", len(uc.salt))
		return fmt.Errorf("admin salt too long: %d bytes (max 60 for bcrypt safety)", len(uc.salt))
	}
	admin, _ := uc.GetUserByUsername(ctx, username)
	if admin != nil {
		uc.logger.Info("Default admin user %s already exists, skipping creation", username)
		return nil
	}

	password := helpers.RandomString(12)
	hashedPassword, err := uc.HashPassword(ctx, password, uc.salt)
	if err != nil {
		uc.logger.Error("Failed to hash password for admin user %s: %v", username, err)
		return err
	}

	err = uc.settingUseCase.Upsert(ctx, &entities.SettingJson{
		ShowProjects: true,
	})

	if err != nil {
		uc.logger.Error("Failed to create default settings for admin user %s: %v", username, err)
		return err
	}

	user := &entities.User{
		Username:  username,
		Password:  hashedPassword,
		Email:     username + "@localhost",
		Role:      entities.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		uc.logger.Error("Failed to create default admin user %s: %v", username, err)
		return err
	}

	uc.logger.Info("Default admin user created successfully: %s", username)

	if err := uc.writeAdminPasswordToFile(username, password); err != nil {
		uc.logger.Warn("Failed to write admin password to secure file: %v", err)
	}

	return nil
}

func (uc *AuthUseCase) RevokeToken(ctx context.Context, userID int, token string) error {
	return uc.revokeTokenRepo.RevokedToken(ctx, userID, token)
}

func (uc *AuthUseCase) IsTokenRevoked(ctx context.Context, userID int, token string) (bool, error) {
	return uc.revokeTokenRepo.IsTokenRevoked(ctx, userID, token)
}

func (uc *AuthUseCase) writeAdminPasswordToFile(username, password string) error {
	credDir := ".admin-credentials"
	if err := os.MkdirAll(credDir, 0700); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	credFile := filepath.Join(credDir, "admin-password.txt")

	content := "ADMIN CREDENTIALS - DELETE AFTER READING\n"
	content += "========================================\n\n"
	content += fmt.Sprintf("Username: %s\n", username)
	content += fmt.Sprintf("Password: %s\n\n", password)
	content += "⚠️  IMPORTANT: Delete this file immediately after noting the password!\n"
	content += "⚠️  This file contains sensitive information.\n\n"
	content += fmt.Sprintf("Command to delete: rm -f %s\n", credFile)

	if err := os.WriteFile(credFile, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	uc.logger.Info("Admin credentials written to secure file: %s (permissions: 0600)", credFile)
	uc.logger.Warn("⚠️  SECURITY: Please retrieve and delete the credentials file immediately!")

	return nil
}
