package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
	"time"
)

type projectRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewProjectRepository(db *sql.DB, logger *logger.Logger) interfaces.ProjectRepository {
	return &projectRepository{db: db, logger: logger}
}

func (repo *projectRepository) GetAll(ctx context.Context, userID int) ([]*entities.Project, error) {
	query := `SELECT project_id, user_id, project_title, project_description, project_short_description, 
	          project_technologies, project_github_url, project_image_url, project_status, 
	          project_created_at, project_updated_at 
	          FROM projects WHERE user_id = ? ORDER BY project_created_at DESC`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.logger.Error("Failed to getall projects: %v", err)
		return nil, domain.NewDatabaseError("project retrieval", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var projects []*entities.Project
	for rows.Next() {
		project := &entities.Project{}
		err := rows.Scan(
			&project.ProjectID,
			&project.UserID,
			&project.Title,
			&project.Description,
			&project.ShortDescription,
			&project.Technologies,
			&project.GithubURL,
			&project.ImageURL,
			&project.Status,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning project: %v", err)
			return nil, domain.NewDatabaseError("project scanning", err)
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		repo.logger.Error("Failed to iterating projects: %v", err)
		return nil, domain.NewDatabaseError("project iteration", err)
	}

	return projects, nil
}

func (repo *projectRepository) GetByID(ctx context.Context, projectID int) (*entities.Project, error) {
	query := `SELECT project_id, user_id, project_title, project_description, project_short_description, 
	          project_technologies, project_github_url, project_image_url, project_status, 
	          project_created_at, project_updated_at 
	          FROM projects WHERE project_id = ?`

	project := &entities.Project{}
	row := repo.db.QueryRowContext(ctx, query, projectID)

	err := row.Scan(
		&project.ProjectID,
		&project.UserID,
		&project.Title,
		&project.Description,
		&project.ShortDescription,
		&project.Technologies,
		&project.GithubURL,
		&project.ImageURL,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewNotFoundError("Project", fmt.Sprint(projectID))
		}
		repo.logger.Error("Failed to getbyid project: %v", err)
		return nil, domain.NewDatabaseError("project retrieval by ID", err)
	}

	return project, nil
}

func (repo *projectRepository) Create(ctx context.Context, project *entities.Project) (*entities.Project, error) {
	query := `INSERT INTO projects (user_id, project_title, project_description, project_short_description, 
	          project_technologies, project_github_url, project_image_url, project_status, 
	          project_created_at, project_updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := repo.db.ExecContext(ctx, query,
		project.UserID,
		project.Title,
		project.Description,
		project.ShortDescription,
		project.Technologies,
		project.GithubURL,
		project.ImageURL,
		project.Status,
		now,
		now,
	)

	if err != nil {
		repo.logger.Error("Failed to create project: %v", err)
		return nil, domain.NewDatabaseError("project creation", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to create project lastinsertid: %v", err)
		return nil, domain.NewDatabaseError("project id retrieval", err)
	}

	project.ProjectID = int(id)
	project.CreatedAt = now
	project.UpdatedAt = now

	return project, nil
}

func (repo *projectRepository) Update(ctx context.Context, projectID int, project *entities.Project) (*entities.Project, error) {
	query := `UPDATE projects SET project_title = ?, project_description = ?, project_short_description = ?, 
	          project_technologies = ?, project_github_url = ?, project_image_url = ?, project_status = ?, 
	          project_updated_at = ? WHERE project_id = ?`

	now := time.Now()
	result, err := repo.db.ExecContext(ctx, query,
		project.Title,
		project.Description,
		project.ShortDescription,
		project.Technologies,
		project.GithubURL,
		project.ImageURL,
		project.Status,
		now,
		projectID,
	)

	if err != nil {
		repo.logger.Error("Failed to update project: %v", err)
		return nil, domain.NewDatabaseError("project update", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Failed to update project rowsaffected: %v", err)
		return nil, domain.NewDatabaseError("project update verification", err)
	}

	if rowsAffected == 0 {
		return nil, domain.NewNotFoundError("Project", fmt.Sprint(project.ProjectID))
	}

	project.UpdatedAt = now
	return repo.GetByID(ctx, projectID)
}

func (repo *projectRepository) Patch(ctx context.Context, projectID int, project *entities.Project) (*entities.Project, error) {
	type field struct {
		Name  string
		Value interface{}
		Set   bool
	}

	var fields []field
	if project.Title != "" {
		fields = append(fields, field{"project_title", project.Title, true})
	}
	if project.Description != "" {
		fields = append(fields, field{"project_description", project.Description, true})
	}
	if project.ShortDescription != "" {
		fields = append(fields, field{"project_short_description", project.ShortDescription, true})
	}
	if project.Technologies != "" {
		fields = append(fields, field{"project_technologies", project.Technologies, true})
	}
	if project.GithubURL != "" {
		fields = append(fields, field{"project_github_url", project.GithubURL, true})
	}
	if project.ImageURL != "" {
		fields = append(fields, field{"project_image_url", project.ImageURL, true})
	}
	if project.Status != "" {
		fields = append(fields, field{"project_status", project.Status, true})
	}

	if len(fields) == 0 {
		return project, nil
	}

	query := "UPDATE projects SET "
	var args []interface{}
	for i, f := range fields {
		query += f.Name + " = ?"
		args = append(args, f.Value)
		if i < len(fields)-1 {
			query += ", "
		}
	}
	query += " WHERE project_id = ?"
	args = append(args, projectID)

	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.logger.Error("Failed to patch project: %v", err)
		return nil, fmt.Errorf("unable to patch project: %w", err)
	}
	return repo.GetByID(ctx, projectID)
}

func (repo *projectRepository) Delete(ctx context.Context, projectID int) error {
	query := `DELETE FROM projects WHERE project_id = ?`

	result, err := repo.db.ExecContext(ctx, query, projectID)
	if err != nil {
		repo.logger.Error("Failed to delete project: %v", err)
		return domain.NewDatabaseError("project deletion", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Failed to delete project rowsaffected: %v", err)
		return domain.NewDatabaseError("project deletion verification", err)
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("Project", fmt.Sprint(projectID))
	}

	return nil
}
