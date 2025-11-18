package entities

import "time"

type Project struct {
	ProjectID        int
	UserID           int
	Title            string
	Description      string
	ShortDescription string
	Technologies     string
	GithubURL        string
	ImageURL         string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (p *Project) IsActive() bool {
	return p.Status == "active"
}

func (p *Project) IsArchived() bool {
	return p.Status == "archived"
}

func (p *Project) IsInactive() bool {
	return p.Status == "inactive"
}

func (p *Project) SetStatus(status string) bool {
	validStatuses := []string{"active", "inactive", "archived"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			p.Status = status
			p.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func (p *Project) HasValidStatus() bool {
	validStatuses := []string{"active", "inactive", "archived"}
	for _, validStatus := range validStatuses {
		if p.Status == validStatus {
			return true
		}
	}
	return false
}

func (p *Project) BelongsToUser(userID int) bool {
	return p.UserID == userID
}

func (p *Project) HasRequiredFields() bool {
	return p.Title != "" && p.Status != "" && p.UserID > 0
}

func (p *Project) MarkAsUpdated() {
	p.UpdatedAt = time.Now()
}
