package utils

import "time"

type Date time.Time

func NewDate(t time.Time) Date {
	return Date(t)
}

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

func ParseDate(dateStr string) (*Date, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	d := Date(t)
	return &d, nil
}
