package interfaces

import "context"

type RevokedTokenRepository interface {
	RevokedToken(ctx context.Context, userID int, token string) error
	IsTokenRevoked(ctx context.Context, userID int, token string) (bool, error)
}
