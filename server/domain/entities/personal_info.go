package entities

import (
	"portfolio/domain/utils"
	"time"
)

type PersonalInfo struct {
	DateOfBirth       *utils.Date
	CreatedAt         time.Time
	UpdatedAt         time.Time
	FirstName         string
	LastName          string
	ProfessionalTitle string
	Bio               string
	Location          string
	ResumeURL         string
	WebsiteURL        string
	LinkedinURL       string
	GithubURL         string
	XURL              string
	PhoneNumber       string
	Interests         string
	ProfilePicture    string
	PersonalInfoID    int
	UserID            int
}

func (pi *PersonalInfo) HasRequiredFields() bool {
	return pi.FirstName != "" && pi.LastName != "" && pi.UserID > 0
}

func (pi *PersonalInfo) GetFullName() string {
	return pi.FirstName + " " + pi.LastName
}

func (pi *PersonalInfo) BelongsToUser(userID int) bool {
	return pi.UserID == userID
}

func (pi *PersonalInfo) MarkAsUpdated() {
	pi.UpdatedAt = time.Now()
}

func (pi *PersonalInfo) HasProfilePicture() bool {
	return pi.ProfilePicture != ""
}

func (pi *PersonalInfo) HasSocialLinks() bool {
	return pi.LinkedinURL != "" || pi.GithubURL != "" || pi.XURL != "" || pi.WebsiteURL != ""
}
