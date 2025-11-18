package sqlite

import (
	"context"
	"database/sql"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
	"time"
)

type userRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewUserRepository(db *sql.DB, logger *logger.Logger) interfaces.UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (repo *userRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := `
        INSERT INTO users (user_username, user_email, user_password, user_role, user_is_active, user_created_at, user_updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	result, err := repo.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		repo.logger.Error("Failed to createuser: %v", err)
		return nil, domain.NewDatabaseError("user creation", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to createuser lastinsertid: %v", err)
		return nil, domain.NewDatabaseError("user id retrieval", err)
	}
	user.ID = int(id)

	return user, nil
}

func (repo *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `
        SELECT user_id, user_username, user_email, user_password, user_role, user_is_active, user_created_at, user_updated_at, user_last_login
        FROM users WHERE user_username = ?
    `
	user := &entities.User{}
	var lastLogin sql.NullTime

	err := repo.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to getbyusername: %v", err)
		return nil, domain.NewDatabaseError("user retrieval by username", err)
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time
	}

	return user, nil
}

func (repo *userRepository) UpdateLastLogin(ctx context.Context, userID int) error {
	query := "UPDATE users SET user_last_login = ? WHERE user_id = ?"
	_, err := repo.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		repo.logger.Error("Failed to updatelastlogin: %v", err)
		return domain.NewDatabaseError("last login update", err)
	}
	return nil
}

func (repo *userRepository) ExistsByID(ctx context.Context, userID int) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE user_id = ?"
	var count int
	err := repo.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyid: %v", err)
		return false, domain.NewDatabaseError("user existence check", err)
	}
	return count > 0, nil
}
