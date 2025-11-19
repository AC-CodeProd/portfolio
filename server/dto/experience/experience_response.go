package dto

import (
	"portfolio/domain/entities"
	"portfolio/shared"
	"time"
)

// @Description Experience represents a work experience entry in the portfolio
type Experience struct {
	ExperienceID      int       `json:"experience_id"`
	UserID            int       `json:"user_id"`
	JobTitle          string    `json:"job_title"`
	CompanyName       string    `json:"company_name"`
	StartDate         string    `json:"start_date"`
	EndDate           string    `json:"end_date,omitempty"`
	Description       string    `json:"description"`
	IsCurrentPosition bool      `json:"is_current_position"`
	DurationInMonths  int       `json:"duration_in_months"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
} //@name Experience

// @Description Response for a list of experiences
type ExperienceListResponse struct {
	Experiences []*Experience `json:"experiences"`
	Meta        *shared.Meta  `json:"meta"`
} //@name ExperienceListResponse

// @Description Response for an experience
type ExperienceResponse struct {
	Experience *Experience  `json:"experience"`
	Meta       *shared.Meta `json:"meta"`
} //@name ExperienceResponse

func FromExperienceEntityToResponse(experience *entities.Experience, meta *shared.Meta) *ExperienceResponse {
	if experience == nil {
		return nil
	}

	_experience := &Experience{
		ExperienceID:      experience.ExperienceID,
		UserID:            experience.UserID,
		JobTitle:          experience.JobTitle,
		CompanyName:       experience.CompanyName,
		StartDate:         experience.StartDate.Format("2006-01-02"),
		Description:       experience.Description,
		IsCurrentPosition: experience.IsCurrentPosition(),
		DurationInMonths:  experience.GetDurationInMonths(),
		CreatedAt:         experience.CreatedAt,
		UpdatedAt:         experience.UpdatedAt,
	}

	if !experience.IsCurrentPosition() {
		_experience.EndDate = experience.EndDate.Format("2006-01-02")
	}

	return &ExperienceResponse{
		Experience: _experience,
		Meta:       meta,
	}
}

func FromExperiencesEntityToResponse(experiences []*entities.Experience, meta *shared.Meta) *ExperienceListResponse {
	if experiences == nil {
		return nil
	}

	var experienceResponses []*Experience
	for _, experience := range experiences {
		_experience := &Experience{
			ExperienceID:      experience.ExperienceID,
			UserID:            experience.UserID,
			JobTitle:          experience.JobTitle,
			CompanyName:       experience.CompanyName,
			StartDate:         experience.StartDate.Format("2006-01-02"),
			Description:       experience.Description,
			IsCurrentPosition: experience.IsCurrentPosition(),
			DurationInMonths:  experience.GetDurationInMonths(),
			CreatedAt:         experience.CreatedAt,
			UpdatedAt:         experience.UpdatedAt,
		}

		if !experience.IsCurrentPosition() {
			_experience.EndDate = experience.EndDate.Format("2006-01-02")
		}

		experienceResponses = append(experienceResponses, _experience)
	}

	return &ExperienceListResponse{
		Experiences: experienceResponses,
		Meta:        meta,
	}
}
