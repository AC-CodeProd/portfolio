package admin

import (
	"encoding/json"
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
		AbstractHandler: AbstractHandler{settingUseCase: settingUseCase},
		projectUseCase:  projectUseCase,
		logger:          logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetAdminProjectsHandler",
			Pattern: "GET /projects",
			Handler: projectHandler.GetProjects,
		},
		{
			Name:    "PostAdminProjectHandler",
			Pattern: "POST /projects",
			Handler: projectHandler.CreateProject,
		},
		{
			Name:    "GetAdminProjectHandler",
			Pattern: "GET /projects/{id}",
			Handler: projectHandler.GetProject,
		},
		{
			Name:    "PutAdminProjectHandler",
			Pattern: "PUT /projects/{id}",
			Handler: projectHandler.UpdateProject,
		},
		{
			Name:    "PatchAdminProjectHandler",
			Pattern: "PATCH /projects/{id}",
			Handler: projectHandler.PatchProject,
		},
		{
			Name:    "DeleteAdminProjectHandler",
			Pattern: "DELETE /projects/{id}",
			Handler: projectHandler.DeleteProject,
		},
	}
}

// GetProjects
//
//	@Summary		Get all admin projects
//	@Description	Retrieve all projects for the authenticated admin user
//	@Tags			Admin Projects
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.ProjectListResponse}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/projects [get]
//	@Security		BearerAuth
func (ph *projectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ph.getUserIDFromContext(w, r)
	if !ok {
		return
	}

	projectsEntity, err := ph.projectUseCase.GetProjectsByUserID(ctx, userID)
	if err != nil {
		ph.logger.Error("Failed to get projects for user %d: %v", userID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := projectDto.FromProjectsEntityToResponse(projectsEntity,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})

	if response == nil {
		ph.logger.Error("No projects found for user %d", userID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Projects", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetProject
//
//	@Summary		Get a specific admin project
//	@Description	Retrieve a specific project by ID for admin management
//	@Tags			Admin Projects
//	@Produce		json
//	@Param			id	path		int	true	"Project ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.ProjectResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/projects/{id} [get]
//	@Security		BearerAuth
func (ph *projectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		ph.logger.Error("Project ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Project ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ph.logger.Error("Invalid project ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid project ID", "id", &err))
		return
	}

	projectEntity, err := ph.projectUseCase.GetProjectByID(ctx, id)
	if err != nil {
		ph.logger.Error("Failed to get project %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	if projectEntity == nil {
		ph.logger.Error("Project not found: %d", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Project", idStr))
		return
	}

	response := projectDto.FromProjectEntityToResponse(projectEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// CreateProject
//
//	@Summary		Create a new project
//	@Description	Create a new project for the authenticated admin user
//	@Tags			Admin Projects
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateProjectRequest	true	"Project creation request"
//	@Success		201		{object}	shared.APIResponse{data=dto.ProjectResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/projects [post]
//	@Security		BearerAuth
func (ph *projectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request projectDto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ph.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		ph.logger.Error("Request validation failed: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	userID, ok := ph.getUserIDFromContext(w, r)
	if !ok {
		ph.logger.Error("Failed to get user ID from context")
		return
	}

	projectEntity, err := ph.projectUseCase.CreateProject(ctx, userID, &request)
	if err != nil {
		ph.logger.Error("Failed to create project: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := projectDto.FromProjectEntityToResponse(projectEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusCreated, response)
}

// UpdateProject
//
//	@Summary		Update an existing project
//	@Description	Update an existing project by ID for the authenticated admin user
//	@Tags			Admin Projects
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Project ID"
//	@Param			request	body		dto.UpdateProjectRequest	true	"Project update request"
//	@Success		200		{object}	shared.APIResponse{data=dto.ProjectResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/projects/{id} [put]
//	@Security		BearerAuth
func (ph *projectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		ph.logger.Error("Project ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Project ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ph.logger.Error("Invalid project ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid project ID", "id", &err))
		return
	}

	var request projectDto.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ph.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		ph.logger.Error("Request validation failed: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	projectEntity, err := ph.projectUseCase.UpdateProject(ctx, id, &request)
	if err != nil {
		ph.logger.Error("Failed to update project %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := projectDto.FromProjectEntityToResponse(projectEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// PatchProject
//
//	@Summary		Partially update a project
//	@Description	Partially update a project by ID for the authenticated admin user
//	@Tags			Admin Projects
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Project ID"
//	@Param			request	body		dto.PatchProjectRequest		true	"Patch project request"
//	@Success		200		{object}	shared.APIResponse{data=dto.ProjectResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/projects/{id} [patch]
//	@Security		BearerAuth
func (ph *projectHandler) PatchProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		ph.logger.Error("Project ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Project ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ph.logger.Error("Invalid project ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid project ID", "id", &err))
		return
	}

	userID, ok := ph.getUserIDFromContext(w, r)
	if !ok {
		ph.logger.Error("Failed to get user ID from context")
		return
	}

	err = ph.projectUseCase.ValidateProjectOwnership(ctx, id, userID)
	if err != nil {
		ph.logger.Error("Project ownership validation failed for project %d and user %d: %v", id, userID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	var req projectDto.PatchProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid JSON in request body", "request", &err))
		return
	}

	projectEntity, err := ph.projectUseCase.PatchProject(ctx, id, &req)
	if err != nil {
		ph.logger.Error("Failed to patch project %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := projectDto.FromProjectEntityToResponse(projectEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// DeleteProject
//
//	@Summary		Delete a project
//	@Description	Delete a project by ID for the authenticated admin user
//	@Tags			Admin Projects
//	@Produce		json
//	@Param			id	path		int	true	"Project ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/projects/{id} [delete]
//	@Security		BearerAuth
func (ph *projectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		ph.logger.Error("Project ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Project ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ph.logger.Error("Invalid project ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid project ID", "id", &err))
		return
	}

	err = ph.projectUseCase.DeleteProject(ctx, id)
	if err != nil {
		ph.logger.Error("Failed to delete project %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
