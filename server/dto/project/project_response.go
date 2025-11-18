package dto

import (
	"portfolio/domain/entities"
	"portfolio/shared"
	"time"
)

// @Description Project represents a project entry in the portfolio
type Project struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	Technologies     string    `json:"technologies"`
	GithubURL        string    `json:"github_url"`
	ImageURL         string    `json:"image_url"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
} // @name Project

// @Description Response for a list of projects
type ProjectListResponse struct {
	Projects []*Project   `json:"projects"`
	Meta     *shared.Meta `json:"meta"`
} //@name ProjectListResponse

// @Description Response for a project
type ProjectResponse struct {
	Project *Project     `json:"project"`
	Meta    *shared.Meta `json:"meta"`
} //@name ProjectResponse

func FromProjectEntityToResponse(project *entities.Project, meta *shared.Meta) *ProjectResponse {
	if project == nil {
		return nil
	}

	return &ProjectResponse{
		Project: &Project{
			ID:               project.ProjectID,
			UserID:           project.UserID,
			Title:            project.Title,
			Description:      project.Description,
			ShortDescription: project.ShortDescription,
			Technologies:     project.Technologies,
			GithubURL:        project.GithubURL,
			ImageURL:         project.ImageURL,
			Status:           project.Status,
			CreatedAt:        project.CreatedAt,
			UpdatedAt:        project.UpdatedAt,
		},
		Meta: meta,
	}
}

func FromProjectsEntityToResponse(projects []*entities.Project, meta *shared.Meta) *ProjectListResponse {
	if projects == nil {
		return nil
	}

	var projectResponses []*Project
	for _, project := range projects {
		_project := &Project{
			ID:               project.ProjectID,
			UserID:           project.UserID,
			Title:            project.Title,
			Description:      project.Description,
			ShortDescription: project.ShortDescription,
			Technologies:     project.Technologies,
			GithubURL:        project.GithubURL,
			ImageURL:         project.ImageURL,
			Status:           project.Status,
			CreatedAt:        project.CreatedAt,
			UpdatedAt:        project.UpdatedAt,
		}

		projectResponses = append(projectResponses, _project)
	}

	return &ProjectListResponse{
		Projects: projectResponses,
		Meta:     meta,
	}
}
