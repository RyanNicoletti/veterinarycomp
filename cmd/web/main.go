package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/handlers"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
)

const portNumber string = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.IsProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.CreateNewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	fmt.Printf("Listening on port %s\n", portNumber)
	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
