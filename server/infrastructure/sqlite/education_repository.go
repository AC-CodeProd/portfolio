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

type educationRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewEducationRepository(db *sql.DB, logger *logger.Logger) interfaces.EducationRepository {
	return &educationRepository{db: db, logger: logger}
}

func (repo *educationRepository) Create(ctx context.Context, education *entities.Education) (*entities.Education, error) {
	query := `INSERT INTO educations (user_id, education_degree, education_institution, 
			  education_start_date, education_end_date, education_description) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	result, err := repo.db.ExecContext(ctx, query,
		education.UserID,
		education.Degree,
		education.Institution,
		education.StartDate.String(),
		education.EndDate.String(),
		education.Description,
	)
	if err != nil {
		repo.logger.Error("Failed to create education: %v", err)
		return nil, domain.NewDatabaseError("create education", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to get education ID after creation: %v", err)
		return nil, domain.NewDatabaseError("get education ID", err)
	}

	education.EducationID = int(id)
	education.CreatedAt = time.Now()
	education.UpdatedAt = time.Now()

	return education, nil
}

func (repo *educationRepository) GetByID(ctx context.Context, educationID int) (*entities.Education, error) {
	var education entities.Education
	query := `SELECT education_id, user_id, education_degree, education_institution, 
			  education_start_date, education_end_date, education_description, 
			  education_created_at, education_updated_at 
			  FROM educations WHERE education_id = ?`

	row := repo.db.QueryRowContext(ctx, query, educationID)
	err := row.Scan(
		&education.EducationID,
		&education.UserID,
		&education.Degree,
		&education.Institution,
		&education.StartDate,
		&education.EndDate,
		&education.Description,
		&education.CreatedAt,
		&education.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to retrieve education by ID: %v", err)
		return nil, domain.NewDatabaseError("retrieve education", err)
	}

	return &education, nil
}

func (repo *educationRepository) GetByUserID(ctx context.Context, userID int) ([]*entities.Education, error) {
	query := `SELECT education_id, user_id, education_degree, education_institution, 
			  education_start_date, education_end_date, education_description,
			  education_created_at, education_updated_at 
			  FROM educations WHERE user_id = ? ORDER BY education_start_date DESC`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.logger.Error("Failed to retrieve educations by user ID: %v", err)
		return nil, domain.NewDatabaseError("retrieve educations by user ID", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var educations []*entities.Education
	for rows.Next() {
		var education entities.Education
		err := rows.Scan(
			&education.EducationID,
			&education.UserID,
			&education.Degree,
			&education.Institution,
			&education.StartDate,
			&education.EndDate,
			&education.Description,
			&education.CreatedAt,
			&education.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning education: %v", err)
			continue
		}

		educations = append(educations, &education)
	}

	return educations, nil
}

func (repo *educationRepository) Update(ctx context.Context, educationID int, education *entities.Education) (*entities.Education, error) {
	query := `UPDATE educations SET education_degree = ?, education_institution = ?, 
			  education_start_date = ?, education_end_date = ?, education_description = ?, 
			  education_updated_at = ? WHERE education_id = ?`

	_, err := repo.db.ExecContext(ctx, query,
		education.Degree,
		education.Institution,
		education.StartDate.String(),
		education.EndDate.String(),
		education.Description,
		time.Now(),
		educationID,
	)
	if err != nil {
		repo.logger.Error("Failed to update education: %v", err)
		return nil, domain.NewDatabaseError("update education", err)
	}

	return repo.GetByID(ctx, educationID)
}

func (repo *educationRepository) Patch(ctx context.Context, educationID int, education *entities.Education) (*entities.Education, error) {
	type field struct {
		Name  string
		Value interface{}
		Set   bool
	}

	var fields []field
	if education.Degree != "" {
		fields = append(fields, field{"education_degree", education.Degree, true})
	}
	if education.Institution != "" {
		fields = append(fields, field{"education_institution", education.Institution, true})
	}
	if !education.StartDate.IsZero() {
		fields = append(fields, field{"education_start_date", education.StartDate.String(), true})
	}
	if !education.EndDate.IsZero() {
		fields = append(fields, field{"education_end_date", education.EndDate.String(), true})
	}
	if education.Description != "" {
		fields = append(fields, field{"education_description", education.Description, true})
	}

	if len(fields) == 0 {
		return education, nil
	}

	query := "UPDATE educations SET "
	var args []interface{}
	for i, f := range fields {
		query += f.Name + " = ?"
		args = append(args, f.Value)
		if i < len(fields)-1 {
			query += ", "
		}
	}
	query += " WHERE education_id = ?"
	args = append(args, educationID)

	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.logger.Error("Failed to patch education: %v", err)
		return nil, fmt.Errorf("unable to patch education: %w", err)
	}
	return repo.GetByID(ctx, educationID)
}

func (repo *educationRepository) Delete(ctx context.Context, educationID int) error {
	query := `DELETE FROM educations WHERE education_id = ?`

	_, err := repo.db.ExecContext(ctx, query, educationID)
	if err != nil {
		repo.logger.Error("Failed to delete education: %v", err)
		return domain.NewDatabaseError("delete education", err)
	}

	return nil
}

func (repo *educationRepository) ExistsByID(ctx context.Context, educationID int) (bool, error) {
	query := `SELECT COUNT(*) FROM educations WHERE education_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, educationID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyid education: %v", err)
		return false, domain.NewDatabaseError("check education existence", err)
	}

	return count > 0, nil
}

func (repo *educationRepository) ExistsByDegreeInstitutionAndUserID(ctx context.Context, degree, institution string, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM educations WHERE education_degree = ? AND education_institution = ? AND user_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, degree, institution, userID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbydegreeinstitutionanduserid education: %v", err)
		return false, domain.NewDatabaseError("check education existence by degree and institution", err)
	}

	return count > 0, nil
}

func (repo *educationRepository) GetCurrentEducations(ctx context.Context, userID int) ([]*entities.Education, error) {
	query := `SELECT education_id, user_id, education_degree, education_institution, 
			  education_start_date, education_end_date, education_description,
			  education_created_at, education_updated_at 
			  FROM educations WHERE user_id = ? AND (education_end_date IS NULL OR education_end_date = '' OR education_end_date = '0000-00-00')
			  ORDER BY education_start_date DESC`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.logger.Error("Failed to getcurrenteducations: %v", err)
		return nil, domain.NewDatabaseError("retrieve current educations", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var educations []*entities.Education
	for rows.Next() {
		var education entities.Education
		err := rows.Scan(
			&education.EducationID,
			&education.UserID,
			&education.Degree,
			&education.Institution,
			&education.StartDate,
			&education.EndDate,
			&education.Description,
			&education.CreatedAt,
			&education.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning current education: %v", err)
			continue
		}

		educations = append(educations, &education)
	}

	return educations, nil
}

func (repo *educationRepository) GetAll(ctx context.Context) ([]*entities.Education, error) {
	query := `SELECT education_id, user_id, education_degree, education_institution, 
			  education_start_date, education_end_date, education_description,
			  education_created_at, education_updated_at 
			  FROM educations ORDER BY education_start_date DESC`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		repo.logger.Error("Failed to getall educations: %v", err)
		return nil, domain.NewDatabaseError("retrieve all educations", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var educations []*entities.Education
	for rows.Next() {
		var education entities.Education
		err := rows.Scan(
			&education.EducationID,
			&education.UserID,
			&education.Degree,
			&education.Institution,
			&education.StartDate,
			&education.EndDate,
			&education.Description,
			&education.CreatedAt,
			&education.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning education: %v", err)
			continue
		}

		educations = append(educations, &education)
	}

	return educations, nil
}
