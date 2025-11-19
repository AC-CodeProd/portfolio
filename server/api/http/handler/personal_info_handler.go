package handler

import (
	"fmt"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
	personalInfoDto "portfolio/dto/personal_info"
	"portfolio/logger"
	"portfolio/shared"
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
			Name:    "GetPersonalInfoHandler",
			Pattern: "GET /personal-info",
			Handler: personalInfoHandler.GetPersonalInfo,
		},
	}
}

// GetPersonalInfo
//
//	@Summary		Get personal information
//	@Description	Retrieve personal information for the portfolio
//	@Tags			Personal Info
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.PersonalInfoResponse}
//	@Failure		404	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/personal-info [get]
func (pih *personalInfoHandler) GetPersonalInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolioOwnerID, err := utils.GetPortfolioOwnerID(pih.settingUseCase, ctx, w)

	if err != nil {
		pih.logger.Error("Failed to get portfolio owner ID: %v", err)
		return
	}
	personalInfo, err := pih.personalInfoUseCase.GetPersonalInfoByUserID(ctx, portfolioOwnerID)
	if err != nil {
		pih.logger.Error("Failed to get personal info: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	if personalInfo == nil {
		pih.logger.Error("Personal info not found for user ID %d", portfolioOwnerID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Personal Info", "current user"))
		return
	}

	user, err := pih.personalInfoUseCase.GetUserByID(ctx, personalInfo.UserID)
	if err != nil {
		pih.logger.Error("Failed to get user by ID %d: %v", personalInfo.UserID, err)
		utils.WriteErrorResponse(w, err)
		return
	}

	if user == nil {
		pih.logger.Error("User not found for ID %d", personalInfo.UserID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("User", fmt.Sprint(personalInfo.UserID)))
		return
	}

	personalInfoResponse := personalInfoDto.NewPersonalInfoResponse(user, personalInfo, &shared.Meta{
		"timestamp":  time.Now().Format(time.RFC3339),
		"request_id": utils.GetRequestIDFromContext(ctx),
	})

	if personalInfoResponse == nil {
		pih.logger.Error("Personal info response is nil for user ID %d", personalInfo.UserID)
		utils.WriteErrorResponse(w, domain.NewNotFoundError("Personal Info", "current user"))
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, personalInfoResponse)
}
