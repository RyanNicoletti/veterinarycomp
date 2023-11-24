package repository

import "github.com/ryannicoletti/veterinarycomp/internal/models"

type DatabaseRepo interface {
	GetAllCompensation() ([]models.Compensation, error)
	InsertCompensation(comp models.Compensation) error
}
