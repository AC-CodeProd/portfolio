package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
	"strings"
	"time"
)

type technologyRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewTechnologyRepository(db *sql.DB, logger *logger.Logger) interfaces.TechnologyRepository {
	return &technologyRepository{db: db, logger: logger}
}

func (repo *technologyRepository) Create(ctx context.Context, technology *entities.Technology) (*entities.Technology, error) {
	query := `INSERT INTO technologies (user_id, technology_name, technology_icon_url) 
			  VALUES (?, ?, ?)`

	result, err := repo.db.ExecContext(ctx, query,
		technology.UserID,
		technology.Name,
		technology.IconURL,
	)
	if err != nil {
		repo.logger.Error("Failed to create technology: %v", err)
		return nil, domain.NewDatabaseError("create technology", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to create technology lastinsertid: %v", err)
		return nil, domain.NewDatabaseError("get technology ID", err)
	}

	technology.TechnologyID = int(id)
	technology.CreatedAt = time.Now()
	technology.UpdatedAt = time.Now()

	return technology, nil
}

func (repo *technologyRepository) GetByID(ctx context.Context, technologyID int) (*entities.Technology, error) {
	var technology entities.Technology
	query := `SELECT technology_id, user_id, technology_name, technology_icon_url, 
			  technology_created_at, technology_updated_at 
			  FROM technologies WHERE technology_id = ?`

	row := repo.db.QueryRowContext(ctx, query, technologyID)
	err := row.Scan(
		&technology.TechnologyID,
		&technology.UserID,
		&technology.Name,
		&technology.IconURL,
		&technology.CreatedAt,
		&technology.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to get technology: %v", err)
		return nil, domain.NewDatabaseError("retrieve technology", err)
	}

	return &technology, nil
}

func (repo *technologyRepository) GetByUserID(ctx context.Context, userID int) ([]*entities.Technology, error) {
	query := `SELECT technology_id, user_id, technology_name, technology_icon_url,
			  technology_created_at, technology_updated_at 
			  FROM technologies WHERE user_id = ? ORDER BY technology_name`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.logger.Error("Failed to getbyuserid technologies: %v", err)
		return nil, domain.NewDatabaseError("retrieve technologies by user ID", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var technologies []*entities.Technology
	for rows.Next() {
		var technology entities.Technology
		err := rows.Scan(
			&technology.TechnologyID,
			&technology.UserID,
			&technology.Name,
			&technology.IconURL,
			&technology.CreatedAt,
			&technology.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning technology: %v", err)
			continue
		}

		technologies = append(technologies, &technology)
	}

	return technologies, nil
}

func (repo *technologyRepository) Update(ctx context.Context, technologyID int, technology *entities.Technology) (*entities.Technology, error) {
	query := `UPDATE technologies SET technology_name = ?, technology_icon_url = ?, 
			  technology_updated_at = ? WHERE technology_id = ?`

	_, err := repo.db.ExecContext(ctx, query,
		technology.Name,
		technology.IconURL,
		time.Now(),
		technologyID,
	)
	if err != nil {
		repo.logger.Error("Failed to update technology: %v", err)
		return nil, domain.NewDatabaseError("update technology", err)
	}

	return repo.GetByID(ctx, technologyID)
}

func (repo *technologyRepository) Patch(ctx context.Context, technologyID int, technology *entities.Technology) (*entities.Technology, error) {
	type field struct {
		Name  string
		Value interface{}
		Set   bool
	}

	var fields []field
	if technology.Name != "" {
		fields = append(fields, field{"technology_name", technology.Name, true})
	}
	if technology.IconURL != "" {
		fields = append(fields, field{"technology_icon_url", technology.IconURL, true})
	}

	if len(fields) == 0 {
		return technology, nil
	}

	query := "UPDATE technologies SET "
	var args []interface{}
	for i, f := range fields {
		query += f.Name + " = ?"
		args = append(args, f.Value)
		if i < len(fields)-1 {
			query += ", "
		}
	}
	query += " WHERE technology_id = ?"
	args = append(args, technologyID)

	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.logger.Error("Failed to patch technology: %v", err)
		return nil, fmt.Errorf("unable to patch technology: %w", err)
	}
	return repo.GetByID(ctx, technologyID)
}

func (repo *technologyRepository) Delete(ctx context.Context, technologyID int) error {
	query := `DELETE FROM technologies WHERE technology_id = ?`

	_, err := repo.db.ExecContext(ctx, query, technologyID)
	if err != nil {
		repo.logger.Error("Failed to delete technology: %v", err)
		return domain.NewDatabaseError("delete technology", err)
	}

	return nil
}

func (repo *technologyRepository) ExistsByID(ctx context.Context, technologyID int) (bool, error) {
	query := `SELECT COUNT(*) FROM technologies WHERE technology_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, technologyID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyid technology: %v", err)
		return false, domain.NewDatabaseError("check technology existence", err)
	}

	return count > 0, nil
}

func (repo *technologyRepository) ExistsByNameAndUserID(ctx context.Context, name string, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM technologies WHERE technology_name = ? AND user_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, name, userID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbynameanduserid technology: %v", err)
		return false, domain.NewDatabaseError("check technology existence by name", err)
	}

	return count > 0, nil
}

func (repo *technologyRepository) GetByNames(ctx context.Context, names []string, userID int) ([]*entities.Technology, error) {
	if len(names) == 0 {
		return []*entities.Technology{}, nil
	}

	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names)+1)
	args[0] = userID

	for i, name := range names {
		placeholders[i] = "?"
		args[i+1] = name
	}

	query := `SELECT technology_id, user_id, technology_name, technology_icon_url,
			  technology_created_at, technology_updated_at 
			  FROM technologies WHERE user_id = ? AND technology_name IN (` +
		strings.Join(placeholders, ",") + `) ORDER BY technology_name`

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		repo.logger.Error("Failed to getbynames technologies: %v", err)
		return nil, domain.NewDatabaseError("retrieve technologies by names", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var technologies []*entities.Technology
	for rows.Next() {
		var technology entities.Technology
		err := rows.Scan(
			&technology.TechnologyID,
			&technology.UserID,
			&technology.Name,
			&technology.IconURL,
			&technology.CreatedAt,
			&technology.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning technology: %v", err)
			continue
		}

		technologies = append(technologies, &technology)
	}

	return technologies, nil
}

func (repo *technologyRepository) GetAll(ctx context.Context) ([]*entities.Technology, error) {
	query := `SELECT technology_id, user_id, technology_name, technology_icon_url,
			  technology_created_at, technology_updated_at 
			  FROM technologies ORDER BY technology_name`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		repo.logger.Error("Failed to getall technologies: %v", err)
		return nil, domain.NewDatabaseError("retrieve all technologies", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var technologies []*entities.Technology
	for rows.Next() {
		var technology entities.Technology
		err := rows.Scan(
			&technology.TechnologyID,
			&technology.UserID,
			&technology.Name,
			&technology.IconURL,
			&technology.CreatedAt,
			&technology.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning technology: %v", err)
			continue
		}

		technologies = append(technologies, &technology)
	}

	return technologies, nil
}
