package repository

import "github.com/ryannicoletti/veterinarycomp/internal/models"

type CompensationRepo interface {
	GetAllCompensation(limit, offset int) ([]models.Compensation, error)
	InsertCompensation(comp models.Compensation) error
	SearchCompensation(locationOrHospital string, rowPerPage, offset int) ([]models.Compensation, error)
	GetTotalSearchCompensationsCount(locationOrHospital string) (int, error)
	GetTotalCompensationsCount() (int, error)
}

type UserRepo interface {
	GetUserById(id int) (models.User, error)
	Authenticate(email, password string) (int, string, error)
}
