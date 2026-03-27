package dto

import (
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"strings"
	"time"
)

type CreateExperienceRequest struct {
	JobTitle    string `json:"job_title" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
} // @name CreateExperienceRequest

// @Description Request to create multiple experiences in bulk
type CreateBulkExperiencesRequest struct {
	Experiences []CreateExperienceRequest `json:"experiences" validate:"required"`
} // @name CreateBulkExperiencesRequest

func (req *CreateBulkExperiencesRequest) Validate() error {
	if len(req.Experiences) == 0 {
		return domain.NewValidationError("At least one experience is required", "experiences", nil)
	}

	if len(req.Experiences) > 50 {
		return domain.NewValidationError("Cannot create more than 50 experiences at once", "experiences", nil)
	}

	nameMap := make(map[string]bool)
	for i, experience := range req.Experiences {
		if err := experience.Validate(); err != nil {
			return domain.NewValidationError("Experience "+string(rune(i+1))+": "+err.Error(), "experiences", &err)
		}

		normalizedName := strings.ToLower(strings.TrimSpace(experience.JobTitle))
		if nameMap[normalizedName] {
			return domain.NewAlreadyExistsError("experience", experience.JobTitle)
		}
		nameMap[normalizedName] = true
	}

	return nil
}

type UpdateExperienceRequest struct {
	JobTitle    string `json:"job_title" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
} // @name UpdateExperienceRequest

type PatchExperienceRequest struct {
	JobTitle    string `json:"job_title,omitempty" validate:"omitempty"`
	CompanyName string `json:"company_name,omitempty" validate:"omitempty"`
	StartDate   string `json:"start_date,omitempty" validate:"omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	Description string `json:"description,omitempty"`
} // @name PatchExperienceRequest

type DeleteExperienceRequest struct {
	ID int `json:"id" validate:"required"`
} // @name DeleteExperienceRequest

func (req *CreateExperienceRequest) Validate() error {
	if strings.TrimSpace(req.JobTitle) == "" {
		return domain.NewRequiredFieldError("job_title")
	}

	if len(req.JobTitle) > 100 {
		return domain.NewValidationError("Job title cannot exceed 100 characters", "job_title", nil)
	}

	if strings.TrimSpace(req.CompanyName) == "" {
		return domain.NewRequiredFieldError("company_name")
	}

	if len(req.CompanyName) > 100 {
		return domain.NewValidationError("Company name cannot exceed 100 characters", "company_name", nil)
	}

	if strings.TrimSpace(req.StartDate) == "" {
		return domain.NewRequiredFieldError("start_date")
	}

	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		return domain.NewInvalidFormatError("start_date", "YYYY-MM-DD")
	}

	if req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
			return domain.NewInvalidFormatError("end_date", "YYYY-MM-DD")
		}

		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		if endDate.Before(startDate) {
			return domain.NewValidationError("End date must be after start date", "end_date", nil)
		}
	}

	if len(req.Description) > 1000 {
		return domain.NewValidationError("Description cannot exceed 1000 characters", "description", nil)
	}

	return nil
}

func (req *UpdateExperienceRequest) Validate() error {
	if strings.TrimSpace(req.JobTitle) == "" {
		return domain.NewRequiredFieldError("job_title")
	}

	if len(req.JobTitle) > 100 {
		return domain.NewValidationError("Job title cannot exceed 100 characters", "job_title", nil)
	}

	if strings.TrimSpace(req.CompanyName) == "" {
		return domain.NewRequiredFieldError("company_name")
	}

	if len(req.CompanyName) > 100 {
		return domain.NewValidationError("Company name cannot exceed 100 characters", "company_name", nil)
	}

	if strings.TrimSpace(req.StartDate) == "" {
		return domain.NewRequiredFieldError("start_date")
	}

	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		return domain.NewInvalidFormatError("start_date", "YYYY-MM-DD")
	}

	if req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
			return domain.NewInvalidFormatError("end_date", "YYYY-MM-DD")
		}

		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		if endDate.Before(startDate) {
			return domain.NewValidationError("End date must be after start date", "end_date", nil)
		}
	}

	if len(req.Description) > 1000 {
		return domain.NewValidationError("Description cannot exceed 1000 characters", "description", nil)
	}

	return nil
}

func (req *CreateExperienceRequest) ToEntity(userID int) (*entities.Experience, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	var endDate time.Time
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		endDate = parsedEndDate
	}

	return &entities.Experience{
		UserID:      userID,
		JobTitle:    strings.TrimSpace(req.JobTitle),
		CompanyName: strings.TrimSpace(req.CompanyName),
		StartDate:   startDate,
		EndDate:     endDate,
		Description: strings.TrimSpace(req.Description),
	}, nil
}

func (req *CreateBulkExperiencesRequest) ToEntities(userID int) ([]*entities.Experience, error) {
	experienceEntities := make([]*entities.Experience, 0, len(req.Experiences))
	for _, expReq := range req.Experiences {
		entity, err := expReq.ToEntity(userID)
		if err != nil {
			return nil, err
		}
		experienceEntities = append(experienceEntities, entity)
	}
	return experienceEntities, nil
}

func (req *UpdateExperienceRequest) ToEntity(experienceID, userID int) (*entities.Experience, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	var endDate time.Time
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		endDate = parsedEndDate
	}

	return &entities.Experience{
		ExperienceID: experienceID,
		UserID:       userID,
		JobTitle:     strings.TrimSpace(req.JobTitle),
		CompanyName:  strings.TrimSpace(req.CompanyName),
		StartDate:    startDate,
		EndDate:      endDate,
		Description:  strings.TrimSpace(req.Description),
	}, nil
}

func (req *PatchExperienceRequest) Validate() error {
	if req.JobTitle != "" {
		if strings.TrimSpace(req.JobTitle) == "" {
			return domain.NewValidationError("Job title cannot be empty", "job_title", nil)
		}
		if len(req.JobTitle) > 100 {
			return domain.NewValidationError("Job title cannot exceed 100 characters", "job_title", nil)
		}
	}

	if req.CompanyName != "" {
		if strings.TrimSpace(req.CompanyName) == "" {
			return domain.NewValidationError("Company name cannot be empty", "company_name", nil)
		}
		if len(req.CompanyName) > 100 {
			return domain.NewValidationError("Company name cannot exceed 100 characters", "company_name", nil)
		}
	}

	if req.StartDate != "" {
		if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
			return domain.NewInvalidFormatError("start_date", "YYYY-MM-DD")
		}
	}

	if req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
			return domain.NewInvalidFormatError("end_date", "YYYY-MM-DD")
		}
	}

	return nil
}

func (req *PatchExperienceRequest) ToEntity(experienceID, userID int) (*entities.Experience, error) {
	experience := &entities.Experience{
		ExperienceID: experienceID,
		UserID:       userID,
	}

	if req.JobTitle != "" {
		experience.JobTitle = strings.TrimSpace(req.JobTitle)
	}

	if req.CompanyName != "" {
		experience.CompanyName = strings.TrimSpace(req.CompanyName)
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %v", err)
		}
		experience.StartDate = startDate
	}

	if req.EndDate != "" {
		if strings.TrimSpace(req.EndDate) != "" {
			endDate, err := time.Parse("2006-01-02", req.EndDate)
			if err != nil {
				return nil, fmt.Errorf("invalid end date format: %v", err)
			}
			experience.EndDate = endDate
		}
	}

	if req.Description != "" {
		experience.Description = strings.TrimSpace(req.Description)
	}

	return experience, nil
}
