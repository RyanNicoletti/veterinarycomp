package models

import (
	"time"
)

type Compensation struct {
	ID                   int
	CompanyName          string
	JobTitle             string
	PracticeType         string
	BoardCertification   string
	Location             string
	YearsExperience      int
	BaseSalary           float32
	SignOnBonus          float32
	Production           float32
	TotalCompensation    float64
	VerificationDocument *Document
	Verified             bool
	Approved             bool
	IsHourly             bool
	CreatedAt            time.Time
}

type Document struct {
	ID          string
	FileName    string
	ContentType string
	FilePath    string
	CreatedAt   time.Time
}

type Pagination struct {
	Next          int
	Previous      int
	RecordPerPage int
	CurrentPage   int
	TotalPage     int
}

type User struct {
	Id        int
	Email     string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time
}
