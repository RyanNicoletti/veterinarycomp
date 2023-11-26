package repository

import "github.com/ryannicoletti/veterinarycomp/internal/models"

type CompensationDatabaseRepo interface {
	GetAllCompensation() ([]models.Compensation, error)
	InsertCompensation(comp models.Compensation) error
	SearchCompensation(locationOrHospital string) ([]models.Compensation, error)
}
