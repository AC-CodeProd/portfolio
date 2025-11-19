package usecases

import (
	"context"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
)

type EducationUseCase struct {
	educationRepo interfaces.EducationRepository
	userRepo      interfaces.UserRepository
	logger        *logger.Logger
}

func NewEducationUseCase(educationRepo interfaces.EducationRepository, userRepo interfaces.UserRepository, logger *logger.Logger) *EducationUseCase {
	return &EducationUseCase{
		educationRepo: educationRepo,
		userRepo:      userRepo,
		logger:        logger,
	}
}

func (uc *EducationUseCase) CreateEducation(ctx context.Context, education *entities.Education) (*entities.Education, error) {
	if !education.HasRequiredFields() {
		uc.logger.Error("Invalid education fields: %v", education)
		return nil, domain.NewValidationError("education", "degree, institution, and user ID are required", nil)
	}

	userExists, err := uc.userRepo.ExistsByID(ctx, education.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if user exists: %v", err)
		return nil, domain.NewInternalError("failed to validate user", err)
	}
	if !userExists {
		uc.logger.Error("User not found for ID %d", education.UserID)
		return nil, domain.NewNotFoundError("User", fmt.Sprint(education.UserID))
	}

	exists, err := uc.educationRepo.ExistsByDegreeInstitutionAndUserID(ctx, education.Degree, education.Institution, education.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if education exists: %v", err)
		return nil, domain.NewInternalError("failed to check education existence", err)
	}
	if exists {
		uc.logger.Error("Education already exists: %s at %s", education.Degree, education.Institution)
		return nil, domain.NewAlreadyExistsError("Education", education.Degree+" at "+education.Institution)
	}

	createdEducation, err := uc.educationRepo.Create(ctx, education)
	if err != nil {
		uc.logger.Error("Failed to create education: %v", err)
		return nil, domain.NewInternalError("failed to create education", err)
	}

	return createdEducation, nil
}

func (uc *EducationUseCase) GetEducationByID(ctx context.Context, educationID int) (*entities.Education, error) {
	if educationID <= 0 {
		uc.logger.Error("Invalid education ID: %d", educationID)
		return nil, domain.NewValidationError("educationID", "education ID must be positive", nil)
	}

	education, err := uc.educationRepo.GetByID(ctx, educationID)
	if err != nil {
		uc.logger.Error("Failed to get education by ID %d: %v", educationID, err)
		return nil, domain.NewInternalError("failed to retrieve education", err)
	}

	if education == nil {
		uc.logger.Error("Education not found for ID %d", educationID)
		return nil, domain.NewNotFoundError("Education", fmt.Sprint(educationID))
	}

	return education, nil
}

func (uc *EducationUseCase) GetEducationsByUserID(ctx context.Context, userID int) ([]*entities.Education, error) {
	if userID <= 0 {
		uc.logger.Error("Invalid user ID: %d", userID)
		return nil, domain.NewValidationError("userID", "user ID must be positive", nil)
	}

	educations, err := uc.educationRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get educations for user %d: %v", userID, err)
		return nil, domain.NewInternalError("failed to retrieve educations", err)
	}

	return educations, nil
}

func (uc *EducationUseCase) UpdateEducation(ctx context.Context, educationID int, education *entities.Education) (*entities.Education, error) {
	if !education.HasRequiredFields() {
		uc.logger.Error("Invalid education fields: %v", education)
		return nil, domain.NewValidationError("education", "degree, institution, and user ID are required", nil)
	}

	exists, err := uc.educationRepo.ExistsByID(ctx, education.EducationID)
	if err != nil {
		uc.logger.Error("Failed to check if education exists: %v", err)
		return nil, domain.NewInternalError("failed to check education existence", err)
	}
	if !exists {
		uc.logger.Error("Education not found for ID %d", education.EducationID)
		return nil, domain.NewNotFoundError("Education", fmt.Sprint(education.EducationID))
	}

	education.MarkAsUpdated()
	updatedEducation, err := uc.educationRepo.Update(ctx, educationID, education)
	if err != nil {
		uc.logger.Error("Failed to update education: %v", err)
		return nil, domain.NewInternalError("failed to update education", err)
	}

	return updatedEducation, nil
}

func (uc *EducationUseCase) PatchEducation(ctx context.Context, educationID int, updates *entities.Education) (*entities.Education, error) {
	if educationID <= 0 {
		uc.logger.Error("Education ID is required")
		return nil, domain.NewValidationError("educationID", "education ID must be positive", nil)
	}

	existingEducation, err := uc.educationRepo.GetByID(ctx, educationID)
	if err != nil {
		uc.logger.Error("Failed to get education by ID %d: %v", educationID, err)
		return nil, domain.NewInternalError("failed to get education", err)
	}

	if existingEducation == nil {
		uc.logger.Error("Education not found: %d", educationID)
		return nil, domain.NewNotFoundError("Education", fmt.Sprint(educationID))
	}

	if updates.Degree != "" {
		existingEducation.Degree = updates.Degree
	}
	if updates.Institution != "" {
		existingEducation.Institution = updates.Institution
	}
	if !updates.StartDate.IsZero() {
		existingEducation.StartDate = updates.StartDate
	}
	if updates.EndDate != nil {
		existingEducation.EndDate = updates.EndDate
	}
	if updates.Description != "" {
		existingEducation.Description = updates.Description
	}

	existingEducation.MarkAsUpdated()

	patchedEducation, err := uc.educationRepo.Patch(ctx, educationID, existingEducation)
	if err != nil {
		uc.logger.Error("Failed to update education: %v", err)
		return nil, domain.NewInternalError("failed to update education", err)
	}

	return patchedEducation, nil
}

func (uc *EducationUseCase) DeleteEducation(ctx context.Context, educationID int) error {
	if educationID <= 0 {
		uc.logger.Error("Invalid education ID: %d", educationID)
		return domain.NewValidationError("educationID", "education ID must be positive", nil)
	}

	exists, err := uc.educationRepo.ExistsByID(ctx, educationID)
	if err != nil {
		uc.logger.Error("Failed to check if education exists: %v", err)
		return domain.NewInternalError("failed to check education existence", err)
	}
	if !exists {
		uc.logger.Error("Education not found for ID %d", educationID)
		return domain.NewNotFoundError("Education", fmt.Sprint(educationID))
	}

	err = uc.educationRepo.Delete(ctx, educationID)
	if err != nil {
		uc.logger.Error("Failed to delete education: %v", err)
		return domain.NewInternalError("failed to delete education", err)
	}

	return nil
}

func (uc *EducationUseCase) GetAllEducations(ctx context.Context) ([]*entities.Education, error) {
	educations, err := uc.educationRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all educations: %v", err)
		return nil, domain.NewInternalError("failed to retrieve educations", err)
	}

	return educations, nil
}

func (uc *EducationUseCase) GetCurrentEducations(ctx context.Context, userID int) ([]*entities.Education, error) {
	if userID <= 0 {
		uc.logger.Error("Invalid user ID: %d", userID)
		return nil, domain.NewValidationError("userID", "user ID must be positive", nil)
	}

	educations, err := uc.educationRepo.GetCurrentEducations(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get current educations for user %d: %v", userID, err)
		return nil, domain.NewInternalError("failed to retrieve current educations", err)
	}

	return educations, nil
}
