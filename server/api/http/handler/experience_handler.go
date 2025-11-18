package handler

import (
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
		AbstractHandler: AbstractHandler{
			settingUseCase: settingUseCase,
		},
		experienceUseCase: experienceUseCase,
		logger:            logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetExperiencesHandler",
			Pattern: "GET /experiences",
			Handler: experienceHandler.GetExperiences,
		},
		{
			Name:    "GetExperienceHandler",
			Pattern: "GET /experiences/{id}",
			Handler: experienceHandler.GetExperience,
		},
	}
}

// GetExperiences
//
//	@Summary		Get all experiences
//	@Description	Retrieve all experiences for the portfolio
//	@Tags			Experiences
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.ExperienceListResponse}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/experiences [get]
func (eh *experienceHandler) GetExperiences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(eh.settingUseCase, ctx, w)

	if err != nil {
		eh.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	experiences, err := eh.experienceUseCase.GetExperiencesByUserID(ctx, portfolioOwnerID)
	if err != nil {
		eh.logger.Error("Failed to get experiences: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := experienceDto.FromExperiencesEntityToResponse(experiences, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if response == nil {
		eh.logger.Error("No experiences found for user ID %d", portfolioOwnerID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Experiences", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetExperience
//
//	@Summary		Get a specific experience
//	@Description	Retrieve a specific experience by ID
//	@Tags			Experiences
//	@Produce		json
//	@Param			id	path		int	true	"Experience ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.ExperienceResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/experiences/{id} [get]
func (eh *experienceHandler) GetExperience(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	experienceIDStr := r.PathValue("id")
	experienceID, err := strconv.Atoi(experienceIDStr)
	if err != nil || experienceID <= 0 {
		eh.logger.Error("Invalid experience ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid experience ID", "id", &err))
		return
	}

	experienceEntity, err := eh.experienceUseCase.GetExperienceByID(ctx, experienceID)
	if err != nil {
		eh.logger.Error("Failed to get experience %d: %v", experienceID, err)
		utils.WriteErrorResponse(w, err)
		return
	}
	portfolioOwnerID, err := utils.GetPortfolioOwnerID(eh.settingUseCase, ctx, w)

	if err != nil {
		eh.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	if experienceEntity.UserID != portfolioOwnerID {
		eh.logger.Error("Unauthorized access to experience %d", experienceID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Experience", strconv.Itoa(experienceID)))
		return
	}

	response := experienceDto.FromExperienceEntityToResponse(experienceEntity, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}
