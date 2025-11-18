package admin

import (
	"encoding/json"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
	educationDto "portfolio/dto/education"
	"portfolio/logger"
	"portfolio/shared"
	"strconv"
	"time"
)

type educationHandler struct {
	AbstractHandler
	educationUseCase *usecases.EducationUseCase
	logger           *logger.Logger
}

func NewEducationHandler(settingUseCase *usecases.SettingUseCase, educationUseCase *usecases.EducationUseCase, logger *logger.Logger) []*routes.NamedRoute {
	educationHandler := educationHandler{
		AbstractHandler:  AbstractHandler{settingUseCase: settingUseCase},
		educationUseCase: educationUseCase,
		logger:           logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetAdminEducationsHandler",
			Pattern: "GET /educations",
			Handler: educationHandler.GetEducations,
		},
		{
			Name:    "PostAdminEducationHandler",
			Pattern: "POST /educations",
			Handler: educationHandler.CreateEducation,
		},
		{
			Name:    "GetAdminEducationHandler",
			Pattern: "GET /educations/{id}",
			Handler: educationHandler.GetEducation,
		},
		{
			Name:    "PutAdminEducationHandler",
			Pattern: "PUT /educations/{id}",
			Handler: educationHandler.UpdateEducation,
		},
		{
			Name:    "PatchAdminEducationHandler",
			Pattern: "PATCH /educations/{id}",
			Handler: educationHandler.PatchEducation,
		},
		{
			Name:    "DeleteAdminEducationHandler",
			Pattern: "DELETE /educations/{id}",
			Handler: educationHandler.DeleteEducation,
		},
	}
}

// GetEducations
//
//	@Summary		Get all admin educations
//	@Description	Retrieve all educations for the authenticated admin user
//	@Tags			Admin Educations
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	shared.APIResponse{data=dto.EducationListResponse}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}	"Unauthorized"
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/educations [get]
func (eh *educationHandler) GetEducations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := eh.getUserIDFromContext(w, r)
	if !ok {
		eh.logger.Error("Failed to get user ID from context")
		return
	}

	educations, err := eh.educationUseCase.GetEducationsByUserID(ctx, userID)
	if err != nil {
		eh.logger.Error("Failed to get educations for user %d: %v", userID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := educationDto.FromEducationsEntityToResponse(educations,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})

	if response == nil {
		eh.logger.Error("No educations found for user %d", userID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Educations", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetEducation
//
//	@Summary		Get a specific admin education
//	@Description	Retrieve a specific education by ID for admin management
//	@Tags			Admin Educations
//	@Produce		json
//	@Param			id	path	int	true	"Education ID"
//	@Security		BearerAuth
//	@Success		200	{object}	shared.APIResponse{data=dto.EducationResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/educations/{id} [get]
func (eh *educationHandler) GetEducation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Education ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Education ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid education ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education ID", "id", &err))
		return
	}

	educationEntity, err := eh.educationUseCase.GetEducationByID(ctx, id)
	if err != nil {
		eh.logger.Error("Failed to get education %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	if educationEntity == nil {
		eh.logger.Error("Education not found: %d", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Education", idStr))
		return
	}

	response := educationDto.FromEducationEntityToResponse(educationEntity,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// CreateEducation
//
//	@Summary		Create a new education
//	@Description	Create a new education for the authenticated admin user
//	@Tags			Admin Educations
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.CreateEducationRequest	true	"Education creation request"
//	@Security		BearerAuth
//	@Success		201	{object}	shared.APIResponse{data=dto.EducationResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/educations [post]
func (eh *educationHandler) CreateEducation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request educationDto.CreateEducationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		eh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	userID, ok := eh.getUserIDFromContext(w, r)
	if !ok {
		eh.logger.Error("Failed to get user ID from context")
		return
	}

	educationEntity, err := request.ToEntity(userID)
	if err != nil {
		eh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education data", "education", &err))
		return
	}

	createdEducation, err := eh.educationUseCase.CreateEducation(ctx, educationEntity)
	if err != nil {
		eh.logger.Error("Failed to create education: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := educationDto.FromEducationEntityToResponse(createdEducation,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})
	utils.WriteSuccessResponse(w, http.StatusCreated, response)
}

// UpdateEducation
//
//	@Summary		Update an existing education
//	@Description	Update an existing education by ID for the authenticated admin user
//	@Tags			Admin Educations
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int							true	"Education ID"
//	@Param			request	body	dto.UpdateEducationRequest	true	"Education update request"
//	@Security		BearerAuth
//	@Success		200	{object}	shared.APIResponse{data=dto.EducationResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/educations/{id} [put]
func (eh *educationHandler) UpdateEducation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Education ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Education ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid education ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education ID", "id", &err))
		return
	}

	var request educationDto.UpdateEducationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		eh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	request.ID = id

	userID, ok := eh.getUserIDFromContext(w, r)
	if !ok {
		eh.logger.Error("Failed to get user ID from context")
		return
	}

	educationEntity, err := request.ToEntity(userID)
	if err != nil {
		eh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education data", "education", &err))
		return
	}

	updatedEducation, err := eh.educationUseCase.UpdateEducation(ctx, id, educationEntity)
	if err != nil {
		eh.logger.Error("Failed to update education %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := educationDto.FromEducationEntityToResponse(updatedEducation,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// PatchEducation
//
//	@Summary		Partially update an education
//	@Description	Partially update an education by ID for the authenticated admin user
//	@Tags			Admin Educations
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Education ID"
//	@Param			request	body		dto.PatchEducationRequest	true	"Patch education request"
//	@Success		200		{object}	shared.APIResponse{data=dto.EducationResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/educations/{id} [patch]
//	@Security		BearerAuth
func (eh *educationHandler) PatchEducation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Education ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Education ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid education ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education ID", "id", &err))
		return
	}

	userID, ok := eh.getUserIDFromContext(w, r)
	if !ok {
		eh.logger.Error("Failed to get user ID from context")
		return
	}

	var request educationDto.PatchEducationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		eh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		eh.logger.Error("Validation error: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	educationEntity, err := request.ToEntity(id, userID)
	if err != nil {
		eh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education data", "education", &err))
		return
	}

	updatedEducation, err := eh.educationUseCase.PatchEducation(ctx, id, educationEntity)
	if err != nil {
		eh.logger.Error("Failed to patch education %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := educationDto.FromEducationEntityToResponse(updatedEducation,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// DeleteEducation
//
//	@Summary		Delete an education
//	@Description	Delete an education by ID for the authenticated admin user
//	@Tags			Admin Educations
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Education ID"
//	@Security		BearerAuth
//	@Success		204	"No Content"
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/educations/{id} [delete]
func (eh *educationHandler) DeleteEducation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Education ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Education ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid education ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education ID", "id", &err))
		return
	}

	err = eh.educationUseCase.DeleteEducation(ctx, id)
	if err != nil {
		eh.logger.Error("Failed to delete education %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
