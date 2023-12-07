package repositoryimpl

import (
	"context"
	"database/sql"
	"encoding/json"
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
func NewPostgresCompensationRepo(conn *sql.DB, a *config.AppConfig) repository.CompensationRepo {
	return &pgCompensationRepo{
		App: a,
		DB:  conn,
	}
}

func (dbRepo *pgCompensationRepo) GetAllCompensation(limit, offset int) ([]models.Compensation, error) {
	// if we cant insert within 3 seconds, cancel the transaction
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select * from compensations order by id limit $1 offset $2`
	rows, err := dbRepo.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var compensations = []models.Compensation{}
	for rows.Next() {
		var compensation = models.Compensation{}
		var dbyte []byte
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
			&dbyte,
			&compensation.Verified,
			&compensation.CreatedAt)
		if err != nil {
			return nil, err
		}
		if len(dbyte) > 0 {
			var d *models.Document
			err := json.Unmarshal(dbyte, &d)
			if err != nil {
				return nil, err
			}
			compensation.VerificationDocument = d
		} else {
			compensation.VerificationDocument = nil
		}

		compensations = append(compensations, compensation)
	}
	return compensations, nil
}

func (dbRepo *pgCompensationRepo) InsertCompensation(comp models.Compensation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	jsonDoc, e := json.Marshal(comp.VerificationDocument)
	if e != nil {
		return e
	}

	query := `INSERT INTO compensations (company_name, job_title, type_of_practice, board_certification, location, years_of_experience, base_salary, sign_on_bonus, production, total_comp, verification_document, verified, date_created)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := dbRepo.DB.ExecContext(ctx, query,
		comp.CompanyName, comp.JobTitle, comp.PracticeType, comp.BoardCertification, comp.Location,
		comp.YearsExperience, comp.BaseSalary, comp.SignOnBonus, comp.Production, comp.TotalCompensation,
		jsonDoc, comp.Verified, time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (dbRepo *pgCompensationRepo) SearchCompensation(locationOrHospital string, rowPerPage, offset int) ([]models.Compensation, error) {
	var compensations = []models.Compensation{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM compensations
    WHERE location ~* '\m` + locationOrHospital + `\M' OR company_name ~* '\m` + locationOrHospital + `\M'
	order by id limit $1 offset $2`
	rows, err := dbRepo.DB.QueryContext(ctx, query, rowPerPage, offset)
	if err != nil {
		return compensations, err
	}
	defer rows.Close()
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
			return compensations, err
		}
		compensations = append(compensations, compensation)
	}
	return compensations, nil
}

func (dbRepo *pgCompensationRepo) GetTotalSearchCompensationsCount(locationOrHospital string) (int, error) {
	query := `SELECT COUNT(*) FROM compensations WHERE location ~* '\m` + locationOrHospital + `\M' OR company_name ~* '\m` + locationOrHospital + `\M'`
	var count int
	if err := dbRepo.DB.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (dbRepo *pgCompensationRepo) GetTotalCompensationsCount() (int, error) {
	var total int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT COUNT(id) FROM compensations`
	err := dbRepo.DB.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (dbRepo *pgCompensationRepo) GetVerificationDocument(ID int) (*models.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var documentData []byte
	query := `SELECT verification_document FROM compensations WHERE id = $1`
	err := dbRepo.DB.QueryRowContext(ctx, query, ID).Scan(&documentData)
	if err != nil {
		return nil, err
	}
	var document models.Document
	err = json.Unmarshal(documentData, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (dbRepo *pgCompensationRepo) GetDocumentMetaDataById(ID int) (models.Compensation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var c models.Compensation
	var b []byte
	query := `SELECT verification_document from compensations WHERE id = $1`
	err := dbRepo.DB.QueryRowContext(ctx, query, ID).Scan(&b)
	if err != nil {
		return c, err
	}
	if len(b) > 0 {
		var d *models.Document
		err := json.Unmarshal(b, &d)
		if err != nil {
			return c, err
		}
		c.VerificationDocument = d
	}
	return c, nil
}

func (dbRepo *pgCompensationRepo) VerifyComp(ID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `UPDATE compensations SET verified = TRUE WHERE id = $1`
	_, err := dbRepo.DB.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}
	return nil
}

func (dbRepo *pgCompensationRepo) GetCompensationByID(ID int) (models.Compensation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * from compensations WHERE id = $1`
	var c models.Compensation
	var verificationDocument sql.NullString
	err := dbRepo.DB.QueryRowContext(ctx, query, ID).Scan(&c.ID,
		&c.CompanyName,
		&c.JobTitle,
		&c.PracticeType,
		&c.BoardCertification,
		&c.Location,
		&c.YearsExperience,
		&c.BaseSalary,
		&c.SignOnBonus,
		&c.Production,
		&c.TotalCompensation,
		&verificationDocument,
		&c.Verified,
		&c.CreatedAt)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (dbRepo *pgCompensationRepo) DeleteCompensationByID(ID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `DELETE FROM compensations WHERE id = $1`
	_, err := dbRepo.DB.ExecContext(ctx, query, ID)
	// handle errors here better-what if no row is deleted?
	if err != nil {
		return err
	}
	return nil
}
