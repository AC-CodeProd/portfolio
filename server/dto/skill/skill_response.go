package dto

import (
	"portfolio/domain/entities"
	"portfolio/shared"
	"time"
)

// @Description Skill represents a skill entry in the portfolio
type Skill struct {
	SkillID   int       `json:"skill_id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name Skill

// @Description Response for a list of skills
type SkillListResponse struct {
	Skills []*Skill     `json:"skills"`
	Meta   *shared.Meta `json:"meta"`
} //@name SkillListResponse

// @Description Response for a skill
type SkillResponse struct {
	Skill *Skill       `json:"skill"`
	Meta  *shared.Meta `json:"meta"`
} //@name SkillResponse

// @Description Response for bulk skill operations
type SkillBulkResponse struct {
	Skills []*Skill           `json:"skills"`
	Errors []*shared.APIError `json:"errors,omitempty"`
	Meta   *shared.Meta       `json:"meta"`
} // @name SkillBulkResponse

func FromSkillsEntityToResponse(skills []*entities.Skill, meta *shared.Meta) *SkillListResponse {
	if skills == nil {
		return nil
	}

	var skillResponses []*Skill

	for _, skill := range skills {
		_skill := &Skill{
			SkillID:   skill.SkillID,
			UserID:    skill.UserID,
			Name:      skill.Name,
			Level:     skill.Level,
			CreatedAt: skill.CreatedAt,
			UpdatedAt: skill.UpdatedAt,
		}
		skillResponses = append(skillResponses, _skill)
	}

	return &SkillListResponse{
		Skills: skillResponses,
		Meta:   meta,
	}
}

func FromSkillEntityToResponse(skill *entities.Skill, meta *shared.Meta) *SkillResponse {
	if skill == nil {
		return nil
	}

	return &SkillResponse{
		Skill: &Skill{
			SkillID:   skill.SkillID,
			UserID:    skill.UserID,
			Name:      skill.Name,
			Level:     skill.Level,
			CreatedAt: skill.CreatedAt,
			UpdatedAt: skill.UpdatedAt,
		},
		Meta: meta,
	}
}

func FromSkillsEntityForBulkToResponse(skills []*entities.Skill, meta *shared.Meta) *SkillBulkResponse {
	if skills == nil {
		return nil
	}

	var skillResponses []*Skill

	for _, skill := range skills {
		_skill := &Skill{
			SkillID:   skill.SkillID,
			UserID:    skill.UserID,
			Name:      skill.Name,
			Level:     skill.Level,
			CreatedAt: skill.CreatedAt,
			UpdatedAt: skill.UpdatedAt,
		}
		skillResponses = append(skillResponses, _skill)
	}

	return &SkillBulkResponse{
		Skills: skillResponses,
		Meta:   meta,
	}
}
