package usecases

import (
	"context"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/domain/utils"
	"portfolio/domain/validation"
	dto "portfolio/dto/personal_info"
	"portfolio/logger"
)

type PersonalInfoUseCase struct {
	personalInfoRepo interfaces.PersonalInfoRepository
	logger           *logger.Logger
}

func NewPersonalInfoUseCase(personalInfoRepo interfaces.PersonalInfoRepository, logger *logger.Logger) *PersonalInfoUseCase {
	return &PersonalInfoUseCase{
		personalInfoRepo: personalInfoRepo,
		logger:           logger,
	}
}

func validatePersonalInfoRequest(request *dto.CreatePersonalInfoRequest) []*domain.DomainError {
	validator := validation.NewValidator()

	validator.Required("first_name", request.FirstName).
		Required("last_name", request.LastName).
		Required("professional_title", request.ProfessionalTitle).
		Required("bio", request.Bio).
		Required("location", request.Location).
		Required("phone_number", request.PhoneNumber).
		Required("interests", request.Interests).
		Required("profile_picture", request.ProfilePicture)

	if request.PhoneNumber != "" {
		validator.Phone("phone_number", request.PhoneNumber)
	}

	if request.DateOfBirth != "" {
		dateOfBirth, err := utils.ParseDate(request.DateOfBirth)
		if err != nil {
			return nil
		}
		validator.DateNotFuture("date_of_birth", dateOfBirth.Time())
	}

	if request.ResumeURL != "" {
		validator.URL("resume_url", request.ResumeURL)
	}
	if request.WebsiteURL != "" {
		validator.URL("website_url", request.WebsiteURL)
	}
	if request.LinkedinURL != "" {
		validator.URL("linkedin_url", request.LinkedinURL)
	}
	if request.GithubURL != "" {
		validator.URL("github_url", request.GithubURL)
	}
	if request.XURL != "" {
		validator.URL("x_url", request.XURL)
	}

	return validator.Errors()
}

func validateUpdatePersonalInfoRequest(request *dto.UpdatePersonalInfoRequest) []*domain.DomainError {
	return validatePersonalInfoRequest(&request.CreatePersonalInfoRequest)
}

func validatePatchPersonalInfoRequest(request *dto.PatchPersonalInfoRequest) []*domain.DomainError {
	validator := validation.NewValidator()

	if request.PhoneNumber != nil && *request.PhoneNumber != "" {
		validator.Phone("phone_number", *request.PhoneNumber)
	}

	if request.ResumeURL != nil && *request.ResumeURL != "" {
		validator.URL("resume_url", *request.ResumeURL)
	}
	if request.WebsiteURL != nil && *request.WebsiteURL != "" {
		validator.URL("website_url", *request.WebsiteURL)
	}
	if request.LinkedinURL != nil && *request.LinkedinURL != "" {
		validator.URL("linkedin_url", *request.LinkedinURL)
	}
	if request.GithubURL != nil && *request.GithubURL != "" {
		validator.URL("github_url", *request.GithubURL)
	}
	if request.XURL != nil && *request.XURL != "" {
		validator.URL("x_url", *request.XURL)
	}

	return validator.Errors()
}

func (uc *PersonalInfoUseCase) ValidateCreatePersonalInfoRequest(request *dto.CreatePersonalInfoRequest) error {
	if request == nil {
		return domain.NewValidationError("Request cannot be nil", "", nil)
	}

	errors := validatePersonalInfoRequest(request)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

func (uc *PersonalInfoUseCase) ValidateUpdatePersonalInfoRequest(request *dto.UpdatePersonalInfoRequest) error {
	if request == nil {
		uc.logger.Error("Request cannot be nil")
		return domain.NewValidationError("Request cannot be nil", "", nil)
	}

	errors := validateUpdatePersonalInfoRequest(request)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

func (uc *PersonalInfoUseCase) ValidatePatchPersonalInfoRequest(request *dto.PatchPersonalInfoRequest) error {
	if request == nil {
		uc.logger.Error("Request cannot be nil")
		return domain.NewValidationError("Request cannot be nil", "", nil)
	}

	errors := validatePatchPersonalInfoRequest(request)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

func (uc *PersonalInfoUseCase) GetUserByID(ctx context.Context, userID int) (*entities.User, error) {
	if userID <= 0 {
		uc.logger.Error("Invalid user ID: %d", userID)
		return nil, domain.NewValidationError("User ID must be positive", "user_id", nil)
	}

	user, err := uc.personalInfoRepo.GetUserByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user by ID: %v", err)
		return nil, domain.NewDatabaseError("retrieving user", err)
	}
	if user == nil {
		uc.logger.Warn("No user found with ID %d", userID)
		return nil, domain.NewNotFoundError("User", fmt.Sprintf("%d", userID))
	}
	return user, nil
}

func (uc *PersonalInfoUseCase) GetPersonalInfoByUserID(ctx context.Context, userID int) (*entities.PersonalInfo, error) {
	if userID <= 0 {
		uc.logger.Error("Invalid user ID provided")
		return nil, domain.NewValidationError("User ID must be positive", "user_id", nil)
	}
	user, err := uc.personalInfoRepo.GetUserByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user by ID: %v", err)
		return nil, domain.NewDatabaseError("retrieving user", err)
	}
	if user == nil {
		uc.logger.Warn("No user found with ID %d", userID)
		return nil, domain.NewNotFoundError("User", fmt.Sprintf("%d", userID))
	}
	personalInfo, err := uc.personalInfoRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get personal info by user ID: %v", err)
		return nil, domain.NewDatabaseError("retrieving personal info", err)
	}
	if personalInfo == nil {
		uc.logger.Warn("No personal info found for user ID %d", userID)
		return nil, domain.NewNotFoundError("Personal Info", fmt.Sprintf("user_id:%d", userID))
	}
	return personalInfo, nil
}

func (uc *PersonalInfoUseCase) CreatePersonalInfo(ctx context.Context, userID int, request *dto.CreatePersonalInfoRequest) (*entities.PersonalInfo, error) {
	if err := uc.ValidateCreatePersonalInfoRequest(request); err != nil {
		uc.logger.Error("Invalid request: %v", err)
		return nil, err
	}

	personalInfo := dto.FromPersonalInfoRequestToEntity(userID, request)
	if personalInfo == nil {
		uc.logger.Error("Failed to convert request to entity")
		return nil, domain.NewInternalError("Failed to convert request to entity", nil)
	}

	createdInfo, err := uc.personalInfoRepo.Create(ctx, personalInfo)
	if err != nil {
		uc.logger.Error("Failed to create personal info: %v", err)
		return nil, domain.NewDatabaseError("creating personal info", err)
	}

	return createdInfo, nil
}

func (uc *PersonalInfoUseCase) UpdatePersonalInfo(ctx context.Context, personalInfoId int, request *dto.UpdatePersonalInfoRequest) (*entities.PersonalInfo, error) {
	if err := uc.ValidateUpdatePersonalInfoRequest(request); err != nil {
		uc.logger.Error("Invalid request: %v", err)
		return nil, err
	}

	if personalInfoId <= 0 {
		uc.logger.Error("Invalid personal info ID provided")
		return nil, domain.NewValidationError("Personal info ID must be positive", "id", nil)
	}

	existingPersonalInfo, err := uc.personalInfoRepo.GetByID(ctx, personalInfoId)
	if err != nil {
		uc.logger.Error("Failed to check if personal info exists by ID: %v", err)
		return nil, domain.NewDatabaseError("checking personal info existence", err)
	}

	if existingPersonalInfo == nil {
		uc.logger.Warn("No personal info found with ID %d", personalInfoId)
		return nil, domain.NewNotFoundError("personal information", fmt.Sprintf("%d", personalInfoId))
	}

	personalInfo := dto.FromUpdatePersonalInfoRequestToEntity(personalInfoId, request)
	if personalInfo == nil {
		uc.logger.Error("Failed to convert request to entity")
		return nil, domain.NewInternalError("Failed to convert request to entity", nil)
	}

	updatedPersonalInfo, err := uc.personalInfoRepo.Update(ctx, personalInfoId, personalInfo)
	if err != nil {
		uc.logger.Error("Failed to update personal info: %v", err)
		return nil, domain.NewDatabaseError("updating personal info", err)
	}

	return updatedPersonalInfo, nil
}

func (uc *PersonalInfoUseCase) PatchPersonalInfo(ctx context.Context, personalInfoId int, request *dto.PatchPersonalInfoRequest) (*entities.PersonalInfo, error) {
	if err := uc.ValidatePatchPersonalInfoRequest(request); err != nil {
		uc.logger.Error("Invalid patch request: %v", err)
		return nil, err
	}

	if personalInfoId <= 0 {
		uc.logger.Error("Invalid personal info ID provided")
		return nil, domain.NewValidationError("Personal info ID must be positive", "id", nil)
	}

	existingPersonalInfo, err := uc.personalInfoRepo.GetByID(ctx, personalInfoId)
	if err != nil {
		uc.logger.Error("Failed to check if personal info exists by ID: %v", err)
		return nil, domain.NewDatabaseError("checking personal info existence", err)
	}

	if existingPersonalInfo == nil {
		uc.logger.Error("No personal info found with ID %d", personalInfoId)
		return nil, domain.NewNotFoundError("personal information", fmt.Sprintf("%d", personalInfoId))
	}

	personalInfo, err := dto.FromPatchPersonalInfoRequestToEntity(personalInfoId, request)
	if err != nil {
		uc.logger.Error("Failed to convert patch request to entity: %v", err)
		return nil, domain.NewInternalError("Failed to convert patch request to entity", err)
	}

	patchedPersonalInfo, err := uc.personalInfoRepo.Patch(ctx, personalInfoId, personalInfo)
	if err != nil {
		uc.logger.Error("Failed to patch personal info: %v", err)
		return nil, domain.NewDatabaseError("patching personal info", err)
	}

	return patchedPersonalInfo, nil
}

func (uc *PersonalInfoUseCase) DeletePersonalInfo(ctx context.Context, personalInfoId int) error {
	if personalInfoId <= 0 {
		uc.logger.Error("Invalid personal info ID provided")
		return domain.NewValidationError("Personal info ID must be positive", "id", nil)
	}

	existingPersonalInfo, err := uc.personalInfoRepo.GetByID(ctx, personalInfoId)
	if err != nil {
		uc.logger.Error("Failed to check if personal info exists by ID: %v", err)
		return domain.NewDatabaseError("checking personal info existence", err)
	}

	if existingPersonalInfo == nil {
		uc.logger.Warn("No personal info found with ID %d", personalInfoId)
		return domain.NewNotFoundError("personal information", fmt.Sprintf("%d", personalInfoId))
	}

	err = uc.personalInfoRepo.Delete(ctx, personalInfoId)
	if err != nil {
		uc.logger.Error("Failed to delete personal info: %v", err)
		return domain.NewDatabaseError("deleting personal info", err)
	}

	return nil
}

func (uc *PersonalInfoUseCase) ExistsPersonalInfoByUserID(ctx context.Context, userID int) (bool, error) {
	if userID <= 0 {
		uc.logger.Error("Invalid user ID provided")
		return false, domain.NewValidationError("User ID must be positive", "user_id", nil)
	}

	exists, err := uc.personalInfoRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to check if personal info exists by user ID: %v", err)
		return false, domain.NewDatabaseError("checking personal info existence by user ID", err)
	}

	return exists, nil
}
