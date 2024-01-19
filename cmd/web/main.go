package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pressly/goose"
	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/driver"
	"github.com/ryannicoletti/veterinarycomp/internal/handlers"
	"github.com/ryannicoletti/veterinarycomp/internal/helpers"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
)

const portNumber string = ":8080"

var (
	app      config.AppConfig
	session  *scs.SessionManager
	infoLog  *log.Logger
	errorLog *log.Logger
	dbHost   = os.Getenv("VETCOMP_DB_HOST")
	dbName   = os.Getenv("VETCOMP_DB_NAME")
	dbUser   = os.Getenv("VETCOMP_DB_USER")
	dbPass   = os.Getenv("VETCOMP_DB_PASSWORD")
	dbPort   = os.Getenv("VETCOMP_DB_PORT")
)

func main() {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads", 0755)
		if err != nil {
			fmt.Println("Error creating upload directory:", err)
			os.Exit(1)
		}
	}

	gob.Register(models.Compensation{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	db, err := start()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Listening on port %s\n", portNumber)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func start() (*driver.DB, error) {
	app.IsProduction = os.Getenv("VETCOMP_IS_PROD") == "true"

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	log.Println("Connecting to database...")

	dbConnectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", dbHost, dbPort, dbName, dbUser, dbPass)
	db, err := driver.ConnectSQL(dbConnectionString)
	if err != nil {
		fmt.Println("Failed to connect to db, ggs...")
	}
	log.Println("Successfully connected to database.")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Failed to create template cache", err)
	}
	app.TemplateCache = tc
	// in production we want to use cache, otherwise, dont need to
	app.UseCache = os.Getenv("VETCOMP_IS_PROD") == "true"

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	// run migrations
	sqlDb, err := driver.NewDatabase(dbConnectionString)
	if err != nil {
		log.Fatal("Failed to create db.")
	}
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlDb, "./db/migrations"); err != nil {
		panic(err)
	}

	return db, nil
}
