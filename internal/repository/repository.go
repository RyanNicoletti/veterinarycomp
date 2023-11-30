package repository

import "github.com/ryannicoletti/veterinarycomp/internal/models"

type CompensationDatabaseRepo interface {
	GetAllCompensation(limit, offset int) ([]models.Compensation, error)
	InsertCompensation(comp models.Compensation) error
	SearchCompensation(locationOrHospital string, rowPerPage, offset int) ([]models.Compensation, error)
	GetTotalSearchCompensationsCount(locationOrHospital string) (int, error)
	GetTotalCompensationsCount() (int, error)
}
