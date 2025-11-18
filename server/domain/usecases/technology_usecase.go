package usecases

import (
	"context"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
)

type TechnologyUseCase struct {
	technologyRepo interfaces.TechnologyRepository
	userRepo       interfaces.UserRepository
	logger         *logger.Logger
}

func NewTechnologyUseCase(technologyRepo interfaces.TechnologyRepository, userRepo interfaces.UserRepository, logger *logger.Logger) *TechnologyUseCase {
	return &TechnologyUseCase{
		technologyRepo: technologyRepo,
		userRepo:       userRepo,
		logger:         logger,
	}
}

func (uc *TechnologyUseCase) CreateTechnology(ctx context.Context, technology *entities.Technology) (*entities.Technology, error) {
	if !technology.HasRequiredFields() {
		uc.logger.Error("Required fields are missing for technology: %v", technology)
		return nil, domain.NewValidationError("technology", "technology name, icon URL, and user ID are required", nil)
	}

	userExists, err := uc.userRepo.ExistsByID(ctx, technology.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if user exists: %v", err)
		return nil, domain.NewInternalError("failed to validate user", err)
	}
	if !userExists {
		uc.logger.Error("User %d does not exist", technology.UserID)
		return nil, domain.NewNotFoundError("User", fmt.Sprint(technology.UserID))
	}

	exists, err := uc.technologyRepo.ExistsByNameAndUserID(ctx, technology.Name, technology.UserID)
	if err != nil {
		uc.logger.Error("Failed to check if technology exists: %v", err)
		return nil, domain.NewInternalError("failed to check technology existence", err)
	}
	if exists {
		uc.logger.Error("Technology %s already exists for user %d", technology.Name, technology.UserID)
		return nil, domain.NewAlreadyExistsError("Technology", technology.Name)
	}

	createdTechnology, err := uc.technologyRepo.Create(ctx, technology)
	if err != nil {
		uc.logger.Error("Failed to create technology: %v", err)
		return nil, domain.NewInternalError("failed to create technology", err)
	}

	return createdTechnology, nil
}

func (uc *TechnologyUseCase) GetTechnologyByID(ctx context.Context, technologyID int) (*entities.Technology, error) {
	if technologyID <= 0 {
		uc.logger.Error("Technology ID is required")
		return nil, domain.NewValidationError("technologyID", "technology ID must be positive", nil)
	}

	technology, err := uc.technologyRepo.GetByID(ctx, technologyID)
	if err != nil {
		uc.logger.Error("Failed to get technology by ID %d: %v", technologyID, err)
		return nil, domain.NewInternalError("failed to retrieve technology", err)
	}

	if technology == nil {
		uc.logger.Error("Technology not found: %d", technologyID)
		return nil, domain.NewNotFoundError("Technology", fmt.Sprint(technologyID))
	}

	return technology, nil
}

func (uc *TechnologyUseCase) GetTechnologiesByUserID(ctx context.Context, userID int) ([]*entities.Technology, error) {
	if userID <= 0 {
		uc.logger.Error("User ID is required")
		return nil, domain.NewValidationError("userID", "user ID must be positive", nil)
	}

	technologies, err := uc.technologyRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get technologies for user %d: %v", userID, err)
		return nil, domain.NewInternalError("failed to retrieve technologies", err)
	}

	return technologies, nil
}

func (uc *TechnologyUseCase) UpdateTechnology(ctx context.Context, technologyID int, technology *entities.Technology) (*entities.Technology, error) {
	if !technology.HasRequiredFields() {
		uc.logger.Error("Required fields are missing for technology: %v", technology)
		return nil, domain.NewValidationError("technology", "technology name, icon URL, and user ID are required", nil)
	}

	exists, err := uc.technologyRepo.ExistsByID(ctx, technology.TechnologyID)
	if err != nil {
		uc.logger.Error("Failed to check if technology exists: %v", err)
		return nil, domain.NewInternalError("failed to check technology existence", err)
	}
	if !exists {
		uc.logger.Error("Technology not found: %d", technology.TechnologyID)
		return nil, domain.NewNotFoundError("Technology", fmt.Sprint(technology.TechnologyID))
	}

	technology.MarkAsUpdated()
	updatedTechnology, err := uc.technologyRepo.Update(ctx, technologyID, technology)
	if err != nil {
		uc.logger.Error("Failed to update technology: %v", err)
		return nil, domain.NewInternalError("failed to update technology", err)
	}

	return updatedTechnology, nil
}

func (uc *TechnologyUseCase) PatchTechnology(ctx context.Context, technologyID int, patchData *entities.Technology) (*entities.Technology, error) {
	if technologyID <= 0 {
		return nil, domain.NewValidationError("technologyID", "technology ID must be positive", nil)
	}

	existingTechnology, err := uc.technologyRepo.GetByID(ctx, technologyID)
	if err != nil {
		uc.logger.Error("Failed to get technology: %v", err)
		return nil, domain.NewInternalError("failed to retrieve technology", err)
	}

	if existingTechnology == nil {
		uc.logger.Error("Technology not found: %d", technologyID)
		return nil, domain.NewNotFoundError("Technology", fmt.Sprint(technologyID))
	}

	if patchData.Name != "" {
		existingTechnology.Name = patchData.Name
	}
	if patchData.IconURL != "" {
		existingTechnology.IconURL = patchData.IconURL
	}

	existingTechnology.MarkAsUpdated()
	patchedTechnology, err := uc.technologyRepo.Patch(ctx, technologyID, existingTechnology)
	if err != nil {
		uc.logger.Error("Failed to patch technology: %v", err)
		return nil, domain.NewInternalError("failed to patch technology", err)
	}

	return patchedTechnology, nil
}

func (uc *TechnologyUseCase) DeleteTechnology(ctx context.Context, technologyID int) error {
	if technologyID <= 0 {
		uc.logger.Error("Technology ID is required")
		return domain.NewValidationError("technologyID", "technology ID must be positive", nil)
	}

	exists, err := uc.technologyRepo.ExistsByID(ctx, technologyID)
	if err != nil {
		uc.logger.Error("Failed to check if technology exists: %v", err)
		return domain.NewInternalError("failed to check technology existence", err)
	}
	if !exists {
		uc.logger.Error("Technology not found: %d", technologyID)
		return domain.NewNotFoundError("Technology", fmt.Sprint(technologyID))
	}

	err = uc.technologyRepo.Delete(ctx, technologyID)
	if err != nil {
		uc.logger.Error("Failed to delete technology: %v", err)
		return domain.NewInternalError("failed to delete technology", err)
	}

	return nil
}

func (uc *TechnologyUseCase) GetAllTechnologies(ctx context.Context) ([]*entities.Technology, error) {
	technologies, err := uc.technologyRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all technologies: %v", err)
		return nil, domain.NewInternalError("failed to retrieve technologies", err)
	}

	return technologies, nil
}

func (uc *TechnologyUseCase) GetTechnologiesByNames(ctx context.Context, names []string, userID int) ([]*entities.Technology, error) {
	if userID <= 0 {
		uc.logger.Error("User ID is required")
		return nil, domain.NewValidationError("userID", "user ID must be positive", nil)
	}

	if len(names) == 0 {
		uc.logger.Error("Technology names are required")
		return []*entities.Technology{}, nil
	}

	technologies, err := uc.technologyRepo.GetByNames(ctx, names, userID)
	if err != nil {
		uc.logger.Error("Failed to get technologies by names for user %d: %v", userID, err)
		return nil, domain.NewInternalError("failed to retrieve technologies by names", err)
	}

	return technologies, nil
}
