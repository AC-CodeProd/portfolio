package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type TechnologyRepository interface {
	Create(ctx context.Context, technology *entities.Technology) (*entities.Technology, error)
	Update(ctx context.Context, technologyID int, technology *entities.Technology) (*entities.Technology, error)
	Patch(ctx context.Context, technologyID int, technology *entities.Technology) (*entities.Technology, error)
	Delete(ctx context.Context, technologyID int) error

	GetByID(ctx context.Context, technologyID int) (*entities.Technology, error)
	GetByUserID(ctx context.Context, userID int) ([]*entities.Technology, error)
	ExistsByID(ctx context.Context, technologyID int) (bool, error)
	ExistsByNameAndUserID(ctx context.Context, name string, userID int) (bool, error)
	GetAll(ctx context.Context) ([]*entities.Technology, error)
	GetByNames(ctx context.Context, names []string, userID int) ([]*entities.Technology, error)
}
