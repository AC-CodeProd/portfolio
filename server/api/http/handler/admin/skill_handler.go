package admin

import (
	"encoding/json"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/usecases"
	skillDto "portfolio/dto/skill"
	"portfolio/logger"
	"portfolio/shared"
	"strconv"
	"time"
)

type skillHandler struct {
	AbstractHandler
	skillUseCase *usecases.SkillUseCase
	logger       *logger.Logger
}

func NewSkillHandler(settingUseCase *usecases.SettingUseCase, skillUseCase *usecases.SkillUseCase, logger *logger.Logger) []*routes.NamedRoute {
	skillHandler := skillHandler{
		AbstractHandler: AbstractHandler{settingUseCase: settingUseCase},
		skillUseCase:    skillUseCase,
		logger:          logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetAdminSkillsHandler",
			Pattern: "GET /skills",
			Handler: skillHandler.GetSkills,
		},
		{
			Name:    "PostAdminSkillHandler",
			Pattern: "POST /skills",
			Handler: skillHandler.CreateSkill,
		},
		{
			Name:    "PostBulkAdminSkillHandler",
			Pattern: "POST /skills/bulk",
			Handler: skillHandler.CreateBulkSkills,
		},
		{
			Name:    "GetAdminSkillHandler",
			Pattern: "GET /skills/{id}",
			Handler: skillHandler.GetSkill,
		},
		{
			Name:    "PutAdminSkillHandler",
			Pattern: "PUT /skills/{id}",
			Handler: skillHandler.UpdateSkill,
		},
		{
			Name:    "PatchAdminSkillHandler",
			Pattern: "PATCH /skills/{id}",
			Handler: skillHandler.PatchSkill,
		},
		{
			Name:    "DeleteAdminSkillHandler",
			Pattern: "DELETE /skills/{id}",
			Handler: skillHandler.DeleteSkill,
		},
	}
}

// GetSkills
//
//	@Summary		Get all admin skills
//	@Description	Retrieve all skills for the authenticated admin user
//	@Tags			Admin Skills
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.SkillListResponse}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/skills [get]
//	@Security		BearerAuth
func (sh *skillHandler) GetSkills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := sh.getUserIDFromContext(w, r)
	if !ok {
		sh.logger.Error("Failed to get user ID from context")
		return
	}

	skills, err := sh.skillUseCase.GetSkillsByUserID(ctx, userID)
	if err != nil {
		sh.logger.Error("Failed to get skills for user %d: %v", userID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := skillDto.FromSkillsEntityToResponse(skills,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})

	if response == nil {
		sh.logger.Warn("No skills found for user %d", userID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Skills", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetSkill
//
//	@Summary		Get a specific admin skill
//	@Description	Retrieve a specific skill by ID for admin management
//	@Tags			Admin Skills
//	@Produce		json
//	@Param			id	path		int	true	"Skill ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.SkillResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/skills/{id} [get]
//	@Security		BearerAuth
func (sh *skillHandler) GetSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		sh.logger.Error("Skill ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Skill ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sh.logger.Error("Invalid skill ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill ID", "id", &err))
		return
	}

	skillEntity, err := sh.skillUseCase.GetSkillByID(ctx, id)
	if err != nil {
		sh.logger.Error("Failed to get skill %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	if skillEntity == nil {
		sh.logger.Error("Skill not found: %d", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Skill", idStr))
		return
	}

	response := skillDto.FromSkillEntityToResponse(skillEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if response == nil {
		sh.logger.Error("No skill found: %d", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Skill", idStr))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// CreateSkill
//
//	@Summary		Create a new skill
//	@Description	Create a new skill for the authenticated admin user
//	@Tags			Admin Skills
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateSkillRequest	true	"Skill creation request"
//	@Success		201		{object}	shared.APIResponse{data=dto.SkillResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/skills [post]
//	@Security		BearerAuth
func (h *skillHandler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request skillDto.CreateSkillRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	userID, ok := h.getUserIDFromContext(w, r)
	if !ok {
		return
	}

	skillEntity, err := request.ToEntity(userID)
	if err != nil {
		h.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill data", "skill", &err))
		return
	}

	createdSkill, err := h.skillUseCase.CreateSkill(ctx, skillEntity)
	if err != nil {
		h.logger.Error("Failed to create skill: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := skillDto.FromSkillEntityToResponse(createdSkill, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	utils.WriteSuccessResponse(w, http.StatusCreated, response)
}

// CreateBulkSkills
//
//	@Summary		Create multiple skills in bulk
//	@Description	Create multiple skills for the authenticated admin user
//	@Tags			Admin Skills
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateBulkSkillsRequest	true	"Bulk skills creation request"
//	@Success		201		{object}	shared.APIResponse{data=dto.SkillListResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/skills/bulk [post]
//	@Security		BearerAuth
func (sh *skillHandler) CreateBulkSkills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request skillDto.CreateBulkSkillsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sh.logger.Error("Failed to decode bulk skills request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	userID, ok := sh.getUserIDFromContext(w, r)
	if !ok {
		return
	}

	skillEntities, err := request.ToEntities(userID)
	if err != nil {
		sh.logger.Error("Failed to convert bulk request to entities: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skills data", "skills", &err))
		return
	}

	var createdSkills []*entities.Skill
	var errs []error

	for i, skillEntity := range skillEntities {
		createdSkill, err := sh.skillUseCase.CreateSkill(ctx, skillEntity)
		if err != nil {
			sh.logger.Error("Failed to create skill at index %d (name: %s): %v", i, skillEntity.Name, err)
			errs = append(errs, err)
		} else {
			createdSkills = append(createdSkills, createdSkill)
		}
	}

	statusCode := http.StatusCreated
	if len(skillEntities) == len(errs) {
		sh.logger.Error("All skills failed to create, returning errors")
		utils.WriteErrorResponse(w, errs...)
		return
	} else if len(errs) > 0 {
		statusCode = http.StatusMultiStatus
	}

	response := skillDto.FromSkillsEntityForBulkToResponse(createdSkills, &shared.Meta{
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

// UpdateSkill
//
//	@Summary		Update an existing skill
//	@Description	Update an existing skill by ID for the authenticated admin user
//	@Tags			Admin Skills
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Skill ID"
//	@Param			request	body		dto.UpdateSkillRequest	true	"Skill update request"
//	@Success		200		{object}	shared.APIResponse{data=dto.SkillResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/skills/{id} [put]
//	@Security		BearerAuth
func (sh *skillHandler) UpdateSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		sh.logger.Error("Skill ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Skill ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sh.logger.Error("Invalid skill ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill ID", "id", &err))
		return
	}

	var request skillDto.UpdateSkillRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	userID, ok := sh.getUserIDFromContext(w, r)
	if !ok {
		sh.logger.Error("Failed to get user ID from context")
		return
	}

	skillEntity, err := request.ToEntity(id, userID)
	if err != nil {
		sh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill data", "skill", &err))
		return
	}

	updatedSkill, err := sh.skillUseCase.UpdateSkill(ctx, id, skillEntity)
	if err != nil {
		sh.logger.Error("Failed to update skill %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := skillDto.FromSkillEntityToResponse(updatedSkill, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// PatchSkill
//
//	@Summary		Partially update a skill
//	@Description	Partially update a skill by ID for the authenticated admin user
//	@Tags			Admin Skills
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Skill ID"
//	@Param			skill	body		dto.PatchSkillRequest		true	"Skill data to update"
//	@Success		200		{object}	dto.SkillResponse
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/skills/{id} [patch]
//	@Security		BearerAuth
func (sh *skillHandler) PatchSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		sh.logger.Error("Skill ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Skill ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sh.logger.Error("Invalid skill ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill ID", "id", &err))
		return
	}

	var request skillDto.PatchSkillRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := request.Validate(); err != nil {
		sh.logger.Error("Invalid request: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("request", err.Error(), nil))
		return
	}

	userID, ok := sh.getUserIDFromContext(w, r)
	if !ok {
		sh.logger.Error("Failed to get user ID from context")
		return
	}

	skillEntity, err := request.ToEntity(id, userID)
	if err != nil {
		sh.logger.Error("Failed to convert request to entity: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill data", "skill", &err))
		return
	}

	patchedSkill, err := sh.skillUseCase.PatchSkill(ctx, id, skillEntity)
	if err != nil {
		sh.logger.Error("Failed to patch skill %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := skillDto.FromSkillEntityToResponse(patchedSkill, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// DeleteSkill
//
//	@Summary		Delete a skill
//	@Description	Delete a skill by ID for the authenticated admin user
//	@Tags			Admin Skills
//	@Produce		json
//	@Param			id	path		int	true	"Skill ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/skills/{id} [delete]
//	@Security		BearerAuth
func (sh *skillHandler) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		sh.logger.Error("Skill ID is required")
		utils.WriteErrorResponse(w, domain.NewValidationError("Skill ID is required", "id", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sh.logger.Error("Invalid skill ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid skill ID", "id", &err))
		return
	}

	err = sh.skillUseCase.DeleteSkill(ctx, id)
	if err != nil {
		sh.logger.Error("Failed to delete skill %d: %v", id, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
