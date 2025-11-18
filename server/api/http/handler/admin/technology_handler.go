package admin

import (
	"encoding/json"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/usecases"
	technologyDto "portfolio/dto/technology"
	"portfolio/logger"
	"portfolio/shared"
	"strconv"
	"time"
)

type technologyHandler struct {
	AbstractHandler
	technologyUseCase *usecases.TechnologyUseCase
	logger            *logger.Logger
}

func NewTechnologyHandler(settingUseCase *usecases.SettingUseCase, technologyUseCase *usecases.TechnologyUseCase, logger *logger.Logger) []*routes.NamedRoute {
	technologyHandler := technologyHandler{
		AbstractHandler:   AbstractHandler{settingUseCase: settingUseCase},
		technologyUseCase: technologyUseCase,
		logger:            logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetAdminTechnologiesHandler",
			Pattern: "GET /technologies",
			Handler: technologyHandler.GetTechnologies,
		},
		{
			Name:    "PostAdminTechnologyHandler",
			Pattern: "POST /technologies",
			Handler: technologyHandler.CreateTechnology,
		},
		{
			Name:    "PostAdminBulkTechnologiesHandler",
			Pattern: "POST /technologies/bulk",
			Handler: technologyHandler.CreateBulkTechnologies,
		},
		{
			Name:    "GetAdminTechnologyHandler",
			Pattern: "GET /technologies/{id}",
			Handler: technologyHandler.GetTechnology,
		},
		{
			Name:    "PutAdminTechnologyHandler",
			Pattern: "PUT /technologies/{id}",
			Handler: technologyHandler.UpdateTechnology,
		},
		{
			Name:    "PatchAdminTechnologyHandler",
			Pattern: "PATCH /technologies/{id}",
			Handler: technologyHandler.PatchTechnology,
		},
		{
			Name:    "DeleteAdminTechnologyHandler",
			Pattern: "DELETE /technologies/{id}",
			Handler: technologyHandler.DeleteTechnology,
		},
	}
}

// GetTechnologies
//
//	@Summary		Get all technologies
//	@Description	Get a list of all technologies for the authenticated admin user
//	@Tags			Admin Technologies
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.TechnologyListResponse}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/technologies [get]
//	@Security		BearerAuth
func (th *technologyHandler) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := th.getUserIDFromContext(w, r)
	if !ok {
		th.logger.Error("Failed to get user ID from context")
		return
	}

	technologies, err := th.technologyUseCase.GetTechnologiesByUserID(ctx, userID)
	if err != nil {
		th.logger.Error("Failed to get technologies for user %d: %v", userID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := technologyDto.FromTechnologiesEntityToResponse(technologies, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if response == nil {
		th.logger.Error("No technologies found for user %d", userID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Technologies", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetTechnology
//
//	@Summary		Get a technology by ID
//	@Description	Get a specific technology by ID for the authenticated admin user
//	@Tags			Admin Technologies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Technology ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.TechnologyResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/technologies/{id} [get]
//	@Security		BearerAuth
func (th *technologyHandler) GetTechnology(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		th.logger.Error("Technology ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Technology ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		th.logger.Error("Invalid technology ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology ID", "id", &err))
		return
	}

	technologyEntity, err := th.technologyUseCase.GetTechnologyByID(ctx, id)
	if err != nil {
		th.logger.Error("Failed to get technology %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	if technologyEntity == nil {
		th.logger.Error("Technology %d not found", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Technology", idStr))
		return
	}

	response := technologyDto.FromTechnologyEntityToResponse(technologyEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// CreateTechnology
//
//	@Summary		Create a new technology
//	@Description	Create a new technology for the authenticated admin user
//	@Tags			Admin Technologies
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateTechnologyRequest	true	"Technology creation request"
//	@Success		201		{object}	shared.APIResponse{data=dto.TechnologyResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/technologies [post]
//	@Security		BearerAuth
func (th *technologyHandler) CreateTechnology(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request technologyDto.CreateTechnologyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		th.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	userID, ok := th.getUserIDFromContext(w, r)
	if !ok {
		th.logger.Error("Failed to get user ID from context")
		return
	}

	technologyEntity, err := request.ToEntity(userID)
	if err != nil {
		th.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology data", "technology", &err))
		return
	}

	createdTechnology, err := th.technologyUseCase.CreateTechnology(ctx, technologyEntity)
	if err != nil {
		th.logger.Error("Failed to create technology: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := technologyDto.FromTechnologyEntityToResponse(createdTechnology, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusCreated, response)
}

// CreateBulkTechnologies
//
//	@Summary		Create multiple technologies in bulk
//	@Description	Create multiple technologies at once for the authenticated admin user
//	@Tags			Admin Technologies
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateBulkTechnologiesRequest	true	"Bulk technology creation request"
//	@Success		201		{object}	shared.APIResponse{data=dto.TechnologyListResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/technologies/bulk [post]
//	@Security		BearerAuth
func (th *technologyHandler) CreateBulkTechnologies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request technologyDto.CreateBulkTechnologiesRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		th.logger.Error("Failed to decode bulk technologies request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		th.logger.Error("Invalid bulk technologies request: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	userID, ok := th.getUserIDFromContext(w, r)
	if !ok {
		th.logger.Error("Failed to get user ID from context")
		return
	}

	technologyEntities, err := request.ToEntities(userID)
	if err != nil {
		th.logger.Error("Failed to convert bulk request to entities: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technologies data", "technologies", &err))
		return
	}

	var createdTechnologies []*entities.Technology
	var errs []error

	for i, technologyEntity := range technologyEntities {
		createdTechnology, err := th.technologyUseCase.CreateTechnology(ctx, technologyEntity)
		if err != nil {
			th.logger.Error("Failed to create technology at index %d (name: %s): %v", i, technologyEntity.Name, err)
			errs = append(errs, err)
		} else {
			createdTechnologies = append(createdTechnologies, createdTechnology)
		}
	}

	statusCode := http.StatusCreated
	if len(technologyEntities) == len(errs) {
		th.logger.Error("All technologies failed to create, returning errors")
		utils.WriteErrorResponse(w, errs...)
		return
	} else if len(errs) > 0 {
		statusCode = http.StatusMultiStatus
	}

	response := technologyDto.FromTechnologiesEntityForBulkToResponse(createdTechnologies, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	domainErrors := make([]*shared.APIError, len(errs))
	for i, err := range errs {
		if domainErr, ok := domain.AsDomainError(err); ok {
			domainErrors[i] = utils.DomainErrorToAPIError(domainErr)
		}
	}
	response.Errors = domainErrors
	utils.WriteSuccessResponse(w, statusCode, response)
}

// UpdateTechnology
//
//	@Summary		Update an existing technology
//	@Description	Update an existing technology by ID for the authenticated admin user
//	@Tags			Admin Technologies
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int										true	"Technology ID"
//	@Param			request	body		dto.UpdateTechnologyRequest	true	"Technology update request"
//	@Success		200		{object}	shared.APIResponse{data=dto.TechnologyResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/technologies/{id} [put]
//	@Security		BearerAuth
func (th *technologyHandler) UpdateTechnology(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		th.logger.Error("Technology ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Technology ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		th.logger.Error("Invalid technology ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology ID", "id", &err))
		return
	}

	var request technologyDto.UpdateTechnologyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		th.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	request.ID = id

	userID, ok := th.getUserIDFromContext(w, r)
	if !ok {
		th.logger.Error("Failed to get user ID from context")
		return
	}

	technologyEntity, err := request.ToEntity(userID)
	if err != nil {
		th.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology data", "technology", &err))
		return
	}

	updatedTechnology, err := th.technologyUseCase.UpdateTechnology(ctx, id, technologyEntity)
	if err != nil {
		th.logger.Error("Failed to update technology %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := technologyDto.FromTechnologyEntityToResponse(updatedTechnology, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// PatchTechnology
//
//	@Summary		Partially update a technology
//	@Description	Partially update a technology by ID for the authenticated admin user
//	@Tags			Admin Technologies
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int								true	"Technology ID"
//	@Param			technology	body		dto.PatchTechnologyRequest	true	"Technology data to update"
//	@Success		200			{object}	dto.TechnologyResponse
//	@Failure		400			{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401			{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404			{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500			{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/technologies/{id} [patch]
//	@Security		BearerAuth
func (th *technologyHandler) PatchTechnology(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		th.logger.Error("Technology ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Technology ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		th.logger.Error("Invalid technology ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology ID", "id", &err))
		return
	}

	var request technologyDto.PatchTechnologyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		th.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		th.logger.Error("Invalid request: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	userID, ok := th.getUserIDFromContext(w, r)
	if !ok {
		th.logger.Error("Failed to get user ID from context")
		return
	}

	technologyEntity, err := request.ToEntity(id, userID)
	if err != nil {
		th.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology data", "technology", &err))
		return
	}

	patchedTechnologyEntity, err := th.technologyUseCase.PatchTechnology(ctx, id, technologyEntity)
	if err != nil {
		th.logger.Error("Failed to patch technology %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := technologyDto.FromTechnologyEntityToResponse(patchedTechnologyEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// DeleteTechnology
//
//	@Summary		Delete a technology
//	@Description	Delete a technology by ID for the authenticated admin user
//	@Tags			Admin Technologies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Technology ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/technologies/{id} [delete]
//	@Security		BearerAuth
func (th *technologyHandler) DeleteTechnology(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		th.logger.Error("Technology ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Technology ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		th.logger.Error("Invalid technology ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology ID", "id", &err))
		return
	}

	err = th.technologyUseCase.DeleteTechnology(ctx, id)
	if err != nil {
		th.logger.Error("Failed to delete technology %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
