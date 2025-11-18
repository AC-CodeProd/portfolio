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

type experienceRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewExperienceRepository(db *sql.DB, logger *logger.Logger) interfaces.ExperienceRepository {
	return &experienceRepository{db: db, logger: logger}
}

func (repo *experienceRepository) Create(ctx context.Context, experience *entities.Experience) (*entities.Experience, error) {
	query := `INSERT INTO experiences (user_id, experience_company_name, experience_job_title, 
			  experience_start_date, experience_end_date, experience_description) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	result, err := repo.db.ExecContext(ctx, query,
		experience.UserID,
		experience.CompanyName,
		experience.JobTitle,
		experience.StartDate.String(),
		experience.EndDate.String(),
		experience.Description,
	)
	if err != nil {
		repo.logger.Error("Failed to create experience: %v", err)
		return nil, domain.NewDatabaseError("create experience", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to create experience lastinsertid: %v", err)
		return nil, domain.NewDatabaseError("get experience ID", err)
	}

	experience.ExperienceID = int(id)
	experience.CreatedAt = time.Now()
	experience.UpdatedAt = time.Now()

	return experience, nil
}

func (repo *experienceRepository) GetByID(ctx context.Context, experienceID int) (*entities.Experience, error) {
	var experience entities.Experience
	query := `SELECT experience_id, user_id, experience_company_name, experience_job_title, 
			  experience_start_date, experience_end_date, experience_description, 
			  experience_created_at, experience_updated_at 
			  FROM experiences WHERE experience_id = ?`

	row := repo.db.QueryRowContext(ctx, query, experienceID)
	err := row.Scan(
		&experience.ExperienceID,
		&experience.UserID,
		&experience.CompanyName,
		&experience.JobTitle,
		&experience.StartDate,
		&experience.EndDate,
		&experience.Description,
		&experience.CreatedAt,
		&experience.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to get experience: %v", err)
		return nil, domain.NewDatabaseError("retrieve experience", err)
	}

	return &experience, nil
}

func (repo *experienceRepository) GetByUserID(ctx context.Context, userID int) ([]*entities.Experience, error) {
	query := `SELECT experience_id, user_id, experience_company_name, experience_job_title, 
			  experience_start_date, experience_end_date, experience_description,
			  experience_created_at, experience_updated_at 
			  FROM experiences WHERE user_id = ? ORDER BY experience_start_date DESC`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.logger.Error("Failed to getbyuserid experiences: %v", err)
		return nil, domain.NewDatabaseError("retrieve experiences by user ID", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var experiences []*entities.Experience
	for rows.Next() {
		var experience entities.Experience
		err := rows.Scan(
			&experience.ExperienceID,
			&experience.UserID,
			&experience.CompanyName,
			&experience.JobTitle,
			&experience.StartDate,
			&experience.EndDate,
			&experience.Description,
			&experience.CreatedAt,
			&experience.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning experience: %v", err)
			continue
		}

		experiences = append(experiences, &experience)
	}

	return experiences, nil
}

func (repo *experienceRepository) Update(ctx context.Context, experienceID int, experience *entities.Experience) (*entities.Experience, error) {
	query := `UPDATE experiences SET experience_company_name = ?, experience_job_title = ?, 
			  experience_start_date = ?, experience_end_date = ?, experience_description = ?, 
			  experience_updated_at = ? WHERE experience_id = ?`

	_, err := repo.db.ExecContext(ctx, query,
		experience.CompanyName,
		experience.JobTitle,
		experience.StartDate.String(),
		experience.EndDate.String(),
		experience.Description,
		time.Now(),
		experienceID,
	)
	if err != nil {
		repo.logger.Error("Failed to update experience: %v", err)
		return nil, domain.NewDatabaseError("update experience", err)
	}

	return repo.GetByID(ctx, experienceID)
}

func (repo *experienceRepository) Patch(ctx context.Context, experienceID int, experience *entities.Experience) (*entities.Experience, error) {
	type field struct {
		Name  string
		Value interface{}
		Set   bool
	}

	var fields []field
	if experience.CompanyName != "" {
		fields = append(fields, field{"experience_company_name", experience.CompanyName, true})
	}
	if experience.JobTitle != "" {
		fields = append(fields, field{"experience_job_title", experience.JobTitle, true})
	}
	if !experience.StartDate.IsZero() {
		fields = append(fields, field{"experience_start_date", experience.StartDate.String(), true})
	}
	if !experience.EndDate.IsZero() {
		fields = append(fields, field{"experience_end_date", experience.EndDate.String(), true})
	}
	if experience.Description != "" {
		fields = append(fields, field{"experience_description", experience.Description, true})
	}

	if len(fields) == 0 {
		return experience, nil
	}

	query := "UPDATE experiences SET "
	var args []interface{}
	for i, f := range fields {
		query += f.Name + " = ?"
		args = append(args, f.Value)
		if i < len(fields)-1 {
			query += ", "
		}
	}
	query += " WHERE experience_id = ?"
	args = append(args, experienceID)

	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.logger.Error("Failed to patch experience: %v", err)
		return nil, fmt.Errorf("unable to patch experience: %w", err)
	}
	return repo.GetByID(ctx, experienceID)
}

func (repo *experienceRepository) Delete(ctx context.Context, experienceID int) error {
	query := `DELETE FROM experiences WHERE experience_id = ?`

	_, err := repo.db.ExecContext(ctx, query, experienceID)
	if err != nil {
		repo.logger.Error("Failed to delete experience: %v", err)
		return domain.NewDatabaseError("delete experience", err)
	}

	return nil
}

func (repo *experienceRepository) ExistsByID(ctx context.Context, experienceID int) (bool, error) {
	query := `SELECT COUNT(*) FROM experiences WHERE experience_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, experienceID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyid experience: %v", err)
		return false, domain.NewDatabaseError("check experience existence", err)
	}

	return count > 0, nil
}

func (repo *experienceRepository) ExistsByJobTitleCompanyAndUserID(ctx context.Context, jobTitle, companyName string, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM experiences WHERE experience_job_title = ? AND experience_company_name = ? AND user_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, jobTitle, companyName, userID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyjobtitlecompanyanduserid experience: %v", err)
		return false, domain.NewDatabaseError("check experience existence by job title and company", err)
	}

	return count > 0, nil
}

func (repo *experienceRepository) GetCurrentExperiences(ctx context.Context, userID int) ([]*entities.Experience, error) {
	query := `SELECT experience_id, user_id, experience_company_name, experience_job_title, 
			  experience_start_date, experience_end_date, experience_description,
			  experience_created_at, experience_updated_at 
			  FROM experiences WHERE user_id = ? AND (experience_end_date IS NULL OR experience_end_date = '' OR experience_end_date = '0000-00-00')
			  ORDER BY experience_start_date DESC`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.logger.Error("Failed to getcurrentexperiences: %v", err)
		return nil, domain.NewDatabaseError("retrieve current experiences", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var experiences []*entities.Experience
	for rows.Next() {
		var experience entities.Experience
		err := rows.Scan(
			&experience.ExperienceID,
			&experience.UserID,
			&experience.CompanyName,
			&experience.JobTitle,
			&experience.StartDate,
			&experience.EndDate,
			&experience.Description,
			&experience.CreatedAt,
			&experience.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning current experience: %v", err)
			continue
		}

		experiences = append(experiences, &experience)
	}

	return experiences, nil
}

func (repo *experienceRepository) GetAll(ctx context.Context) ([]*entities.Experience, error) {
	query := `SELECT experience_id, user_id, experience_company_name, experience_job_title, 
			  experience_start_date, experience_end_date, experience_description,
			  experience_created_at, experience_updated_at 
			  FROM experiences ORDER BY experience_start_date DESC`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		repo.logger.Error("Failed to getall experiences: %v", err)
		return nil, domain.NewDatabaseError("retrieve all experiences", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var experiences []*entities.Experience
	for rows.Next() {
		var experience entities.Experience
		err := rows.Scan(
			&experience.ExperienceID,
			&experience.UserID,
			&experience.CompanyName,
			&experience.JobTitle,
			&experience.StartDate,
			&experience.EndDate,
			&experience.Description,
			&experience.CreatedAt,
			&experience.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning experience: %v", err)
			continue
		}

		experiences = append(experiences, &experience)
	}

	return experiences, nil
}
