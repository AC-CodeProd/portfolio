package entities

import "time"

type Setting struct {
	SettingJson      []byte
	SettingKey       string
	SettingCreatedAt string
	SettingUpdatedAt string
}

type SettingJson struct {
	ShowProjects     bool
	PortfolioOwnerID int
	SiteName         string
	SiteDescription  string
	ContactEmail     string
	Theme            string
	Language         string
	MaintenanceMode  bool
	UpdatedAt        time.Time
}
