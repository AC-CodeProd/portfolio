package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type SkillRepository interface {
	Create(ctx context.Context, skill *entities.Skill) (*entities.Skill, error)
	Update(ctx context.Context, skillID int, skill *entities.Skill) (*entities.Skill, error)
	Patch(ctx context.Context, skillID int, skill *entities.Skill) (*entities.Skill, error)
	Delete(ctx context.Context, skillID int) error

	GetByUserID(ctx context.Context, userID int) ([]*entities.Skill, error)
	GetByID(ctx context.Context, skillID int) (*entities.Skill, error)
	ExistsByID(ctx context.Context, skillID int) (bool, error)
	ExistsByNameAndUserID(ctx context.Context, name string, userID int) (bool, error)
	GetAll(ctx context.Context) ([]*entities.Skill, error)
}
