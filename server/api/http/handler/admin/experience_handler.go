package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
	experienceDto "portfolio/dto/experience"
	"portfolio/logger"
	"portfolio/shared"
	"strconv"
	"time"
)

type experienceHandler struct {
	AbstractHandler
	experienceUseCase *usecases.ExperienceUseCase
	logger            *logger.Logger
}

func NewExperienceHandler(settingUseCase *usecases.SettingUseCase, experienceUseCase *usecases.ExperienceUseCase, logger *logger.Logger) []*routes.NamedRoute {
	experienceHandler := experienceHandler{
		AbstractHandler:   AbstractHandler{settingUseCase: settingUseCase},
		experienceUseCase: experienceUseCase,
		logger:            logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetAdminExperiencesHandler",
			Pattern: "GET /experiences",
			Handler: experienceHandler.GetExperiences,
		},
		{
			Name:    "PostAdminExperienceHandler",
			Pattern: "POST /experiences",
			Handler: experienceHandler.CreateExperience,
		},
		{
			Name:    "GetAdminExperienceHandler",
			Pattern: "GET /experiences/{id}",
			Handler: experienceHandler.GetExperience,
		},
		{
			Name:    "PutAdminExperienceHandler",
			Pattern: "PUT /experiences/{id}",
			Handler: experienceHandler.UpdateExperience,
		},
		{
			Name:    "PatchAdminExperienceHandler",
			Pattern: "PATCH /experiences/{id}",
			Handler: experienceHandler.PatchExperience,
		},
		{
			Name:    "DeleteAdminExperienceHandler",
			Pattern: "DELETE /experiences/{id}",
			Handler: experienceHandler.DeleteExperience,
		},
	}
}

// GetExperiences
//
//	@Summary		Get all experiences
//	@Description	Get a list of all experiences for the authenticated admin user
//	@Tags			Admin Experiences
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.ExperienceListResponse}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/experiences [get]
//	@Security		BearerAuth
func (eh *experienceHandler) GetExperiences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := eh.getUserIDFromContext(w, r)
	if !ok {
		eh.logger.Error("Failed to get user ID from context")
		return
	}

	experiences, err := eh.experienceUseCase.GetExperiencesByUserID(ctx, userID)
	if err != nil {
		eh.logger.Error("Failed to get experiences for user %d: %v", userID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := experienceDto.FromExperiencesEntityToResponse(experiences,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})

	if response == nil {
		eh.logger.Error("No experiences found for user %d", userID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Experiences", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetExperience
//
//	@Summary		Get an experience by ID
//	@Description	Get a specific experience by ID for the authenticated admin user
//	@Tags			Admin Experiences
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Experience ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.ExperienceResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/experiences/{id} [get]
//	@Security		BearerAuth
func (eh *experienceHandler) GetExperience(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Experience ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Experience ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid experience ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience ID", "id", &err))
		return
	}

	experienceEntity, err := eh.experienceUseCase.GetExperienceByID(ctx, id)
	if err != nil {
		eh.logger.Error("Failed to get experience %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	if experienceEntity == nil {
		eh.logger.Error("Experience not found: %d", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Experience", idStr))
		return
	}

	response := experienceDto.FromExperienceEntityToResponse(experienceEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// CreateExperience
//
//	@Summary		Create a new experience
//	@Description	Create a new experience for the authenticated admin user
//	@Tags			Admin Experiences
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateExperienceRequest	true	"Experience creation request"
//	@Success		201		{object}	shared.APIResponse{data=dto.ExperienceResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/experiences [post]
//	@Security		BearerAuth
func (eh *experienceHandler) CreateExperience(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request experienceDto.CreateExperienceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		eh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		eh.logger.Error("Invalid request: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	userID, ok := eh.getUserIDFromContext(w, r)
	if !ok {
		eh.logger.Error("Failed to get user ID from context")
		return
	}

	experienceEntity, err := request.ToEntity(userID)
	if err != nil {
		eh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience data", "experience", &err))
		return
	}

	fmt.Printf("_____________________________%+v\n", experienceEntity)

	createdExperience, err := eh.experienceUseCase.CreateExperience(ctx, experienceEntity)
	if err != nil {
		eh.logger.Error("Failed to create experience: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := experienceDto.FromExperienceEntityToResponse(createdExperience, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusCreated, response)
}

// UpdateExperience
//
//	@Summary		Update an existing experience
//	@Description	Update an existing experience by ID for the authenticated admin user
//	@Tags			Admin Experiences
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Experience ID"
//	@Param			request	body		dto.UpdateExperienceRequest	true	"Experience update request"
//	@Success		200		{object}	shared.APIResponse{data=dto.ExperienceResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/experiences/{id} [put]
//	@Security		BearerAuth
func (eh *experienceHandler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Experience ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Experience ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid experience ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience ID", "id", &err))
		return
	}

	var request experienceDto.UpdateExperienceRequest
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

	experienceEntity, err := request.ToEntity(id, userID)
	if err != nil {
		eh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience data", "experience", &err))
		return
	}

	updatedExperience, err := eh.experienceUseCase.UpdateExperience(ctx, id, experienceEntity)
	if err != nil {
		eh.logger.Error("Failed to update experience %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := experienceDto.FromExperienceEntityToResponse(updatedExperience, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// PatchExperience
//
//	@Summary		Partially update an experience
//	@Description	Partially update an experience by ID for the authenticated admin user
//	@Tags			Admin Experiences
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Experience ID"
//	@Param			request	body		dto.PatchExperienceRequest	true	"Patch experience request"
//	@Success		200		{object}	shared.APIResponse{data=dto.ExperienceResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/experiences/{id} [patch]
//	@Security		BearerAuth
func (eh *experienceHandler) PatchExperience(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Experience ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Experience ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid experience ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience ID", "id", &err))
		return
	}

	userID, ok := eh.getUserIDFromContext(w, r)
	if !ok {
		eh.logger.Error("Failed to get user ID from context")
		return
	}

	var request experienceDto.PatchExperienceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		eh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	experienceEntity, err := request.ToEntity(id, userID)
	if err != nil {
		eh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience data", "experience", &err))
		return
	}

	updatedExperience, err := eh.experienceUseCase.PatchExperience(ctx, id, experienceEntity)
	if err != nil {
		eh.logger.Error("Failed to patch experience %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := experienceDto.FromExperienceEntityToResponse(updatedExperience, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// DeleteExperience
//
//	@Summary		Delete an experience
//	@Description	Delete an experience by ID for the authenticated admin user
//	@Tags			Admin Experiences
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Experience ID"
//	@Success		204	"No Content"
//	@Router			/admin/experiences/{id} [delete]
//	@Security		BearerAuth
func (eh *experienceHandler) DeleteExperience(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		eh.logger.Error("Experience ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Experience ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		eh.logger.Error("Invalid experience ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience ID", "id", &err))
		return
	}

	err = eh.experienceUseCase.DeleteExperience(ctx, id)
	if err != nil {
		eh.logger.Error("Failed to delete experience %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
