package dto

import (
	"portfolio/domain"
	"portfolio/domain/entities"
	"strings"
	"time"
)

type CreateEducationRequest struct {
	Degree      string     `json:"degree" validate:"required"`
	Institution string     `json:"institution" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description *string    `json:"description,omitempty"`
} // @name CreateEducationRequest

type UpdateEducationRequest struct {
	ID          int        `json:"id" validate:"required"`
	Degree      string     `json:"degree" validate:"required"`
	Institution string     `json:"institution" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description *string    `json:"description,omitempty"`
} // @name UpdateEducationRequest

type PatchEducationRequest struct {
	Degree      *string    `json:"degree,omitempty" validate:"omitempty"`
	Institution *string    `json:"institution,omitempty" validate:"omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty" validate:"omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description *string    `json:"description,omitempty"`
} // @name PatchEducationRequest

type DeleteEducationRequest struct {
	ID int `json:"id" validate:"required"`
} // @name DeleteEducationRequest

func (req *CreateEducationRequest) Validate() error {
	if strings.TrimSpace(req.Degree) == "" {
		return domain.NewRequiredFieldError("degree")
	}

	if strings.TrimSpace(req.Institution) == "" {
		return domain.NewRequiredFieldError("institution")
	}

	if req.StartDate.IsZero() {
		return domain.NewRequiredFieldError("start_date")
	}

	return nil
}

func (req *UpdateEducationRequest) Validate() error {
	if req.ID <= 0 {
		return domain.NewRequiredFieldError("id")
	}

	if strings.TrimSpace(req.Degree) == "" {
		return domain.NewRequiredFieldError("degree")
	}

	if strings.TrimSpace(req.Institution) == "" {
		return domain.NewRequiredFieldError("institution")
	}

	if req.StartDate.IsZero() {
		return domain.NewRequiredFieldError("start_date")
	}

	return nil
}

func (req *CreateEducationRequest) ToEntity(userID int) (*entities.Education, error) {
	education := &entities.Education{
		UserID:      userID,
		Degree:      req.Degree,
		Institution: req.Institution,
		StartDate:   req.StartDate,
		Description: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.EndDate != nil {
		education.EndDate = req.EndDate
	}

	if req.Description != nil {
		education.Description = *req.Description
	}

	return education, nil
}

func (req *UpdateEducationRequest) ToEntity(userID int) (*entities.Education, error) {
	education := &entities.Education{
		EducationID: req.ID,
		UserID:      userID,
		Degree:      req.Degree,
		Institution: req.Institution,
		StartDate:   req.StartDate,
		Description: "",
		UpdatedAt:   time.Now(),
	}

	if req.EndDate != nil {
		education.EndDate = req.EndDate
	}

	if req.Description != nil {
		education.Description = *req.Description
	}

	return education, nil
}

func (req *PatchEducationRequest) Validate() error {
	return nil
}

func (req *PatchEducationRequest) ToEntity(id, userID int) (*entities.Education, error) {
	education := &entities.Education{
		EducationID: id,
		UserID:      userID,
		UpdatedAt:   time.Now(),
	}

	if req.Degree != nil {
		education.Degree = *req.Degree
	}
	if req.Institution != nil {
		education.Institution = *req.Institution
	}
	if req.StartDate != nil {
		education.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		education.EndDate = req.EndDate
	}
	if req.Description != nil {
		education.Description = *req.Description
	}

	return education, nil
}
