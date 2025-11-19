package entities

import (
	"time"
)

type Skill struct {
	SkillID   int
	UserID    int
	Name      string
	Level     int // 1-5
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Skill) HasRequiredFields() bool {
	return s.Name != "" && s.UserID > 0 && s.Level >= 1 && s.Level <= 5
}

func (s *Skill) BelongsToUser(userID int) bool {
	return s.UserID == userID
}

func (s *Skill) MarkAsUpdated() {
	s.UpdatedAt = time.Now()
}

func (s *Skill) IsValidLevel() bool {
	return s.Level >= 1 && s.Level <= 5
}

func (s *Skill) GetLevelDescription() string {
	switch s.Level {
	case 1:
		return "Beginner"
	case 2:
		return "Basic"
	case 3:
		return "Intermediate"
	case 4:
		return "Advanced"
	case 5:
		return "Expert"
	default:
		return "Unknown"
	}
}

func (s *Skill) SetLevel(level int) bool {
	if level >= 1 && level <= 5 {
		s.Level = level
		s.MarkAsUpdated()
		return true
	}
	return false
}
