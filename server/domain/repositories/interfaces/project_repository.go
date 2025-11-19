package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *entities.Project) (*entities.Project, error)
	Update(ctx context.Context, projectID int, project *entities.Project) (*entities.Project, error)
	Patch(ctx context.Context, projectID int, project *entities.Project) (*entities.Project, error)
	Delete(ctx context.Context, projectID int) error

	GetByID(ctx context.Context, projectID int) (*entities.Project, error)
	GetAll(ctx context.Context, userID int) ([]*entities.Project, error)
}
