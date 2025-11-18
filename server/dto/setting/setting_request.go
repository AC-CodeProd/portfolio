package dto

import (
	"portfolio/domain/entities"
	"strings"
)

// @Description UpdateSettingRequest represents the payload to update settings
type UpdateSettingRequest struct {
	ShowProjects     bool   `json:"show_projects"`
	PortfolioOwnerID int    `json:"portfolio_owner_id" validate:"required,min=1"`
	SiteName         string `json:"site_name" validate:"required,min=1,max=100"`
	SiteDescription  string `json:"site_description" validate:"required,min=10,max=500"`
	ContactEmail     string `json:"contact_email" validate:"required,email"`
	Theme            string `json:"theme" validate:"required,oneof=light dark auto"`
	Language         string `json:"language" validate:"required,oneof=fr en es"`
	MaintenanceMode  bool   `json:"maintenance_mode"`
} // @name UpdateSettingRequest

type SocialLinks struct {
	GitHub   string `json:"github" validate:"omitempty,url"`
	LinkedIn string `json:"linkedin" validate:"omitempty,url"`
	Twitter  string `json:"twitter" validate:"omitempty,url"`
	Website  string `json:"website" validate:"omitempty,url"`
}

func (r *UpdateSettingRequest) Validate() error {
	r.Sanitize()
	return nil
}

func (r *UpdateSettingRequest) Sanitize() {
	r.SiteName = strings.TrimSpace(r.SiteName)
	r.SiteDescription = strings.TrimSpace(r.SiteDescription)
	r.ContactEmail = strings.TrimSpace(strings.ToLower(r.ContactEmail))
}

func (r *UpdateSettingRequest) ToEntity() *entities.SettingJson {
	return &entities.SettingJson{
		ShowProjects:     r.ShowProjects,
		PortfolioOwnerID: r.PortfolioOwnerID,
		SiteName:         r.SiteName,
		SiteDescription:  r.SiteDescription,
		ContactEmail:     r.ContactEmail,
		Theme:            r.Theme,
		Language:         r.Language,
		MaintenanceMode:  r.MaintenanceMode,
	}
}
