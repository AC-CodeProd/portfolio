package handler

import (
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
	settingDto "portfolio/dto/setting"
	"portfolio/logger"
	"portfolio/shared"
	"time"
)

type settingHandler struct {
	settingUseCase *usecases.SettingUseCase
	logger         *logger.Logger
}

func NewSettingHandler(settingUseCase *usecases.SettingUseCase, logger *logger.Logger) []*routes.NamedRoute {
	settingHandler := settingHandler{
		settingUseCase: settingUseCase,
		logger:         logger,
	}

	return []*routes.NamedRoute{
		{
			Name:    "GetPublicSettingsHandler",
			Pattern: "GET /settings",
			Handler: settingHandler.GetSettings,
		},
	}
}

// GetSettings
//
//	@Summary		Get public application settings
//	@Description	Retrieve public application settings for portfolio display
//	@Tags			Settings
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.SettingResponse}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/v1/settings [get]
func (sh *settingHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	// Get request ID from context (set by middleware)
	ctx := r.Context()
	requestID := utils.GetRequestIDFromContext(ctx)

	settings, err := sh.settingUseCase.GetSettings(ctx)
	if err != nil {
		sh.logger.Error("Failed to retrieve public settings: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to retrieve settings", err))
		return
	}

	response := settingDto.FromSettingEntityToResponse(settings,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": requestID,
		},
	)
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}
