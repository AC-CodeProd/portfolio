package dto

import (
	"portfolio/domain/entities"
	"portfolio/shared"
)

// @Description Education represents an education entry in the portfolio
type Education struct {
	ID          int     `json:"id"`
	Degree      string  `json:"degree"`
	Institution string  `json:"institution"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
} //@name Education

// @Description Response for a list of educations
type EducationListResponse struct {
	Educations []*Education `json:"educations"`
	Meta       *shared.Meta `json:"meta"`
} //@name EducationListResponse

// @Description Response for a education
type EducationResponse struct {
	Education *Education   `json:"education"`
	Meta      *shared.Meta `json:"meta"`
} //@name EducationResponse

func FromEducationEntityToResponse(education *entities.Education, meta *shared.Meta) *EducationResponse {
	if education == nil {
		return nil
	}

	_education := &Education{
		ID:          education.EducationID,
		Degree:      education.Degree,
		Institution: education.Institution,
		StartDate:   education.StartDate.Format("01/2006"),
		Description: education.Description,
		CreatedAt:   education.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   education.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if education.EndDate != nil && !education.EndDate.IsZero() {
		endDate := education.EndDate.Format("01/2006")
		_education.EndDate = &endDate
	}

	return &EducationResponse{
		Education: _education,
		Meta:      meta,
	}
}

func FromEducationsEntityToResponse(educations []*entities.Education, meta *shared.Meta) *EducationListResponse {
	if educations == nil {
		return nil
	}

	var educationResponses []*Education
	for _, education := range educations {
		_education := &Education{
			ID:          education.EducationID,
			Degree:      education.Degree,
			Institution: education.Institution,
			StartDate:   education.StartDate.Format("01/2006"),
			Description: education.Description,
			CreatedAt:   education.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   education.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if education.EndDate != nil && !education.EndDate.IsZero() {
			endDate := education.EndDate.Format("01/2006")
			_education.EndDate = &endDate
		}
		educationResponses = append(educationResponses, _education)
	}

	return &EducationListResponse{
		Educations: educationResponses,
		Meta:       meta,
	}
}
