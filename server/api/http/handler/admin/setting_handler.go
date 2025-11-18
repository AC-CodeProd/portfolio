package admin

import (
	"encoding/json"
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
			Name:    "GetSettingsHandler",
			Pattern: "GET /settings",
			Handler: settingHandler.GetSettings,
		},
		{
			Name:    "UpdateSettingsHandler",
			Pattern: "PUT /settings",
			Handler: settingHandler.UpdateSettings,
		},
		{
			Name:    "ResetSettingsHandler",
			Pattern: "POST /settings/reset",
			Handler: settingHandler.ResetSettings,
		},
	}
}

// GetSettings
//
//	@Summary		Get application settings
//	@Description	Retrieve current application settings
//	@Tags			Admin Settings
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.SettingResponse}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/settings [get]
//	@Security		BearerAuth
func (sh *settingHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	settings, err := sh.settingUseCase.GetSettings(ctx)
	if err != nil {
		sh.logger.Error("Failed to retrieve settings: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to retrieve settings", err))
		return
	}

	response := settingDto.FromSettingEntityToResponse(settings,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// UpdateSettings
//
//	@Summary		Update application settings
//	@Description	Update application settings
//	@Tags			Admin Settings
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UpdateSettingRequest	true	"Setting update request"
//	@Success		200		{object}	shared.APIResponse{data=dto.SettingResponse}
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/settings [put]
//	@Security		BearerAuth
func (sh *settingHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req settingDto.UpdateSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sh.logger.Error("Failed to decode request body: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("body", "Invalid request body", &err))
		return
	}

	if err := req.Validate(); err != nil {
		sh.logger.Error("Settings validation error: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("settings", err.Error(), nil))
		return
	}

	settingEntity := req.ToEntity()

	if err := sh.settingUseCase.Upsert(ctx, settingEntity); err != nil {
		sh.logger.Error("Failed to update settings: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to update settings", err))
		return
	}

	updatedSettings, err := sh.settingUseCase.GetSettings(ctx)
	if err != nil {
		sh.logger.Error("Failed to retrieve updated settings: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to retrieve updated settings", err))
		return
	}

	response := settingDto.FromSettingEntityToResponse(updatedSettings,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

// ResetSettings
//
//	@Summary		Reset settings to default
//	@Description	Reset application settings to default values
//	@Tags			Admin Settings
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=dto.SettingResponse}
//	@Failure		500	{object}	shared.APIResponse{errors=[]shared.APIError}
//	@Router			/admin/settings/reset [post]
//	@Security		BearerAuth
func (sh *settingHandler) ResetSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defaultSettings := &settingDto.UpdateSettingRequest{
		ShowProjects:     true,
		PortfolioOwnerID: 1,
		SiteName:         "My Portfolio",
		SiteDescription:  "Personal portfolio showcasing my projects and skills",
		ContactEmail:     "",
		Theme:            "light",
		Language:         "en",
		MaintenanceMode:  false,
	}

	settingEntity := defaultSettings.ToEntity()

	if err := sh.settingUseCase.Upsert(ctx, settingEntity); err != nil {
		sh.logger.Error("Failed to reset settings: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to reset settings", err))
		return
	}

	resetSettings, err := sh.settingUseCase.GetSettings(ctx)
	if err != nil {
		sh.logger.Error("Failed to retrieve reset settings: %v", err)
		utils.WriteErrorResponse(w, domain.NewInternalError("Failed to retrieve reset settings", err))
		return
	}

	response := settingDto.FromSettingEntityToResponse(resetSettings,
		&shared.Meta{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": utils.GetRequestIDFromContext(ctx),
		})
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}
