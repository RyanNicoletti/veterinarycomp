package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/driver"
	"github.com/ryannicoletti/veterinarycomp/internal/handlers"
	"github.com/ryannicoletti/veterinarycomp/internal/helpers"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
)

const portNumber string = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
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
	app.IsProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	log.Println("Connecting to database...")

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=vetcomp user=ryannicoletti password=")
	if err != nil {
		fmt.Println("Failed to connect to db, ggs...")
	}
	log.Println("Successfully connected to database.")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Failed to create template cache", err)
	}
	app.TemplateCache = tc
	// change to true when in prod
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
