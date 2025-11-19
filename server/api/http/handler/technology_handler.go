package handler

import (
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
	dto "portfolio/dto/technology"
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
		AbstractHandler: AbstractHandler{
			settingUseCase: settingUseCase,
		},
		technologyUseCase: technologyUseCase,
		logger:            logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetTechnologiesHandler",
			Pattern: "GET /technologies",
			Handler: technologyHandler.GetTechnologies,
		},
		{
			Name:    "GetTechnologyHandler",
			Pattern: "GET /technologies/{id}",
			Handler: technologyHandler.GetTechnology,
		},
	}
}

// GetTechnologies
//
//	@Summary		Get all technologies
//	@Description	Retrieve all technologies for the portfolio
//	@Tags			Technologies
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.TechnologyListResponse}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/technologies [get]
func (th *technologyHandler) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(th.settingUseCase, ctx, w)

	if err != nil {
		th.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	technologies, err := th.technologyUseCase.GetTechnologiesByUserID(ctx, portfolioOwnerID)
	if err != nil {
		th.logger.Error("Failed to get technologies: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	response := dto.FromTechnologiesEntityToResponse(technologies, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if response == nil {
		th.logger.Error("Technologies response is nil for user ID %d", portfolioOwnerID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Technologies", ""))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// GetTechnology
//
//	@Summary		Get a specific technology
//	@Description	Retrieve a specific technology by ID
//	@Tags			Technologies
//	@Produce		json
//	@Param			id	path		int	true	"Technology ID"
//	@Success		200	{object}	shared.APIResponse{data=dto.TechnologyResponse}
//	@Failure		400	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/technologies/{id} [get]
func (th *technologyHandler) GetTechnology(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	technologyIDStr := r.PathValue("id")
	technologyID, err := strconv.Atoi(technologyIDStr)
	if err != nil || technologyID <= 0 {
		th.logger.Error("Invalid technology ID format: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid technology ID", "id", &err))
		return
	}

	technology, err := th.technologyUseCase.GetTechnologyByID(ctx, technologyID)
	if err != nil {
		th.logger.Error("Failed to get technology by ID %d: %v", technologyID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	portfolioOwnerID, err := utils.GetPortfolioOwnerID(th.settingUseCase, ctx, w)

	if err != nil {
		th.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}

	if technology.UserID != portfolioOwnerID {
		th.logger.Error("Unauthorized access to technology %d", technologyID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Technology", strconv.Itoa(technologyID)))
		return
	}

	response := dto.FromTechnologyEntityToResponse(technology, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	utils.WriteSuccessResponse(w, http.StatusOK, response)
}
