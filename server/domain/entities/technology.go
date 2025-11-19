package entities

import (
	"strings"
	"time"
)

type Technology struct {
	TechnologyID int
	UserID       int
	Name         string
	IconURL      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t *Technology) HasRequiredFields() bool {
	return t.Name != "" && t.IconURL != "" && t.UserID > 0
}

func (t *Technology) BelongsToUser(userID int) bool {
	return t.UserID == userID
}

func (t *Technology) MarkAsUpdated() {
	t.UpdatedAt = time.Now()
}

func (t *Technology) GetNormalizedName() string {
	return strings.ToLower(strings.TrimSpace(t.Name))
}

func (t *Technology) HasValidIconURL() bool {
	return strings.TrimSpace(t.IconURL) != ""
}

func (t *Technology) SetName(name string) {
	t.Name = strings.TrimSpace(name)
	t.MarkAsUpdated()
}

func (t *Technology) SetIconURL(iconURL string) {
	t.IconURL = strings.TrimSpace(iconURL)
	t.MarkAsUpdated()
}

func (t *Technology) IsWebIcon() bool {
	url := strings.ToLower(t.IconURL)
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func (t *Technology) GetIconFileName() string {
	parts := strings.Split(t.IconURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
