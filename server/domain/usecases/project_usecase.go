package usecases

import (
	"context"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	dto "portfolio/dto/project"
	"portfolio/logger"
	"strconv"
)

type ProjectUseCase struct {
	projectRepo interfaces.ProjectRepository
	userRepo    interfaces.UserRepository
	settingRepo interfaces.SettingRepository
	logger      *logger.Logger
}

func NewProjectUseCase(projectRepo interfaces.ProjectRepository, userRepo interfaces.UserRepository, settingRepo interfaces.SettingRepository, logger *logger.Logger) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo: projectRepo,
		userRepo:    userRepo,
		settingRepo: settingRepo,
		logger:      logger,
	}
}

func (uc *ProjectUseCase) GetProjectsByUserID(ctx context.Context, userID int) ([]*entities.Project, error) {
	return uc.projectRepo.GetAll(ctx, userID)
}

func (uc *ProjectUseCase) GetProjectByID(ctx context.Context, projectID int) (*entities.Project, error) {
	return uc.projectRepo.GetByID(ctx, projectID)
}

func (uc *ProjectUseCase) CreateProject(ctx context.Context, userID int, req *dto.CreateProjectRequest) (*entities.Project, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	project := &entities.Project{
		UserID:           userID,
		Title:            req.Title,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Technologies:     req.Technologies,
		GithubURL:        req.GithubURL,
		ImageURL:         req.ImageURL,
		Status:           req.Status,
	}

	if !project.HasRequiredFields() {
		uc.logger.Error("Required fields are missing for project: %v", project)
		return nil, domain.NewValidationError("Required fields are missing", "project", nil)
	}

	if !project.HasValidStatus() {
		uc.logger.Error("Invalid project status for project: %v", project)
		return nil, domain.NewValidationError("Invalid project status", "status", nil)
	}

	createdProject, err := uc.projectRepo.Create(ctx, project)
	if err != nil {
		uc.logger.Error("Failed to create project: %v", err)
		return nil, err
	}

	return createdProject, nil
}

func (uc *ProjectUseCase) UpdateProject(ctx context.Context, projectID int, req *dto.UpdateProjectRequest) (*entities.Project, error) {
	if err := req.Validate(); err != nil {
		uc.logger.Error("Invalid update project request: %v", err)
		return nil, err
	}

	existingProject, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		uc.logger.Error("Failed to get existing project: %v", err)
		return nil, err
	}

	existingProject.Title = req.Title
	existingProject.Description = req.Description
	existingProject.ShortDescription = req.ShortDescription
	existingProject.Technologies = req.Technologies
	existingProject.GithubURL = req.GithubURL
	existingProject.ImageURL = req.ImageURL

	if !existingProject.SetStatus(req.Status) {
		uc.logger.Error("Invalid project status for project: %v", existingProject)
		return nil, domain.NewValidationError("Invalid project status", "status", nil)
	}

	if !existingProject.HasRequiredFields() {
		uc.logger.Error("Required fields are missing after update for project: %v", existingProject)
		return nil, domain.NewValidationError("Required fields are missing after update", "project", nil)
	}

	updatedProject, err := uc.projectRepo.Update(ctx, projectID, existingProject)
	if err != nil {
		uc.logger.Error("Failed to update project: %v", err)
		return nil, err
	}

	return updatedProject, nil
}

func (uc *ProjectUseCase) PatchProject(ctx context.Context, projectID int, req *dto.PatchProjectRequest) (*entities.Project, error) {
	if err := req.Validate(); err != nil {
		uc.logger.Error("Invalid patch project request: %v", err)
		return nil, err
	}

	existingProject, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		uc.logger.Error("Failed to get existing project: %v", err)
		return nil, err
	}

	if req.Title != "" {
		existingProject.Title = req.Title
	}
	if req.Description != "" {
		existingProject.Description = req.Description
	}
	if req.ShortDescription != "" {
		existingProject.ShortDescription = req.ShortDescription
	}
	if req.Technologies != "" {
		existingProject.Technologies = req.Technologies
	}
	if req.GithubURL != "" {
		existingProject.GithubURL = req.GithubURL
	}
	if req.ImageURL != "" {
		existingProject.ImageURL = req.ImageURL
	}
	if req.Status != "" {
		if !existingProject.SetStatus(req.Status) {
			uc.logger.Error("Invalid project status for project: %v", existingProject)
			return nil, domain.NewValidationError("Invalid project status", "status", nil)
		}
	} else {
		existingProject.MarkAsUpdated()
	}

	patchedProject, err := uc.projectRepo.Patch(ctx, projectID, existingProject)
	if err != nil {
		uc.logger.Error("Failed to patch project: %v", err)
		return nil, err
	}

	return patchedProject, nil
}

func (uc *ProjectUseCase) DeleteProject(ctx context.Context, projectID int) error {
	_, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		uc.logger.Error("Failed to get existing project: %v", err)
		return err
	}

	return uc.projectRepo.Delete(ctx, projectID)
}

func (uc *ProjectUseCase) ValidateProjectOwnership(ctx context.Context, projectID, userID int) error {
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		uc.logger.Error("Failed to get project by ID: %v", err)
		return err
	}

	if !project.BelongsToUser(userID) {
		uc.logger.Error("User %d does not have permission to access project %d", userID, projectID)
		return domain.NewForbiddenError("You don't have permission to access this project")
	}

	return nil
}

func (uc *ProjectUseCase) ValidateCreateProjectRequest(req *dto.CreateProjectRequest) error {
	return req.Validate()
}

func (uc *ProjectUseCase) ValidateUpdateProjectRequest(req *dto.UpdateProjectRequest) error {
	return req.Validate()
}

func (uc *ProjectUseCase) ValidatePatchProjectRequest(req *dto.PatchProjectRequest) error {
	return req.Validate()
}

func (uc *ProjectUseCase) ValidateProjectID(projectIDStr string) (int, error) {
	if projectIDStr == "" {
		uc.logger.Error("Project ID is required")
		return 0, domain.NewValidationError("Project ID is required", "project_id", nil)
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		uc.logger.Error("Failed to convert project ID to integer: %v", err)
		return 0, domain.NewValidationError("Project ID must be a valid integer", "project_id", nil)
	}

	if projectID <= 0 {
		uc.logger.Error("Project ID must be a positive integer: %d", projectID)
		return 0, domain.NewValidationError("Project ID must be a positive integer", "project_id", nil)
	}

	return projectID, nil
}
