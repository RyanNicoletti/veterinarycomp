package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
	"github.com/ryannicoletti/veterinarycomp/internal/repository"
)

type pgCompensationRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// returns a pointer to an instance of postgresDBRepo
// this works because postgresDBRepo implements DatabaseRepo, which
// is the declared return type of the function
func NewPostgresCompensationRepo(conn *sql.DB, a *config.AppConfig) repository.CompensationDatabaseRepo {
	return &pgCompensationRepo{
		App: a,
		DB:  conn,
	}
}

func (dbRepo *pgCompensationRepo) GetAllCompensation() ([]models.Compensation, error) {
	// if we cant insert within 3 seconds, cancel the transaction
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select * from compensations`
	rows, err := dbRepo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var compensations = []models.Compensation{}
	for rows.Next() {
		var compensation = models.Compensation{}
		err := rows.Scan(&compensation.ID,
			&compensation.CompanyName,
			&compensation.JobTitle,
			&compensation.PracticeType,
			&compensation.BoardCertification,
			&compensation.Location,
			&compensation.YearsExperience,
			&compensation.BaseSalary,
			&compensation.SignOnBonus,
			&compensation.Production,
			&compensation.TotalCompensation,
			&compensation.VerificationDocument,
			&compensation.Verified,
			&compensation.CreatedAt)
		if err != nil {
			return nil, err
		}
		compensations = append(compensations, compensation)
	}
	return compensations, nil
}

func (dbRepo *pgCompensationRepo) InsertCompensation(comp models.Compensation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into compensations (company_name, job_title, type_of_practice, board_certification, location, years_of_experience, base_salary, sign_on_bonus, production, total_comp, verification_document, verified, date_created)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := dbRepo.DB.ExecContext(ctx, query, comp.CompanyName, comp.JobTitle, comp.PracticeType, comp.BoardCertification, comp.Location, comp.YearsExperience, comp.BaseSalary, comp.SignOnBonus, comp.Production, comp.TotalCompensation, comp.VerificationDocument, comp.Verified, time.Now())
	if err != nil {
		return err
	}

	return nil
}
