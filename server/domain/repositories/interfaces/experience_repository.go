package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type ExperienceRepository interface {
	Create(ctx context.Context, experience *entities.Experience) (*entities.Experience, error)
	Update(ctx context.Context, experienceID int, experience *entities.Experience) (*entities.Experience, error)
	Patch(ctx context.Context, experienceID int, experience *entities.Experience) (*entities.Experience, error)
	Delete(ctx context.Context, experienceID int) error

	GetByID(ctx context.Context, experienceID int) (*entities.Experience, error)
	GetByUserID(ctx context.Context, userID int) ([]*entities.Experience, error)
	ExistsByID(ctx context.Context, experienceID int) (bool, error)
	ExistsByJobTitleCompanyAndUserID(ctx context.Context, jobTitle, companyName string, userID int) (bool, error)
	GetAll(ctx context.Context) ([]*entities.Experience, error)
	GetCurrentExperiences(ctx context.Context, userID int) ([]*entities.Experience, error)
}
