package dto

import (
	"portfolio/domain/entities"
	"portfolio/domain/validation"
	"strings"
	"time"
)

type CreateProjectRequest struct {
	Title            string `json:"title" validate:"required,max=200"`
	Description      string `json:"description,omitempty" validate:"omitempty,max=2000"`
	ShortDescription string `json:"short_description,omitempty" validate:"omitempty,max=500"`
	Technologies     string `json:"technologies,omitempty" validate:"omitempty,max=500"`
	GithubURL        string `json:"github_url,omitempty" validate:"omitempty,url"`
	ImageURL         string `json:"image_url,omitempty" validate:"omitempty,url"`
	Status           string `json:"status" validate:"required,oneof=active inactive archived"`
} // @name CreateProjectRequest

func (r *CreateProjectRequest) Validate() error {
	validator := validation.NewValidator()

	validator.Required("title", r.Title).MaxLength("title", r.Title, 200)

	if r.Description != "" {
		validator.MaxLength("description", r.Description, 2000)
	}

	if r.ShortDescription != "" {
		validator.MaxLength("short_description", r.ShortDescription, 500)
	}

	if r.Technologies != "" {
		validator.MaxLength("technologies", r.Technologies, 500)
	}

	if r.GithubURL != "" && strings.TrimSpace(r.GithubURL) != "" {
		validator.URL("github_url", r.GithubURL)
	}

	if r.ImageURL != "" && strings.TrimSpace(r.ImageURL) != "" {
		validator.URL("image_url", r.ImageURL)
	}

	validator.Required("status", r.Status)
	validStatuses := []string{"active", "inactive", "archived"}
	isValidStatus := false
	for _, status := range validStatuses {
		if r.Status == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		validator.Custom("status", false, "Status must be one of: active, inactive, archived")
	}

	if validator.HasErrors() {
		return validator.FirstError()
	}

	r.Sanitize()
	return nil
}

func (r *CreateProjectRequest) Sanitize() {
	r.Title = strings.TrimSpace(r.Title)
	r.Status = strings.TrimSpace(strings.ToLower(r.Status))

	if r.Description != "" {
		trimmed := strings.TrimSpace(r.Description)
		if trimmed == "" {
			r.Description = ""
		} else {
			r.Description = trimmed
		}
	}

	if r.ShortDescription != "" {
		trimmed := strings.TrimSpace(r.ShortDescription)
		if trimmed == "" {
			r.ShortDescription = ""
		} else {
			r.ShortDescription = trimmed
		}
	}

	if r.Technologies != "" {
		trimmed := strings.TrimSpace(r.Technologies)
		if trimmed == "" {
			r.Technologies = ""
		} else {
			r.Technologies = trimmed
		}
	}

	if r.GithubURL != "" {
		trimmed := strings.TrimSpace(r.GithubURL)
		if trimmed == "" {
			r.GithubURL = ""
		} else {
			r.GithubURL = trimmed
		}
	}

	if r.ImageURL != "" {
		trimmed := strings.TrimSpace(r.ImageURL)
		if trimmed == "" {
			r.ImageURL = ""
		} else {
			r.ImageURL = trimmed
		}
	}
}

type UpdateProjectRequest struct {
	Title            string `json:"title" validate:"required,max=200"`
	Description      string `json:"description,omitempty" validate:"omitempty,max=2000"`
	ShortDescription string `json:"short_description,omitempty" validate:"omitempty,max=500"`
	Technologies     string `json:"technologies,omitempty" validate:"omitempty,max=500"`
	GithubURL        string `json:"github_url,omitempty" validate:"omitempty,url"`
	ImageURL         string `json:"image_url,omitempty" validate:"omitempty,url"`
	Status           string `json:"status" validate:"required,oneof=active inactive archived"`
}

type PatchProjectRequest struct {
	Title            string `json:"title,omitempty" validate:"omitempty,max=200"`
	Description      string `json:"description,omitempty" validate:"omitempty,max=2000"`
	ShortDescription string `json:"short_description,omitempty" validate:"omitempty,max=500"`
	Technologies     string `json:"technologies,omitempty" validate:"omitempty,max=500"`
	GithubURL        string `json:"github_url,omitempty" validate:"omitempty,url"`
	ImageURL         string `json:"image_url,omitempty" validate:"omitempty,url"`
	Status           string `json:"status,omitempty" validate:"omitempty,oneof=active inactive archived"`
}

func (r *UpdateProjectRequest) Validate() error {
	validator := validation.NewValidator()

	validator.Required("title", r.Title).MaxLength("title", r.Title, 200)
	validator.Required("status", r.Status)

	if r.Description != "" {
		validator.MaxLength("description", r.Description, 2000)
	}

	if r.ShortDescription != "" {
		validator.MaxLength("short_description", r.ShortDescription, 500)
	}

	if r.Technologies != "" {
		validator.MaxLength("technologies", r.Technologies, 500)
	}

	if r.GithubURL != "" && strings.TrimSpace(r.GithubURL) != "" {
		validator.URL("github_url", r.GithubURL)
	}

	if r.ImageURL != "" && strings.TrimSpace(r.ImageURL) != "" {
		validator.URL("image_url", r.ImageURL)
	}

	validStatuses := []string{"active", "inactive", "archived"}
	isValidStatus := false
	for _, status := range validStatuses {
		if r.Status == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		validator.Custom("status", false, "Status must be one of: active, inactive, archived")
	}

	if validator.HasErrors() {
		return validator.FirstError()
	}

	r.Sanitize()
	return nil
}

func (r *PatchProjectRequest) Validate() error {
	validator := validation.NewValidator()

	if r.Title != "" {
		validator.Required("title", r.Title).MaxLength("title", r.Title, 200)
	}

	if r.Description != "" {
		validator.MaxLength("description", r.Description, 2000)
	}

	if r.ShortDescription != "" {
		validator.MaxLength("short_description", r.ShortDescription, 500)
	}

	if r.Technologies != "" {
		validator.MaxLength("technologies", r.Technologies, 500)
	}

	if r.GithubURL != "" && strings.TrimSpace(r.GithubURL) != "" {
		validator.URL("github_url", r.GithubURL)
	}

	if r.ImageURL != "" && strings.TrimSpace(r.ImageURL) != "" {
		validator.URL("image_url", r.ImageURL)
	}

	if r.Status != "" {
		validStatuses := []string{"active", "inactive", "archived"}
		isValidStatus := false
		for _, status := range validStatuses {
			if r.Status == status {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			validator.Custom("status", false, "Status must be one of: active, inactive, archived")
		}
	}

	if validator.HasErrors() {
		return validator.FirstError()
	}

	r.Sanitize()
	return nil
}

func (r *UpdateProjectRequest) Sanitize() {
	r.Title = strings.TrimSpace(r.Title)
	r.Status = strings.TrimSpace(strings.ToLower(r.Status))

	if r.Description != "" {
		trimmed := strings.TrimSpace(r.Description)
		if trimmed == "" {
			r.Description = ""
		} else {
			r.Description = trimmed
		}
	}

	if r.ShortDescription != "" {
		trimmed := strings.TrimSpace(r.ShortDescription)
		if trimmed == "" {
			r.ShortDescription = ""
		} else {
			r.ShortDescription = trimmed
		}
	}

	if r.Technologies != "" {
		trimmed := strings.TrimSpace(r.Technologies)
		if trimmed == "" {
			r.Technologies = ""
		} else {
			r.Technologies = trimmed
		}
	}

	if r.GithubURL != "" {
		trimmed := strings.TrimSpace(r.GithubURL)
		if trimmed == "" {
			r.GithubURL = ""
		} else {
			r.GithubURL = trimmed
		}
	}

	if r.ImageURL != "" {
		trimmed := strings.TrimSpace(r.ImageURL)
		if trimmed == "" {
			r.ImageURL = ""
		} else {
			r.ImageURL = trimmed
		}
	}
}

func (r *PatchProjectRequest) Sanitize() {
	if r.Title != "" {
		r.Title = strings.TrimSpace(r.Title)
	}
	if r.Description != "" {
		r.Description = strings.TrimSpace(r.Description)
	}
	if r.ShortDescription != "" {
		r.ShortDescription = strings.TrimSpace(r.ShortDescription)
	}
	if r.Technologies != "" {
		r.Technologies = strings.TrimSpace(r.Technologies)
	}
	if r.GithubURL != "" {
		r.GithubURL = strings.TrimSpace(r.GithubURL)
	}
	if r.ImageURL != "" {
		r.ImageURL = strings.TrimSpace(r.ImageURL)
	}
	if r.Status != "" {
		r.Status = strings.TrimSpace(r.Status)
	}
}

func (r *CreateProjectRequest) ToEntity(userID int) (*entities.Project, error) {
	now := time.Now()
	return &entities.Project{
		UserID:           userID,
		Title:            strings.TrimSpace(r.Title),
		Description:      r.Description,
		ShortDescription: r.ShortDescription,
		Technologies:     r.Technologies,
		GithubURL:        r.GithubURL,
		ImageURL:         r.ImageURL,
		Status:           r.Status,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (r *UpdateProjectRequest) ToEntity(id, userID int) (*entities.Project, error) {
	project := &entities.Project{
		ProjectID: id,
		UserID:    userID,
		Title:     strings.TrimSpace(r.Title),
		Status:    r.Status,
		UpdatedAt: time.Now(),
	}

	if r.Description != "" {
		project.Description = r.Description
	}
	if r.ShortDescription != "" {
		project.ShortDescription = r.ShortDescription
	}
	if r.Technologies != "" {
		project.Technologies = r.Technologies
	}
	if r.GithubURL != "" {
		project.GithubURL = r.GithubURL
	}
	if r.ImageURL != "" {
		project.ImageURL = r.ImageURL
	}

	return project, nil
}

func (r *PatchProjectRequest) ToEntity(id, userID int) (*entities.Project, error) {
	project := &entities.Project{
		ProjectID: id,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}

	if r.Title != "" {
		project.Title = strings.TrimSpace(r.Title)
	}
	if r.Description != "" {
		project.Description = r.Description
	}
	if r.ShortDescription != "" {
		project.ShortDescription = r.ShortDescription
	}
	if r.Technologies != "" {
		project.Technologies = r.Technologies
	}
	if r.GithubURL != "" {
		project.GithubURL = r.GithubURL
	}
	if r.ImageURL != "" {
		project.ImageURL = r.ImageURL
	}
	if r.Status != "" {
		project.Status = r.Status
	}

	return project, nil
}
