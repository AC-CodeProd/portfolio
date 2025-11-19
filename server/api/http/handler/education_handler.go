package handler

import (
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
		AbstractHandler: AbstractHandler{
			settingUseCase: settingUseCase,
		},
		educationUseCase: educationUseCase,
		logger:           logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetEducationsHandler",
			Pattern: "GET /educations",
			Handler: educationHandler.GetEducations,
		},
		{
			Name:    "GetEducationHandler",
			Pattern: "GET /educations/{id}",
			Handler: educationHandler.GetEducation,
		},
	}
}

// GetEducations
//
//	@Summary		Get all educations
//	@Description	Retrieve all educations for the portfolio
//	@Tags			Educations
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.EducationListResponse}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/educations [get]
func (eh *educationHandler) GetEducations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(eh.settingUseCase, ctx, w)

	if err != nil {
		eh.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	educationsEntity, err := eh.educationUseCase.GetEducationsByUserID(ctx, portfolioOwnerID)
	if err != nil {
		eh.logger.Error("Failed to get educations: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := educationDto.FromEducationsEntityToResponse(educationsEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})
	if response == nil {
		eh.logger.Error("No educations found for user ID %d", portfolioOwnerID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Educations", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetEducation
//
//	@Summary		Get a specific education
//	@Description	Retrieve a specific education by ID
//	@Tags			Educations
//	@Produce		json
//	@Param			id	path		int	true	"Education ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.EducationResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/educations/{id} [get]
func (eh *educationHandler) GetEducation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	educationIDStr := r.PathValue("id")
	educationID, err := strconv.Atoi(educationIDStr)
	if err != nil || educationID <= 0 {
		eh.logger.Error("Invalid education ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid education ID", "id", &err))
		return
	}

	educationEntity, err := eh.educationUseCase.GetEducationByID(ctx, educationID)
	if err != nil {
		eh.logger.Error("Failed to get education: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}
	portfolioOwnerID, err := utils.GetPortfolioOwnerID(eh.settingUseCase, ctx, w)

	if err != nil {
		eh.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	if educationEntity.UserID != portfolioOwnerID {
		eh.logger.Error("Unauthorized access to education %d", educationID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Education", strconv.Itoa(educationID)))
		return
	}

	response := educationDto.FromEducationEntityToResponse(educationEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}
