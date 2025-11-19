package dto

import (
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/utils"
)

type CreatePersonalInfoRequest struct {
	DateOfBirth       string `json:"date_of_birth"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	ProfessionalTitle string `json:"professional_title"`
	Bio               string `json:"bio"`
	Location          string `json:"location"`
	ResumeURL         string `json:"resume_url"`
	WebsiteURL        string `json:"website_url"`
	LinkedinURL       string `json:"linkedin_url"`
	GithubURL         string `json:"github_url"`
	XURL              string `json:"x_url"`
	PhoneNumber       string `json:"phone_number"`
	Interests         string `json:"interests"`
	ProfilePicture    string `json:"profile_picture"`
} // @name CreatePersonalInfoRequest

type UpdatePersonalInfoRequest struct {
	CreatePersonalInfoRequest
}

type PatchPersonalInfoRequest struct {
	FirstName         *string `json:"first_name,omitempty"`
	LastName          *string `json:"last_name,omitempty"`
	ProfessionalTitle *string `json:"professional_title,omitempty"`
	Bio               *string `json:"bio,omitempty"`
	Location          *string `json:"location,omitempty"`
	ResumeURL         *string `json:"resume_url,omitempty"`
	WebsiteURL        *string `json:"website_url,omitempty"`
	LinkedinURL       *string `json:"linkedin_url,omitempty"`
	GithubURL         *string `json:"github_url,omitempty"`
	XURL              *string `json:"x_url,omitempty"`
	DateOfBirth       *string `json:"date_of_birth,omitempty"`
	PhoneNumber       *string `json:"phone_number,omitempty"`
	Interests         *string `json:"interests,omitempty"`
	ProfilePicture    *string `json:"profile_picture,omitempty"`
}

func FromPersonalInfoRequestToEntity(userID int, req *CreatePersonalInfoRequest) *entities.PersonalInfo {
	if req == nil {
		return nil
	}

	dateOfBirth, err := utils.ParseDate(req.DateOfBirth)
	if err != nil {
		return nil
	}

	return &entities.PersonalInfo{
		UserID:            userID,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		ProfessionalTitle: req.ProfessionalTitle,
		Bio:               req.Bio,
		Location:          req.Location,
		ResumeURL:         req.ResumeURL,
		WebsiteURL:        req.WebsiteURL,
		LinkedinURL:       req.LinkedinURL,
		GithubURL:         req.GithubURL,
		XURL:              req.XURL,
		DateOfBirth:       dateOfBirth,
		PhoneNumber:       req.PhoneNumber,
		Interests:         req.Interests,
		ProfilePicture:    req.ProfilePicture,
	}
}

func FromUpdatePersonalInfoRequestToEntity(userID int, req *UpdatePersonalInfoRequest) *entities.PersonalInfo {
	if req == nil {
		return nil
	}

	dateOfBirth, err := utils.ParseDate(req.DateOfBirth)
	if err != nil {
		return nil
	}

	return &entities.PersonalInfo{
		UserID:            userID,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		ProfessionalTitle: req.ProfessionalTitle,
		Bio:               req.Bio,
		Location:          req.Location,
		ResumeURL:         req.ResumeURL,
		WebsiteURL:        req.WebsiteURL,
		LinkedinURL:       req.LinkedinURL,
		GithubURL:         req.GithubURL,
		XURL:              req.XURL,
		DateOfBirth:       dateOfBirth,
		PhoneNumber:       req.PhoneNumber,
		Interests:         req.Interests,
		ProfilePicture:    req.ProfilePicture,
	}
}

func FromPatchPersonalInfoRequestToEntity(userID int, req *PatchPersonalInfoRequest) (*entities.PersonalInfo, error) {
	if req == nil {
		return nil, domain.NewValidationError("Invalid request", "request", nil)
	}

	dateOfBirth, err := utils.ParseDate(getStringValue(req.DateOfBirth))
	if err != nil && req.DateOfBirth != nil {
		return nil, domain.NewInvalidFormatError("date_of_birth", "YYYY-MM-DD")
	}

	return &entities.PersonalInfo{
		UserID:            userID,
		FirstName:         getStringValue(req.FirstName),
		LastName:          getStringValue(req.LastName),
		ProfessionalTitle: getStringValue(req.ProfessionalTitle),
		Bio:               getStringValue(req.Bio),
		Location:          getStringValue(req.Location),
		ResumeURL:         getStringValue(req.ResumeURL),
		WebsiteURL:        getStringValue(req.WebsiteURL),
		LinkedinURL:       getStringValue(req.LinkedinURL),
		GithubURL:         getStringValue(req.GithubURL),
		XURL:              getStringValue(req.XURL),
		DateOfBirth:       dateOfBirth,
		PhoneNumber:       getStringValue(req.PhoneNumber),
		Interests:         getStringValue(req.Interests),
		ProfilePicture:    getStringValue(req.ProfilePicture),
	}, nil
}

func getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}
