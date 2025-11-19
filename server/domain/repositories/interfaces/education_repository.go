package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type EducationRepository interface {
	Create(ctx context.Context, education *entities.Education) (*entities.Education, error)
	Update(ctx context.Context, educationID int, education *entities.Education) (*entities.Education, error)
	Patch(ctx context.Context, educationID int, education *entities.Education) (*entities.Education, error)
	Delete(ctx context.Context, educationID int) error

	GetByUserID(ctx context.Context, userID int) ([]*entities.Education, error)
	GetByID(ctx context.Context, educationID int) (*entities.Education, error)
	ExistsByID(ctx context.Context, educationID int) (bool, error)
	ExistsByDegreeInstitutionAndUserID(ctx context.Context, degree, institution string, userID int) (bool, error)
	GetAll(ctx context.Context) ([]*entities.Education, error)
	GetCurrentEducations(ctx context.Context, userID int) ([]*entities.Education, error)
}
