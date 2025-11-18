package sqlite

import (
	"context"
	"database/sql"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
)

type revokedTokenRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewRevokedTokenRepository(db *sql.DB, logger *logger.Logger) interfaces.RevokedTokenRepository {
	return &revokedTokenRepository{db: db, logger: logger}
}

func (repo *revokedTokenRepository) RevokedToken(ctx context.Context, userID int, token string) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO revoked_tokens (user_id, token) VALUES (?, ?)", userID, token)
	if err != nil {
		repo.logger.Error("failed to revoke token for user %d: %v", userID, err)
	}
	return err
}

func (repo *revokedTokenRepository) IsTokenRevoked(ctx context.Context, userID int, token string) (bool, error) {
	var count int
	err := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM revoked_tokens WHERE user_id = ? AND token = ?", userID, token).Scan(&count)
	if err != nil {
		repo.logger.Error("failed to check if token is revoked for user %d: %v", userID, err)
		return false, err
	}
	return count > 0, nil
}
