package entities

import (
	"fmt"
	"time"
)

type Education struct {
	EducationID int
	UserID      int
	Degree      string
	Institution string
	StartDate   time.Time
	EndDate     *time.Time
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e *Education) HasRequiredFields() bool {
	return e.Degree != "" && e.Institution != "" && e.UserID > 0
}

func (e *Education) BelongsToUser(userID int) bool {
	return e.UserID == userID
}

func (e *Education) MarkAsUpdated() {
	e.UpdatedAt = time.Now()
}

func (e *Education) IsCurrentStudies() bool {
	return e.EndDate == nil || e.EndDate.IsZero()
}

func (e *Education) GetDurationInMonths() int {
	var endDate time.Time
	if e.IsCurrentStudies() {
		endDate = time.Now()
	} else {
		endDate = *e.EndDate
	}

	start := e.StartDate
	if start.IsZero() {
		return 0
	}

	years := endDate.Year() - start.Year()
	months := int(endDate.Month()) - int(start.Month())

	return years*12 + months
}

func (e *Education) GetFormattedDuration() string {
	months := e.GetDurationInMonths()
	if months == 0 {
		return "Less than a month"
	}

	years := months / 12
	remainingMonths := months % 12

	if years == 0 {
		if remainingMonths == 1 {
			return "1 month"
		}
		return fmt.Sprintf("%d months", remainingMonths)
	}

	if remainingMonths == 0 {
		if years == 1 {
			return "1 year"
		}
		return fmt.Sprintf("%d years", years)
	}

	yearText := "year"
	if years > 1 {
		yearText = "years"
	}
	monthText := "month"
	if remainingMonths > 1 {
		monthText = "months"
	}

	return fmt.Sprintf("%d %s %d %s", years, yearText, remainingMonths, monthText)
}

func (e *Education) SetEndDate(endDate *time.Time) {
	e.EndDate = endDate
	e.MarkAsUpdated()
}

func (e *Education) GetFullDescription() string {
	return fmt.Sprintf("%s at %s", e.Degree, e.Institution)
}
