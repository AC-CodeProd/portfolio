package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/domain/utils"
	"portfolio/logger"
	"strings"
	"time"
)

type personalInfoRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewPersonalInfoRepository(db *sql.DB, logger *logger.Logger) interfaces.PersonalInfoRepository {
	return &personalInfoRepository{db: db, logger: logger}
}

func (repo *personalInfoRepository) Get(ctx context.Context) (*entities.PersonalInfo, error) {
	var info entities.PersonalInfo
	var dateOfBirth time.Time

	query := `SELECT personal_info_id, user_id, personal_info_first_name, personal_info_last_name, personal_info_professional_title, personal_info_bio, personal_info_location, personal_info_resume_url, personal_info_website_url, personal_info_linkedin_url, personal_info_github_url, personal_info_x_url, personal_info_date_of_birth, personal_info_phone_number, personal_info_interests, personal_info_profile_picture, personal_info_created_at, personal_info_updated_at FROM personal_infos LIMIT 1`
	row := repo.db.QueryRowContext(ctx, query)

	err := row.Scan(
		&info.PersonalInfoID,
		&info.UserID,
		&info.FirstName,
		&info.LastName,
		&info.ProfessionalTitle,
		&info.Bio,
		&info.Location,
		&info.ResumeURL,
		&info.WebsiteURL,
		&info.LinkedinURL,
		&info.GithubURL,
		&info.XURL,
		&dateOfBirth,
		&info.PhoneNumber,
		&info.Interests,
		&info.ProfilePicture,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to get personalinfo: %v", err)
		return nil, domain.NewDatabaseError("retrieve personal information", err)
	}

	date := utils.NewDate(dateOfBirth)
	info.DateOfBirth = &date

	return &info, nil
}

func (repo *personalInfoRepository) GetByID(ctx context.Context, personalInfoID int) (*entities.PersonalInfo, error) {
	var personalInfo entities.PersonalInfo

	var dateOfBirth time.Time
	query := `SELECT personal_info_id, user_id, personal_info_first_name, personal_info_last_name, personal_info_professional_title, personal_info_bio, personal_info_location, personal_info_resume_url, personal_info_website_url, personal_info_linkedin_url, personal_info_github_url, personal_info_x_url, personal_info_date_of_birth, personal_info_phone_number, personal_info_interests, personal_info_profile_picture, personal_info_created_at, personal_info_updated_at FROM personal_infos WHERE personal_info_id = ?`

	row := repo.db.QueryRowContext(ctx, query, personalInfoID)
	err := row.Scan(
		&personalInfo.PersonalInfoID,
		&personalInfo.UserID,
		&personalInfo.FirstName,
		&personalInfo.LastName,
		&personalInfo.ProfessionalTitle,
		&personalInfo.Bio,
		&personalInfo.Location,
		&personalInfo.ResumeURL,
		&personalInfo.WebsiteURL,
		&personalInfo.LinkedinURL,
		&personalInfo.GithubURL,
		&personalInfo.XURL,
		&dateOfBirth,
		&personalInfo.PhoneNumber,
		&personalInfo.Interests,
		&personalInfo.ProfilePicture,
		&personalInfo.CreatedAt,
		&personalInfo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to get personalinfo: %v", err)
		return nil, domain.NewDatabaseError("retrieve personal information", err)
	}

	date := utils.NewDate(dateOfBirth)
	personalInfo.DateOfBirth = &date

	return &personalInfo, nil
}

func (repo *personalInfoRepository) Create(ctx context.Context, personalInfo *entities.PersonalInfo) (*entities.PersonalInfo, error) {
	query := `INSERT INTO personal_infos (user_id, personal_info_first_name, personal_info_last_name, personal_info_professional_title, personal_info_bio, personal_info_location, personal_info_resume_url, personal_info_website_url, personal_info_linkedin_url, personal_info_github_url, personal_info_x_url, personal_info_date_of_birth, personal_info_phone_number, personal_info_interests, personal_info_profile_picture) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	var dateOfBirth *time.Time
	if personalInfo.DateOfBirth != nil {
		t := personalInfo.DateOfBirth.Time()
		dateOfBirth = &t
	}

	_, err := repo.db.ExecContext(ctx, query,
		personalInfo.UserID,
		personalInfo.FirstName,
		personalInfo.LastName,
		personalInfo.ProfessionalTitle,
		personalInfo.Bio,
		personalInfo.Location,
		personalInfo.ResumeURL,
		personalInfo.WebsiteURL,
		personalInfo.LinkedinURL,
		personalInfo.GithubURL,
		personalInfo.XURL,
		dateOfBirth,
		personalInfo.PhoneNumber,
		personalInfo.Interests,
		personalInfo.ProfilePicture,
	)
	if err != nil {
		repo.logger.Error("Failed to create personalinfo: %v", err)
		return nil, domain.NewDatabaseError("create personal information", err)
	}

	return repo.Get(ctx)
}

func (repo *personalInfoRepository) Update(ctx context.Context, personalInfoId int, personalInfo *entities.PersonalInfo) (*entities.PersonalInfo, error) {
	query := `UPDATE personal_infos SET 
		personal_info_first_name = ?, 
		personal_info_last_name = ?, 
		personal_info_professional_title = ?, 
		personal_info_bio = ?, 
		personal_info_location = ?, 
		personal_info_resume_url = ?, 
		personal_info_website_url = ?, 
		personal_info_linkedin_url = ?, 
		personal_info_github_url = ?, 
		personal_info_x_url = ?, 
		personal_info_date_of_birth = ?, 
		personal_info_phone_number = ?, 
		personal_info_interests = ?, 
		personal_info_profile_picture = ? 
	WHERE personal_info_id = ?`

	// Convert utils.Date to time.Time for database storage
	var dateOfBirth *time.Time
	if personalInfo.DateOfBirth != nil {
		t := personalInfo.DateOfBirth.Time()
		dateOfBirth = &t
	}

	_, err := repo.db.ExecContext(ctx, query,
		personalInfo.FirstName,
		personalInfo.LastName,
		personalInfo.ProfessionalTitle,
		personalInfo.Bio,
		personalInfo.Location,
		personalInfo.ResumeURL,
		personalInfo.WebsiteURL,
		personalInfo.LinkedinURL,
		personalInfo.GithubURL,
		personalInfo.XURL,
		dateOfBirth,
		personalInfo.PhoneNumber,
		personalInfo.Interests,
		personalInfo.ProfilePicture,
		personalInfoId,
	)
	if err != nil {
		repo.logger.Error("Failed to update personalinfo: %v", err)
		return nil, domain.NewDatabaseError("update personal information", err)
	}
	return repo.Get(ctx)
}

func (repo *personalInfoRepository) Patch(ctx context.Context, personalInfoId int, personalInfo *entities.PersonalInfo) (*entities.PersonalInfo, error) {
	type field struct {
		Name  string
		Value interface{}
		Set   bool
	}

	var dateOfBirth *time.Time
	if personalInfo.DateOfBirth != nil {
		t := personalInfo.DateOfBirth.Time()
		dateOfBirth = &t
	}

	fields := []field{
		{"personal_info_first_name", personalInfo.FirstName, personalInfo.FirstName != ""},
		{"personal_info_last_name", personalInfo.LastName, personalInfo.LastName != ""},
		{"personal_info_professional_title", personalInfo.ProfessionalTitle, personalInfo.ProfessionalTitle != ""},
		{"personal_info_bio", personalInfo.Bio, personalInfo.Bio != ""},
		{"personal_info_location", personalInfo.Location, personalInfo.Location != ""},
		{"personal_info_resume_url", personalInfo.ResumeURL, personalInfo.ResumeURL != ""},
		{"personal_info_website_url", personalInfo.WebsiteURL, personalInfo.WebsiteURL != ""},
		{"personal_info_linkedin_url", personalInfo.LinkedinURL, personalInfo.LinkedinURL != ""},
		{"personal_info_github_url", personalInfo.GithubURL, personalInfo.GithubURL != ""},
		{"personal_info_x_url", personalInfo.XURL, personalInfo.XURL != ""},
		{"personal_info_date_of_birth", dateOfBirth, personalInfo.DateOfBirth != nil},
		{"personal_info_phone_number", personalInfo.PhoneNumber, personalInfo.PhoneNumber != ""},
		{"personal_info_interests", personalInfo.Interests, personalInfo.Interests != ""},
		{"personal_info_profile_picture", personalInfo.ProfilePicture, personalInfo.ProfilePicture != ""},
	}

	setClauses := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields)+1)

	for _, f := range fields {
		if f.Set {
			setClauses = append(setClauses, f.Name+" = ?")
			args = append(args, f.Value)
		}
	}

	if len(setClauses) == 0 {
		return repo.Get(ctx)
	}

	query := "UPDATE personal_infos SET " + strings.Join(setClauses, ", ") + " WHERE personal_info_id = ?"
	args = append(args, personalInfoId)
	fmt.Println("Executing query:", query, "with args:", args)
	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.logger.Error("Failed to patch personalinfo: %v", err)
		return nil, fmt.Errorf("unable to patch personal information: %w", err)
	}
	return repo.Get(ctx)
}

func (repo *personalInfoRepository) Delete(ctx context.Context, personalInfoId int) error {
	query := `DELETE FROM personal_infos WHERE personal_info_id = ?`
	_, err := repo.db.ExecContext(ctx, query, personalInfoId)
	if err != nil {
		repo.logger.Error("Failed to delete personalinfo: %v", err)
		return domain.NewDatabaseError("delete personal information", err)
	}
	return nil
}

func (repo *personalInfoRepository) GetByUserID(ctx context.Context, userID int) (*entities.PersonalInfo, error) {
	query := `SELECT personal_info_id, user_id, personal_info_first_name, personal_info_last_name, personal_info_professional_title, personal_info_bio, personal_info_location, personal_info_resume_url, personal_info_website_url, personal_info_linkedin_url, personal_info_github_url, personal_info_x_url, personal_info_date_of_birth, personal_info_phone_number, personal_info_interests, personal_info_profile_picture, personal_info_created_at, personal_info_updated_at FROM personal_infos WHERE user_id = ?`
	row := repo.db.QueryRowContext(ctx, query, userID)

	info := &entities.PersonalInfo{}
	var dateOfBirth time.Time

	err := row.Scan(
		&info.PersonalInfoID,
		&info.UserID,
		&info.FirstName,
		&info.LastName,
		&info.ProfessionalTitle,
		&info.Bio,
		&info.Location,
		&info.ResumeURL,
		&info.WebsiteURL,
		&info.LinkedinURL,
		&info.GithubURL,
		&info.XURL,
		&dateOfBirth,
		&info.PhoneNumber,
		&info.Interests,
		&info.ProfilePicture,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to getbyuserid personalinfo: %v", err)
		return nil, domain.NewDatabaseError("retrieve personal information by user ID", err)
	}

	date := utils.NewDate(dateOfBirth)
	info.DateOfBirth = &date

	return info, nil
}

func (repo *personalInfoRepository) ExistsByID(ctx context.Context, personalInfoId int) (bool, error) {
	query := `SELECT COUNT(*) FROM personal_infos WHERE personal_info_id = ?`
	var count int
	err := repo.db.QueryRowContext(ctx, query, personalInfoId).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyid personalinfo: %v", err)
		return false, domain.NewDatabaseError("check if personal information exists by ID", err)
	}
	return count > 0, nil
}

func (repo *personalInfoRepository) ExistsByUserID(ctx context.Context, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM personal_infos WHERE user_id = ?`
	var count int
	err := repo.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		repo.logger.Error("Failed to existsbyuserid personalinfo: %v", err)
		return false, domain.NewDatabaseError("check if personal information exists by user ID", err)
	}
	return count > 0, nil
}

func (repo *personalInfoRepository) GetUserByID(ctx context.Context, userID int) (*entities.User, error) {
	query := `SELECT user_id, user_username, user_email, user_role, user_is_active, user_created_at, user_updated_at FROM users WHERE user_id = ?`
	user := &entities.User{}
	err := repo.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		repo.logger.Error("Failed to getuserbyid: %v", err)
		return nil, domain.NewDatabaseError("retrieve user by ID", err)
	}
	return user, nil
}
