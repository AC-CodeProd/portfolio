package dto

import (
	"portfolio/domain"
	"portfolio/domain/entities"
	"strings"
	"time"
)

// @Description Request to create a new technology
type CreateTechnologyRequest struct {
	Name    string `json:"name" validate:"required"`
	IconURL string `json:"icon_url" validate:"required,url"`
} // @name CreateTechnologyRequest

// @Description Request to update an existing technology
type UpdateTechnologyRequest struct {
	ID      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	IconURL string `json:"icon_url" validate:"required,url"`
} // @name UpdateTechnologyRequest

// @Description Request to patch an existing technology
type PatchTechnologyRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,max=100"`
	IconURL *string `json:"icon_url,omitempty" validate:"omitempty,url"`
} // @name PatchTechnologyRequest

// @Description Request to delete a technology
type DeleteTechnologyRequest struct {
	ID int `json:"id" validate:"required"`
} // @name DeleteTechnologyRequest

type CreateBulkTechnologiesRequest struct {
	Technologies []CreateTechnologyRequest `json:"technologies" validate:"required"`
}

func (req *CreateTechnologyRequest) Validate() error {
	if strings.TrimSpace(req.Name) == "" {
		return domain.NewValidationError("technology name cannot be empty", "technology name", nil)
	}
	if len(req.Name) > 100 {
		return domain.NewValidationError("technology name cannot exceed 100 characters", "technology name", nil)
	}
	if strings.TrimSpace(req.IconURL) == "" {
		return domain.NewValidationError("technology icon_url cannot be empty", "technology icon_url", nil)
	}
	return nil
}

func (req *CreateTechnologyRequest) ToEntity(userID int) (*entities.Technology, error) {
	return &entities.Technology{
		UserID:    userID,
		Name:      req.Name,
		IconURL:   req.IconURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (req *CreateBulkTechnologiesRequest) Validate() error {
	if len(req.Technologies) == 0 {
		return domain.NewRequiredFieldError("at least one technology")
	}

	if len(req.Technologies) > 50 {
		return domain.NewValidationError("Cannot create more than 50 technologies at once", "technologies", nil)
	}

	nameMap := make(map[string]bool)
	for i, tech := range req.Technologies {
		if err := tech.Validate(); err != nil {
			return domain.NewValidationError("Technology "+string(rune(i+1))+": "+err.Error(), "technologies", &err)
		}

		normalizedName := strings.ToLower(strings.TrimSpace(tech.Name))
		if nameMap[normalizedName] {
			return domain.NewAlreadyExistsError("technology", tech.Name)
		}
		nameMap[normalizedName] = true
	}

	return nil
}

func (req *CreateBulkTechnologiesRequest) ToEntities(userID int) ([]*entities.Technology, error) {
	technologyEntities := make([]*entities.Technology, len(req.Technologies))
	now := time.Now()

	for i, techReq := range req.Technologies {
		technologyEntities[i] = &entities.Technology{
			Name:      strings.TrimSpace(techReq.Name),
			IconURL:   strings.TrimSpace(techReq.IconURL),
			UserID:    userID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	return technologyEntities, nil
}

func (req *UpdateTechnologyRequest) ToEntity(userID int) (*entities.Technology, error) {
	return &entities.Technology{
		TechnologyID: req.ID,
		UserID:       userID,
		Name:         req.Name,
		IconURL:      req.IconURL,
		UpdatedAt:    time.Now(),
	}, nil
}

func (req *PatchTechnologyRequest) Validate() error {
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return domain.NewValidationError("technology name cannot be empty", "technology name", nil)
		}
		if len(*req.Name) > 100 {
			return domain.NewValidationError("technology name cannot exceed 100 characters", "technology name", nil)
		}
	}

	return nil
}

func (req *PatchTechnologyRequest) ToEntity(id, userID int) (*entities.Technology, error) {
	technology := &entities.Technology{
		TechnologyID: id,
		UserID:       userID,
		UpdatedAt:    time.Now(),
	}

	if req.Name != nil {
		technology.Name = strings.TrimSpace(*req.Name)
	}
	if req.IconURL != nil {
		technology.IconURL = strings.TrimSpace(*req.IconURL)
	}

	return technology, nil
}
