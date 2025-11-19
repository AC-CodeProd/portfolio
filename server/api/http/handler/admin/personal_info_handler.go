package admin

import (
	"encoding/json"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/usecases"
	personalInfoDto "portfolio/dto/personal_info"
	"portfolio/logger"
	"portfolio/shared"
	"strconv"
	"time"
)

type personalInfoHandler struct {
	AbstractHandler
	personalInfoUseCase *usecases.PersonalInfoUseCase
	logger              *logger.Logger
}

func NewPersonalInfoHandler(settingUseCase *usecases.SettingUseCase, personalInfoUseCase *usecases.PersonalInfoUseCase, logger *logger.Logger) []*routes.NamedRoute {
	personalInfoHandler := personalInfoHandler{
		AbstractHandler: AbstractHandler{
			settingUseCase: settingUseCase,
		},
		personalInfoUseCase: personalInfoUseCase,
		logger:              logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetAdminPersonalInfoHandler",
			Pattern: "GET /personal-info",
			Handler: personalInfoHandler.GetPersonalInfo,
		},
		{
			Name:    "PostAdminPersonalInfoHandler",
			Pattern: "POST /personal-info",
			Handler: personalInfoHandler.CreatePersonalInfo,
		},
		{
			Name:    "PutAdminPersonalInfoHandler",
			Pattern: "PUT /personal-info/{id}",
			Handler: personalInfoHandler.UpdatePersonalInfo,
		},
		{
			Name:    "PatchAdminPersonalInfoHandler",
			Pattern: "PATCH /personal-info/{id}",
			Handler: personalInfoHandler.PatchPersonalInfo,
		},
		{
			Name:    "DeleteAdminPersonalInfoHandler",
			Pattern: "DELETE /personal-info/{id}",
			Handler: personalInfoHandler.DeletePersonalInfo,
		},
	}
}

func (pih *personalInfoHandler) getUserFromContext(w http.ResponseWriter, r *http.Request) *entities.User {

	userID, ok := pih.getUserIDFromContext(w, r)
	if !ok {
		pih.logger.Error("Failed to get user ID from context")
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid user ID in context", "user_id", nil))
		return nil
	}

	user, err := pih.personalInfoUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		pih.logger.Error("Failed to get user by ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("user", "ID"))
		return nil
	}
	return user
}

// GetPersonalInfo
//
//	@Summary		Get admin personal information
//	@Description	Retrieve personal information for admin management
//	@Tags			Admin Personal Info
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.PersonalInfoResponse}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/personal-info [get]
//	@Security		BearerAuth
func (pih *personalInfoHandler) GetPersonalInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	portfolioOwnerID, err := utils.GetPortfolioOwnerID(pih.settingUseCase, ctx, w)

	if err != nil {
		pih.logger.Error("Error retrieving portfolio owner ID: %v", err)
		return
	}
	personalInfo, err := pih.personalInfoUseCase.GetPersonalInfoByUserID(ctx, portfolioOwnerID)
	if err != nil {
		pih.logger.Error("Error retrieving personal information: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Error retrieving personal information", err))
		return
	}
	if personalInfo == nil {
		pih.logger.Error("Personal information not found for user ID: %d", portfolioOwnerID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("personal information", "user"))
		return
	}

	user := pih.getUserFromContext(w, r)
	if user == nil {
		return
	}

	personalInfoResponse := personalInfoDto.FromPersonalInfoEntityToResponse(user, personalInfo, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusOK, personalInfoResponse)
}

// CreatePersonalInfo
//
//	@Summary		Create personal information
//	@Description	Create new personal information for admin user
//	@Tags			Admin Personal Info
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreatePersonalInfoRequest	true	"Personal info creation request"
//	@Success		201		{object}	shared.APIResponse{data=dto.PersonalInfoResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/personal-info [post]
//	@Security		BearerAuth
func (pih *personalInfoHandler) CreatePersonalInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request personalInfoDto.CreatePersonalInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		pih.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := pih.personalInfoUseCase.ValidateCreatePersonalInfoRequest(&request); err != nil {
		pih.logger.Error("Validation error: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid personal information", "body", &err))
		return
	}

	user := pih.getUserFromContext(w, r)
	if user == nil {
		pih.logger.Error("Failed to get user from context")
		return
	}

	personalInfo, err := pih.personalInfoUseCase.CreatePersonalInfo(ctx, user.ID, &request)
	if err != nil {
		pih.logger.Error("Failed to create personal information: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to create personal information", err))
		return
	}

	personalInfoResponse := personalInfoDto.FromPersonalInfoEntityToResponse(user, personalInfo, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	utils.WriteSuccessResponse(w, http.StatusCreated, personalInfoResponse)
}

// UpdatePersonalInfo
//
//	@Summary		Update personal information
//	@Description	Update existing personal information by ID
//	@Tags			Admin Personal Info
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Personal Info ID"
//	@Param			request	body		dto.UpdatePersonalInfoRequest	true	"Personal info update request"
//	@Success		200		{object}	shared.APIResponse{data=dto.PersonalInfoResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/personal-info/{id} [put]
//	@Security		BearerAuth
func (pih *personalInfoHandler) UpdatePersonalInfo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pih.logger.Error("Invalid ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewInvalidFormatError("id", "integer"))
		return
	}
	ctx := r.Context()

	userID, ok := pih.getUserIDFromContext(w, r)
	if !ok {
		return
	}

	var request personalInfoDto.UpdatePersonalInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		pih.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	if err := pih.personalInfoUseCase.ValidateUpdatePersonalInfoRequest(&request); err != nil {
		pih.logger.Error("Validation error: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid personal information", "body", &err))
		return
	}

	personalInfo, err := pih.personalInfoUseCase.UpdatePersonalInfo(ctx, id, &request)
	if err != nil {
		pih.logger.Error("Failed to update personal info: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to update personal information", err))
		return
	}

	if personalInfo == nil {
		pih.logger.Error("Personal information not found: %d", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("personal information", "ID"))
		return
	}

	user, err := pih.personalInfoUseCase.GetUserByID(ctx, userID)
	if err != nil {
		pih.logger.Error("Failed to get user by ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("user", "ID"))
		return
	}
	personalInfoResponse := personalInfoDto.FromPersonalInfoEntityToResponse(user, personalInfo, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	utils.WriteSuccessResponse(w, http.StatusOK, personalInfoResponse)
}

// PatchPersonalInfo
//
//	@Summary		Partially update personal information
//	@Description	Partially update existing personal information by ID
//	@Tags			Admin Personal Info
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Personal Info ID"
//	@Param			request	body		dto.PatchPersonalInfoRequest	true	"Personal info patch request"
//	@Success		200		{object}	shared.APIResponse{data=dto.PersonalInfoResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/personal-info/{id} [patch]
//	@Security		BearerAuth
func (pih *personalInfoHandler) PatchPersonalInfo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pih.logger.Error("Invalid ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewInvalidFormatError("id", "integer"))
		return
	}
	ctx := r.Context()

	userID, ok := pih.getUserIDFromContext(w, r)
	if !ok {
		pih.logger.Error("Failed to get user ID from context")
		return
	}

	var request personalInfoDto.PatchPersonalInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		pih.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid request body", "body", &err))
		return
	}

	personalInfo, err := pih.personalInfoUseCase.PatchPersonalInfo(ctx, id, &request)
	if err != nil {
		pih.logger.Error("Failed to patch personal info: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to patch personal information", err))
		return
	}

	if personalInfo == nil {
		pih.logger.Error("Personal information not found: %d", id)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("personal information", "ID"))
		return
	}

	user, err := pih.personalInfoUseCase.GetUserByID(ctx, userID)
	if err != nil {
		pih.logger.Error("Failed to get user by ID: %v", err)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("user", "ID"))
		return
	}
	personalInfoResponse := personalInfoDto.FromPersonalInfoEntityToResponse(user, personalInfo, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	utils.WriteSuccessResponse(w, http.StatusOK, personalInfoResponse)
}

// DeletePersonalInfo
//
//	@Summary		Delete personal information
//	@Description	Delete personal information by ID
//	@Tags			Admin Personal Info
//	@Produce		json
//	@Param			id	path	int	true	"Personal Info ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		401	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/personal-info/{id} [delete]
//	@Security		BearerAuth
func (pih *personalInfoHandler) DeletePersonalInfo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pih.logger.Error("Invalid ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewInvalidFormatError("id", "integer"))
		return
	}
	ctx := r.Context()

	if err := pih.personalInfoUseCase.DeletePersonalInfo(ctx, id); err != nil {
		pih.logger.Error("Failed to delete personal info: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to delete personal information", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
