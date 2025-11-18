package dto

import (
	"portfolio/domain"
	"portfolio/domain/entities"
	"strings"
	"time"
)

// @Description Request to create a new skill
type CreateSkillRequest struct {
	Name  string `json:"name" validate:"required"`
	Level int    `json:"level" validate:"required"`
} // @name CreateSkillRequest

// @Description Request to update an existing skill
type UpdateSkillRequest struct {
	Name  string `json:"name" validate:"required"`
	Level int    `json:"level" validate:"required"`
} // @name UpdateSkillRequest

// @Description Request to patch an existing skill
type PatchSkillRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,max=100"`
	Level *int    `json:"level,omitempty" validate:"omitempty,min=1,max=5"`
} // @name PatchSkillRequest

// @Description Request to create multiple skills in bulk
type CreateBulkSkillsRequest struct {
	Skills []CreateSkillRequest `json:"skills" validate:"required"`
} // @name CreateBulkSkillsRequest

func (req *CreateSkillRequest) Validate() error {
	if strings.TrimSpace(req.Name) == "" {
		return domain.NewRequiredFieldError("name")
	}

	if len(req.Name) > 100 {
		return domain.NewValidationError("Skill name cannot exceed 100 characters", "name", nil)
	}

	if req.Level < 1 || req.Level > 5 {
		return domain.NewValidationError("Skill level must be between 1 and 5", "level", nil)
	}

	return nil
}

func (req *CreateSkillRequest) ToEntity(userID int) (*entities.Skill, error) {
	now := time.Now()
	return &entities.Skill{
		Name:      strings.TrimSpace(req.Name),
		Level:     req.Level,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (req *CreateBulkSkillsRequest) Validate() error {
	if len(req.Skills) == 0 {
		return domain.NewValidationError("At least one skill is required", "skills", nil)
	}

	if len(req.Skills) > 50 {
		return domain.NewValidationError("Cannot create more than 50 skills at once", "skills", nil)
	}

	nameMap := make(map[string]bool)
	for i, skill := range req.Skills {
		if err := skill.Validate(); err != nil {
			return domain.NewValidationError("Skill "+string(rune(i+1))+": "+err.Error(), "skills", &err)
		}

		normalizedName := strings.ToLower(strings.TrimSpace(skill.Name))
		if nameMap[normalizedName] {
			return domain.NewAlreadyExistsError("skill", skill.Name)
		}
		nameMap[normalizedName] = true
	}

	return nil
}

func (req *CreateBulkSkillsRequest) ToEntities(userID int) ([]*entities.Skill, error) {
	skillEntities := make([]*entities.Skill, len(req.Skills))
	now := time.Now()

	for i, skillReq := range req.Skills {
		skillEntities[i] = &entities.Skill{
			Name:      strings.TrimSpace(skillReq.Name),
			Level:     skillReq.Level,
			UserID:    userID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	return skillEntities, nil
}

func (req *UpdateSkillRequest) Validate() error {
	if strings.TrimSpace(req.Name) == "" {
		return domain.NewRequiredFieldError("name")
	}

	if len(req.Name) > 100 {
		return domain.NewValidationError("Skill name cannot exceed 100 characters", "name", nil)
	}

	if req.Level < 1 || req.Level > 5 {
		return domain.NewValidationError("Skill level must be between 1 and 5", "level", nil)
	}

	return nil
}

func (req *UpdateSkillRequest) ToEntity(id, userID int) (*entities.Skill, error) {
	return &entities.Skill{
		SkillID:   id,
		Name:      strings.TrimSpace(req.Name),
		Level:     req.Level,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}, nil
}

func (req *PatchSkillRequest) Validate() error {
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return domain.NewValidationError("Skill name cannot be empty", "name", nil)
		}
		if len(*req.Name) > 100 {
			return domain.NewValidationError("Skill name cannot exceed 100 characters", "name", nil)
		}
	}

	if req.Level != nil {
		if *req.Level < 1 || *req.Level > 5 {
			return domain.NewValidationError("Skill level must be between 1 and 5", "level", nil)
		}
	}

	return nil
}

func (req *PatchSkillRequest) ToEntity(id, userID int) (*entities.Skill, error) {
	skill := &entities.Skill{
		SkillID:   id,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}

	if req.Name != nil {
		skill.Name = strings.TrimSpace(*req.Name)
	}
	if req.Level != nil {
		skill.Level = *req.Level
	}

	return skill, nil
}
