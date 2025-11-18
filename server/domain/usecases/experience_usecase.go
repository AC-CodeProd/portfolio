package usecases

import (
	"context"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
)

type ExperienceUseCase struct {
	experienceRepo interfaces.ExperienceRepository
	userRepo       interfaces.UserRepository
	logger         *logger.Logger
}

func NewExperienceUseCase(experienceRepo interfaces.ExperienceRepository, userRepo interfaces.UserRepository, logger *logger.Logger) *ExperienceUseCase {
	return &ExperienceUseCase{
		experienceRepo: experienceRepo,
		userRepo:       userRepo,
		logger:         logger,
	}
}

func (uc *ExperienceUseCase) CreateExperience(ctx context.Context, experience *entities.Experience) (*entities.Experience, error) {
	if !experience.HasRequiredFields() {
		uc.logger.Error("Invalid experience fields: %v", experience)
		return nil, domain.NewValidationError("experience", "job title, company name, and user ID are required", nil)
	}

	userExists, err := uc.userRepo.ExistsByID(ctx, experience.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if user exists: %v", err)
		return nil, domain.NewInternalError("failed to validate user", err)
	}
	if !userExists {
		uc.logger.Error("User not found for ID %d", experience.UserID)
		return nil, domain.NewNotFoundError("User", fmt.Sprint(experience.UserID))
	}

	exists, err := uc.experienceRepo.ExistsByJobTitleCompanyAndUserID(ctx, experience.JobTitle, experience.CompanyName, experience.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if experience exists: %v", err)
		return nil, domain.NewInternalError("failed to check experience existence", err)
	}
	if exists {
		uc.logger.Error("Experience already exists: %s at %s", experience.JobTitle, experience.CompanyName)
		return nil, domain.NewAlreadyExistsError("Experience", experience.JobTitle+" at "+experience.CompanyName)
	}

	createdExperience, err := uc.experienceRepo.Create(ctx, experience)
	if err != nil {
		uc.logger.Error("Failed to create experience: %v", err)
		return nil, domain.NewInternalError("failed to create experience", err)
	}

	return createdExperience, nil
}

func (uc *ExperienceUseCase) GetExperienceByID(ctx context.Context, experienceID int) (*entities.Experience, error) {
	if experienceID <= 0 {
		uc.logger.Error("Invalid experience ID: %d", experienceID)
		return nil, domain.NewValidationError("experienceID", "experience ID must be positive", nil)
	}

	experience, err := uc.experienceRepo.GetByID(ctx, experienceID)
	if err != nil {
		uc.logger.Error("Failed to get experience by ID %d: %v", experienceID, err)
		return nil, domain.NewInternalError("failed to retrieve experience", err)
	}

	if experience == nil {
		uc.logger.Error("Experience not found for ID %d", experienceID)
		return nil, domain.NewNotFoundError("Experience", fmt.Sprint(experienceID))
	}

	return experience, nil
}

func (uc *ExperienceUseCase) GetExperiencesByUserID(ctx context.Context, userID int) ([]*entities.Experience, error) {
	if userID <= 0 {
		uc.logger.Error("Invalid user ID: %d", userID)
		return nil, domain.NewValidationError("userID", "user ID must be positive", nil)
	}

	experiences, err := uc.experienceRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get experiences for user %d: %v", userID, err)
		return nil, domain.NewInternalError("failed to retrieve experiences", err)
	}

	return experiences, nil
}

func (uc *ExperienceUseCase) UpdateExperience(ctx context.Context, experienceID int, experience *entities.Experience) (*entities.Experience, error) {
	if !experience.HasRequiredFields() {
		uc.logger.Error("Invalid experience fields: %v", experience)
		return nil, domain.NewValidationError("experience", "job title, company name, and user ID are required", nil)
	}

	exists, err := uc.experienceRepo.ExistsByID(ctx, experience.ExperienceID)
	if err != nil {
		uc.logger.Error("Failed to check if experience exists: %v", err)
		return nil, domain.NewInternalError("failed to check experience existence", err)
	}
	if !exists {
		uc.logger.Error("Experience not found for ID %d", experience.ExperienceID)
		return nil, domain.NewNotFoundError("Experience", fmt.Sprint(experience.ExperienceID))
	}

	experience.MarkAsUpdated()
	updatedExperience, err := uc.experienceRepo.Update(ctx, experienceID, experience)
	if err != nil {
		uc.logger.Error("Failed to update experience: %v", err)
		return nil, domain.NewInternalError("failed to update experience", err)
	}

	return updatedExperience, nil
}

func (uc *ExperienceUseCase) PatchExperience(ctx context.Context, experienceID int, updates *entities.Experience) (*entities.Experience, error) {
	if experienceID <= 0 {
		uc.logger.Error("Experience ID is required")
		return nil, domain.NewValidationError("experienceID", "experience ID must be positive", nil)
	}

	existingExperience, err := uc.experienceRepo.GetByID(ctx, experienceID)
	if err != nil {
		uc.logger.Error("Failed to get experience by ID %d: %v", experienceID, err)
		return nil, domain.NewInternalError("failed to get experience", err)
	}

	if existingExperience == nil {
		uc.logger.Error("Experience not found: %d", experienceID)
		return nil, domain.NewNotFoundError("Experience", fmt.Sprint(experienceID))
	}

	if updates.JobTitle != "" {
		existingExperience.JobTitle = updates.JobTitle
	}
	if updates.CompanyName != "" {
		existingExperience.CompanyName = updates.CompanyName
	}
	if !updates.StartDate.IsZero() {
		existingExperience.StartDate = updates.StartDate
	}
	if updates.EndDate.IsZero() {
		existingExperience.EndDate = updates.EndDate
	}
	if updates.Description != "" {
		existingExperience.Description = updates.Description
	}

	existingExperience.MarkAsUpdated()

	patchedExperience, err := uc.experienceRepo.Patch(ctx, experienceID, existingExperience)
	if err != nil {
		uc.logger.Error("Failed to update experience: %v", err)
		return nil, domain.NewInternalError("failed to update experience", err)
	}

	return patchedExperience, nil
}

func (uc *ExperienceUseCase) DeleteExperience(ctx context.Context, experienceID int) error {
	if experienceID <= 0 {
		uc.logger.Error("Invalid experience ID: %d", experienceID)
		return domain.NewValidationError("experienceID", "experience ID must be positive", nil)
	}

	exists, err := uc.experienceRepo.ExistsByID(ctx, experienceID)
	if err != nil {
		uc.logger.Error("Failed to check if experience exists: %v", err)
		return domain.NewInternalError("failed to check experience existence", err)
	}
	if !exists {
		uc.logger.Error("Experience not found for ID %d", experienceID)
		return domain.NewNotFoundError("Experience", fmt.Sprint(experienceID))
	}

	err = uc.experienceRepo.Delete(ctx, experienceID)
	if err != nil {
		uc.logger.Error("Failed to delete experience: %v", err)
		return domain.NewInternalError("failed to delete experience", err)
	}

	return nil
}

func (uc *ExperienceUseCase) GetAllExperiences(ctx context.Context) ([]*entities.Experience, error) {
	experiences, err := uc.experienceRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all experiences: %v", err)
		return nil, domain.NewInternalError("failed to retrieve experiences", err)
	}

	return experiences, nil
}

func (uc *ExperienceUseCase) GetCurrentExperiences(ctx context.Context, userID int) ([]*entities.Experience, error) {
	if userID <= 0 {
		uc.logger.Error("Invalid user ID: %d", userID)
		return nil, domain.NewValidationError("userID", "user ID must be positive", nil)
	}

	experiences, err := uc.experienceRepo.GetCurrentExperiences(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get current experiences for user %d: %v", userID, err)
		return nil, domain.NewInternalError("failed to retrieve current experiences", err)
	}

	return experiences, nil
}
