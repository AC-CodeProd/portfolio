package dto

import (
	"portfolio/domain/entities"
	"portfolio/shared"
)

// @Description PersonalInfo represents the personal information of a user in the portfolio
type PersonalInfo struct {
	PersonalInfoID    int    `json:"personal_info_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	ProfessionalTitle string `json:"professional_title"`
	Description       string `json:"description"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	LinkedIn          string `json:"linkedin"`
	GitHub            string `json:"github"`
	X                 string `json:"x"`
	WebsiteURL        string `json:"website_url"`
	ResumeURL         string `json:"resume_url"`
	Location          string `json:"location"`
	DateOfBirth       string `json:"date_of_birth"`
	Interests         string `json:"interests"`
	ProfilePicture    string `json:"profile_picture"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
} // @name PersonalInfo

// @Description Response for personal information
type PersonalInfoResponse struct {
	PersonalInfo PersonalInfo `json:"personal_info"`
	Meta         *shared.Meta `json:"meta"`
} //@name PersonalInfoResponse

func FromPersonalInfoEntityToResponse(user *entities.User, personalInfo *entities.PersonalInfo, meta *shared.Meta) *PersonalInfoResponse {
	if personalInfo == nil || user == nil {
		return nil
	}

	return &PersonalInfoResponse{
		PersonalInfo: PersonalInfo{
			PersonalInfoID:    personalInfo.PersonalInfoID,
			FirstName:         personalInfo.FirstName,
			LastName:          personalInfo.LastName,
			ProfessionalTitle: personalInfo.ProfessionalTitle,
			Description:       personalInfo.Bio,
			Email:             user.Email,
			Phone:             personalInfo.PhoneNumber,
			LinkedIn:          personalInfo.LinkedinURL,
			GitHub:            personalInfo.GithubURL,
			X:                 personalInfo.XURL,
			WebsiteURL:        personalInfo.WebsiteURL,
			ResumeURL:         personalInfo.ResumeURL,
			Location:          personalInfo.Location,
			DateOfBirth:       personalInfo.DateOfBirth.String(),
			Interests:         personalInfo.Interests,
			ProfilePicture:    personalInfo.ProfilePicture,
			CreatedAt:         personalInfo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:         personalInfo.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Meta: meta,
	}
}

func NewPersonalInfoResponse(user *entities.User, personalInfo *entities.PersonalInfo, meta *shared.Meta) *PersonalInfoResponse {
	return FromPersonalInfoEntityToResponse(user, personalInfo, meta)
}
