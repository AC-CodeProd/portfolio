package handler

import (
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
	projectDto "portfolio/dto/project"
	"portfolio/logger"
	"portfolio/shared"
	"strconv"
	"time"
)

type projectHandler struct {
	AbstractHandler
	projectUseCase *usecases.ProjectUseCase
	logger         *logger.Logger
}

func NewProjectHandler(settingUseCase *usecases.SettingUseCase, projectUseCase *usecases.ProjectUseCase, logger *logger.Logger) []*routes.NamedRoute {
	projectHandler := projectHandler{
		AbstractHandler: AbstractHandler{
			settingUseCase: settingUseCase,
		},
		projectUseCase: projectUseCase,
		logger:         logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetProjectsHandler",
			Pattern: "GET /projects",
			Handler: projectHandler.GetProjects,
		},
		{
			Name:    "GetProjectHandler",
			Pattern: "GET /projects/{id}",
			Handler: projectHandler.GetProject,
		},
	}
}

// GetProjects
//
//	@Summary		Get all active projects
//	@Description	Retrieve all active projects for the portfolio
//	@Tags			Projects
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.ProjectListResponse}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/projects [get]
func (ph *projectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(ph.settingUseCase, ctx, w)

	if err != nil {
		ph.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	projects, err := ph.projectUseCase.GetProjectsByUserID(ctx, portfolioOwnerID)
	if err != nil {
		ph.logger.Error("Failed to get projects: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := projectDto.FromProjectsEntityToResponse(projects, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if response == nil {
		ph.logger.Error("No projects found for user ID %d", portfolioOwnerID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Projects", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetProject
//
//	@Summary		Get a specific project
//	@Description	Retrieve a specific project by ID
//	@Tags			Projects
//	@Produce		json
//	@Param			id	path		int	true	"Project ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.ProjectResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/projects/{id} [get]
func (ph *projectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectIDStr := r.PathValue("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil || projectID <= 0 {
		ph.logger.Error("Invalid project ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid project ID", "id", &err))
		return
	}

	project, err := ph.projectUseCase.GetProjectByID(ctx, projectID)
	if err != nil {
		ph.logger.Error("Failed to get project: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(ph.settingUseCase, ctx, w)

	if err != nil {
		ph.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	if project.Status != "active" || project.UserID != portfolioOwnerID {
		ph.logger.Error("Unauthorized access to project %d", projectID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Project", strconv.Itoa(projectID)))
		return
	}

	response := projectDto.FromProjectEntityToResponse(project, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}
