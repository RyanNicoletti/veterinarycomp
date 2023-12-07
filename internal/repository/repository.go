package repository

import "github.com/ryannicoletti/veterinarycomp/internal/models"

type CompensationRepo interface {
	GetAllCompensation(limit, offset int) ([]models.Compensation, error)
	GetCompensationByID(ID int) (models.Compensation, error)
	DeleteCompensationByID(ID int) error
	InsertCompensation(comp models.Compensation) error
	SearchCompensation(locationOrHospital string, rowPerPage, offset int) ([]models.Compensation, error)
	GetTotalSearchCompensationsCount(locationOrHospital string) (int, error)
	GetTotalCompensationsCount() (int, error)
	GetVerificationDocument(ID int) (*models.Document, error)
	GetDocumentMetaDataById(ID int) (models.Compensation, error)
	DeleteCompensationDocumentByID(ID int) error
	VerifyComp(ID int) error
}

type UserRepo interface {
	GetUserById(id int) (models.User, error)
	Authenticate(email, password string) (int, string, error)
	IsAdmin(id int) (bool, error)
}
