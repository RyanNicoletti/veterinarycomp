package repositoryimpl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
	"github.com/ryannicoletti/veterinarycomp/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type pgUserRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// returns a pointer to an instance of postgresDBRepo
// this works because postgresDBRepo implements DatabaseRepo, which
// is the declared return type of the function
func NewPostgresUserRepo(conn *sql.DB, a *config.AppConfig) repository.UserRepo {
	return &pgUserRepo{
		App: a,
		DB:  conn,
	}
}

func (dbRepo *pgUserRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * from users where id = $1`

	row := dbRepo.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(&u.Id, &u.Email, &u.Password, &u.CreatedAt)

	if err != nil {
		return u, err
	}
	return u, nil
}

func (dbRepo *pgUserRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `select id, password from users where email = $1`
	row := dbRepo.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Incorrect email or password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
