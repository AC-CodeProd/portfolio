package dto

import (
	"portfolio/domain/entities"
	"portfolio/shared"
	"time"
)

// @Description Setting represents the settings of the portfolio
type Setting struct {
	ShowProjects     bool   `json:"show_projects"`
	PortfolioOwnerID int    `json:"portfolio_owner_id"`
	SiteName         string `json:"site_name"`
	SiteDescription  string `json:"site_description"`
	ContactEmail     string `json:"contact_email"`
	Theme            string `json:"theme"`
	Language         string `json:"language"`
	MaintenanceMode  bool   `json:"maintenance_mode"`
	UpdatedAt        string `json:"updated_at"`
} // @name Setting

// @Description Response for a setting
type SettingResponse struct {
	Setting *Setting     `json:"setting"`
	Meta    *shared.Meta `json:"meta"`
} //@name SettingResponse

func FromSettingEntityToResponse(setting *entities.SettingJson, meta *shared.Meta) *SettingResponse {
	if setting == nil {
		return nil
	}

	return &SettingResponse{
		Setting: &Setting{
			ShowProjects:     setting.ShowProjects,
			PortfolioOwnerID: setting.PortfolioOwnerID,
			SiteName:         setting.SiteName,
			SiteDescription:  setting.SiteDescription,
			ContactEmail:     setting.ContactEmail,
			Theme:            setting.Theme,
			Language:         setting.Language,
			MaintenanceMode:  setting.MaintenanceMode,
			UpdatedAt:        setting.UpdatedAt.Format(time.RFC3339),
		},
		Meta: meta,
	}
}
