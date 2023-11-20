package models

import "time"

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
	ProductionSalary     float32
	TotalSalary          float32
	VerificationDocument Document
	Verified             bool
	CreatedAt            time.Time
}

type Document struct {
	Data        []byte
	FileName    string
	ContentType string
	CreatedAt   time.Time
}
