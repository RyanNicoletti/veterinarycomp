package main

import (
	"log"
	"net/http"

	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/handlers"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig

func main() {
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}
