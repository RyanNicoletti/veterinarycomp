package repository

import "github.com/ryannicoletti/veterinarycomp/internal/models"

type CompensationDatabaseRepo interface {
	GetAllCompensation(limit, offset int) ([]models.Compensation, error)
	InsertCompensation(comp models.Compensation) error
	SearchCompensation(locationOrHospital string) ([]models.Compensation, error)
	GetTotalCompensationsCount() (int, error)
}
