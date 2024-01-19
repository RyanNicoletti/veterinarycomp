package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConnector = &DB{}

const maxOpenDbConnections = 10
const maxIdleDbConnections = 5
const maxDbLifeTime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(maxDbLifeTime)
	db.SetMaxIdleConns(maxIdleDbConnections)
	db.SetMaxOpenConns(maxOpenDbConnections)

	dbConnector.SQL = db
	return dbConnector, nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
