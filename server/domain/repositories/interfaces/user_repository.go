package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, entity *entities.User) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	UpdateLastLogin(ctx context.Context, userID int) error
	ExistsByID(ctx context.Context, userID int) (bool, error)
}
