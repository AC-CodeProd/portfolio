package dto

import (
	"portfolio/domain/entities"
	"portfolio/shared"
	"time"
)

// @Description Technology represents a technology entry in the portfolio
type Technology struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	IconURL   string    `json:"icon_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name Technology

// @Description Response for a list of technologies
type TechnologyListResponse struct {
	Technologies []*Technology `json:"technologies"`
	Meta         *shared.Meta  `json:"meta"`
} //@name TechnologyListResponse

// @Description Response for a technology
type TechnologyResponse struct {
	Technology *Technology  `json:"technology"`
	Meta       *shared.Meta `json:"meta"`
} //@name TechnologyResponse

// @Description Response for bulk technology operations
type TechnologyBulkResponse struct {
	Technologies []*Technology      `json:"technologies"`
	Errors       []*shared.APIError `json:"errors,omitempty"`
	Meta         *shared.Meta       `json:"meta"`
} // @name TechnologyBulkResponse

func FromTechnologiesEntityToResponse(technologies []*entities.Technology, meta *shared.Meta) *TechnologyListResponse {
	if technologies == nil {
		return nil
	}

	var technologyResponses []*Technology

	for _, technology := range technologies {
		_technology := &Technology{
			ID:        technology.TechnologyID,
			UserID:    technology.UserID,
			Name:      technology.Name,
			IconURL:   technology.IconURL,
			CreatedAt: technology.CreatedAt,
			UpdatedAt: technology.UpdatedAt,
		}
		technologyResponses = append(technologyResponses, _technology)
	}

	return &TechnologyListResponse{
		Technologies: technologyResponses,
		Meta:         meta,
	}
}

func FromTechnologyEntityToResponse(technology *entities.Technology, meta *shared.Meta) *TechnologyResponse {
	if technology == nil {
		return nil
	}

	return &TechnologyResponse{
		Technology: &Technology{
			ID:        technology.TechnologyID,
			UserID:    technology.UserID,
			Name:      technology.Name,
			IconURL:   technology.IconURL,
			CreatedAt: technology.CreatedAt,
			UpdatedAt: technology.UpdatedAt,
		},
		Meta: meta,
	}
}

func FromTechnologiesEntityForBulkToResponse(technologies []*entities.Technology, meta *shared.Meta) *TechnologyBulkResponse {
	if technologies == nil {
		return nil
	}

	var technologyResponses []*Technology

	for _, technology := range technologies {
		_technology := &Technology{
			ID:        technology.TechnologyID,
			UserID:    technology.UserID,
			Name:      technology.Name,
			IconURL:   technology.IconURL,
			CreatedAt: technology.CreatedAt,
			UpdatedAt: technology.UpdatedAt,
		}
		technologyResponses = append(technologyResponses, _technology)
	}

	return &TechnologyBulkResponse{
		Technologies: technologyResponses,
		Meta:         meta,
	}
}
