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

type skillRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewSkillRepository(db *sql.DB, logger *logger.Logger) interfaces.SkillRepository {
	return &skillRepository{db: db, logger: logger}
}

func (repo *skillRepository) Create(ctx context.Context, skill *entities.Skill) (*entities.Skill, error) {
	query := `INSERT INTO skills (user_id, skill_name, skill_level) VALUES (?, ?, ?)`

	result, err := repo.db.ExecContext(ctx, query,
		skill.UserID,
		skill.Name,
		skill.Level,
	)
	if err != nil {
		repo.logger.Error("Failed to create skill: %v", err)
		return nil, domain.NewDatabaseError("create skill", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to create skill lastinsertid: %v", err)
		return nil, domain.NewDatabaseError("get skill ID", err)
	}

	skill.SkillID = int(id)
	skill.CreatedAt = time.Now()
	skill.UpdatedAt = time.Now()

	return skill, nil
}

func (repo *skillRepository) GetByID(ctx context.Context, skillID int) (*entities.Skill, error) {
	var skill entities.Skill
	query := `SELECT skill_id, user_id, skill_name, skill_level, skill_created_at, skill_updated_at 
			  FROM skills WHERE skill_id = ?`

	row := repo.db.QueryRowContext(ctx, query, skillID)
	err := row.Scan(
		&skill.SkillID,
		&skill.UserID,
		&skill.Name,
		&skill.Level,
		&skill.CreatedAt,
		&skill.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to get skill: %v", err)
		return nil, domain.NewDatabaseError("retrieve skill", err)
	}

	return &skill, nil
}

func (repo *skillRepository) GetByUserID(ctx context.Context, userID int) ([]*entities.Skill, error) {
	query := `SELECT skill_id, user_id, skill_name, skill_level, skill_created_at, skill_updated_at 
			  FROM skills WHERE user_id = ? ORDER BY skill_name`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.logger.Error("Failed to getbyuserid skills: %v", err)
		return nil, domain.NewDatabaseError("retrieve skills by user ID", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var skills []*entities.Skill
	for rows.Next() {
		var skill entities.Skill
		err := rows.Scan(
			&skill.SkillID,
			&skill.UserID,
			&skill.Name,
			&skill.Level,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning skill: %v", err)
			continue
		}
		skills = append(skills, &skill)
	}

	return skills, nil
}

func (repo *skillRepository) Update(ctx context.Context, skillID int, skill *entities.Skill) (*entities.Skill, error) {
	query := `UPDATE skills SET skill_name = ?, skill_level = ?, skill_updated_at = ? WHERE skill_id = ?`

	_, err := repo.db.ExecContext(ctx, query,
		skill.Name,
		skill.Level,
		time.Now(),
		skillID,
	)
	if err != nil {
		repo.logger.Error("Failed to update skill: %v", err)
		return nil, domain.NewDatabaseError("update skill", err)
	}

	return repo.GetByID(ctx, skillID)
}

func (repo *skillRepository) Patch(ctx context.Context, skillID int, skill *entities.Skill) (*entities.Skill, error) {
	type field struct {
		Name  string
		Value interface{}
		Set   bool
	}

	var fields []field
	if skill.Name != "" {
		fields = append(fields, field{"skill_name", skill.Name, true})
	}
	if skill.Level > 0 {
		fields = append(fields, field{"skill_level", skill.Level, true})
	}

	if len(fields) == 0 {
		return skill, nil
	}

	query := "UPDATE skills SET "
	var args []interface{}
	for i, f := range fields {
		query += f.Name + " = ?"
		args = append(args, f.Value)
		if i < len(fields)-1 {
			query += ", "
		}
	}
	query += " WHERE skill_id = ?"
	args = append(args, skillID)

	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.logger.Error("Failed to patch skill: %v", err)
		return nil, fmt.Errorf("unable to patch skill: %w", err)
	}
	return repo.GetByID(ctx, skillID)
}

func (repo *skillRepository) Delete(ctx context.Context, skillID int) error {
	query := `DELETE FROM skills WHERE skill_id = ?`

	_, err := repo.db.ExecContext(ctx, query, skillID)
	if err != nil {
		repo.logger.Error("Failed to delete skill: %v", err)
		return domain.NewDatabaseError("delete skill", err)
	}

	return nil
}

func (repo *skillRepository) ExistsByID(ctx context.Context, skillID int) (bool, error) {
	query := `SELECT COUNT(*) FROM skills WHERE skill_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, skillID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyid skill: %v", err)
		return false, domain.NewDatabaseError("check skill existence", err)
	}

	return count > 0, nil
}

func (repo *skillRepository) ExistsByNameAndUserID(ctx context.Context, name string, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM skills WHERE skill_name = ? AND user_id = ?`
	var count int

	err := repo.db.QueryRowContext(ctx, query, name, userID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbynameanduserid skill: %v", err)
		return false, domain.NewDatabaseError("check skill existence by name", err)
	}

	return count > 0, nil
}

func (repo *skillRepository) GetAll(ctx context.Context) ([]*entities.Skill, error) {
	query := `SELECT skill_id, user_id, skill_name, skill_level, skill_created_at, skill_updated_at 
			  FROM skills ORDER BY skill_name`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		repo.logger.Error("Failed to getall skills: %v", err)
		return nil, domain.NewDatabaseError("retrieve all skills", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Error("Failed to closing rows: %v", err)
		}
	}()

	var skills []*entities.Skill
	for rows.Next() {
		var skill entities.Skill
		err := rows.Scan(
			&skill.SkillID,
			&skill.UserID,
			&skill.Name,
			&skill.Level,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)
		if err != nil {
			repo.logger.Error("Failed to scanning skill: %v", err)
			continue
		}
		skills = append(skills, &skill)
	}

	return skills, nil
}
