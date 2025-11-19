package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type PersonalInfoRepository interface {
	Get(ctx context.Context) (*entities.PersonalInfo, error)
	Create(ctx context.Context, personalInfo *entities.PersonalInfo) (*entities.PersonalInfo, error)
	Update(ctx context.Context, personalInfoId int, personalInfo *entities.PersonalInfo) (*entities.PersonalInfo, error)
	Patch(ctx context.Context, personalInfoId int, personalInfo *entities.PersonalInfo) (*entities.PersonalInfo, error)
	Delete(ctx context.Context, personalInfoId int) error

	GetByID(ctx context.Context, personalInfoId int) (*entities.PersonalInfo, error)
	ExistsByUserID(ctx context.Context, userID int) (bool, error)
	GetByUserID(ctx context.Context, userID int) (*entities.PersonalInfo, error)
	GetUserByID(ctx context.Context, userID int) (*entities.User, error)
}
